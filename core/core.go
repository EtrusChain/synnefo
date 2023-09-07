package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/EtrusChain/synnefo/config"
	"github.com/EtrusChain/synnefo/core/bootstrap"
	"github.com/EtrusChain/synnefo/p2p"
	"github.com/EtrusChain/synnefo/peering"
	"github.com/ipfs/go-datastore"
	goprocess "github.com/jbenet/goprocess"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	ddht "github.com/libp2p/go-libp2p-kad-dht/dual"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	psrouter "github.com/libp2p/go-libp2p-pubsub-router"
	"github.com/libp2p/go-libp2p/core/connmgr"
	ic "github.com/libp2p/go-libp2p/core/crypto"
	p2phost "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/metrics"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	p2pbhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	madns "github.com/multiformats/go-multiaddr-dns"
	"github.com/multiformats/go-multihash"
	"github.com/syndtr/goleveldb/leveldb"
)

type ProvideManyRouter interface {
	ProvideMany(ctx context.Context, keys []multihash.Multihash) error
}

// synnefoNode is IPFS Core module. It represents an IPFS instance.
type SynnefoNode struct {
	// Self
	Identity peer.ID // the local node's identity

	// Services
	Reporter  *metrics.BandwidthCounter `optional:"true"`
	Discovery mdns.Service              `optional:"true"`

	// Discovery            mdns.Service              `optional:"true"`
	PubKey     ic.PubKey
	PrivateKey ic.PrivKey `optional:"true"` // the local node's private Key

	// Online
	PeerHost        p2phost.Host            `optional:"true"` // the network host (server+client)
	Peering         *peering.PeeringService `optional:"true"`
	Bootstrapper    io.Closer               `optional:"true"` // the periodic bootstrapper
	DNSResolver     *madns.Resolver         // the DNS resolver
	ResourceManager network.ResourceManager `optional:"true"`
	Routing         ProvideManyRouter       `optional:"true"` // the routing system. recommend ipfs-dht

	PubSub   *pubsub.PubSub             `optional:"true"`
	PSRouter *psrouter.PubsubValueStore `optional:"true"`

	DHT       *ddht.DHT       `optional:"true"`
	DHTClient routing.Routing `name:"dhtc" optional:"true"`
	Db        *leveldb.DB     `optional:"true"`

	P2P *p2p.P2P `optional:"true"`

	Process goprocess.Process
	ctx     context.Context

	stop func() error

	// Flags
	IsOnline bool `optional:"true"` // Online is set when networking is enabled.
	IsDaemon bool `optional:"true"` // Daemon is set when running on a long-running daemon.
}

// Close calls Close() on the App object
func (n *SynnefoNode) Close() error {
	return n.stop()
}

func (n *SynnefoNode) Context() context.Context {
	if n.ctx == nil {
		n.ctx = context.TODO()
	}
	return n.ctx
}

func CreateLibp2pHost(ctx context.Context) (string, error) {
	db, err := leveldb.OpenFile("user/db", nil)
	if err != nil {
		return "", err
	}

	defer db.Close()

	privKey, _, err := ic.GenerateKeyPair(ic.Ed25519, -1)
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return "", err
	}
	// Create a list of libp2p options, including the DHT option
	opts := []libp2p.Option{
		libp2p.DisableRelay(),     // Disable relay (optional)
		libp2p.EnableNATService(), // Enable NAT service (optional)
		libp2p.EnableNATService(), // Enable NAT port mapping (optional)
		libp2p.DefaultTransports,  // Use default transports (optional)
		libp2p.NATPortMap(),
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Identity(privKey),
		libp2p.Ping(true),
		libp2p.Routing(func(h p2phost.Host) (routing.PeerRouting, error) {
			dht, err := dht.New(ctx, h)

			return dht, err
		}),
		// libp2p.EnableAutoRelayWithPeerSource(),  // Enable auto relay (optional)
		// libp2p.Security("synnefo", ctx),
	}

	fmt.Println(privKey)
	// Create the libp2p host with the DHT option
	host, err := libp2p.New(opts...)
	if err != nil {
		return "", err
	}

	// Attach the DHT to the host
	dht, err := dht.New(ctx, host)
	if err != nil {
		return "", err
	}

	// Attach the DHT to the host as a routing option
	host.Peerstore().AddAddrs(host.ID(), host.Addrs(), peerstore.PermanentAddrTTL)
	err = dht.Bootstrap(ctx)
	if err != nil {
		return "", err
	}

	privateKeyBytes, err := ic.MarshalPrivateKey(privKey)
	if err != nil {
		return "", err
	}

	db.Put([]byte("privKey"), privateKeyBytes, nil)
	db.Put([]byte("pubKey"), []byte(host.ID().String()), nil)

	return host.ID().String(), nil
}

func (n *SynnefoNode) Bootstrap(cfg bootstrap.BootstrapConfig) error {

	if n.Bootstrapper != nil {
		n.Bootstrapper.Close() // stop previous bootstrap process.
	}
	// if the caller did not specify a bootstrap peer function, get the
	// freshest bootstrap peers from config. this responds to live changes.
	if cfg.BootstrapPeers == nil {
		cfg.BootstrapPeers = func() []peer.AddrInfo {
			ps, err := n.loadBootstrapPeers()
			if err != nil {
				return nil
			}
			return ps
		}
	}
	if cfg.SaveBackupBootstrapPeers == nil {
		cfg.SaveBackupBootstrapPeers = func(ctx context.Context, peerList []peer.AddrInfo) {
			err := n.saveTempBootstrapPeers(ctx, peerList)
			if err != nil {
				return
			}
		}
	}
	if cfg.LoadBackupBootstrapPeers == nil {
		cfg.LoadBackupBootstrapPeers = func(ctx context.Context) []peer.AddrInfo {
			peerList, err := n.loadTempBootstrapPeers(ctx)
			if err != nil {
				return nil
			}
			return peerList
		}
	}

	repoConf, err := n.Repo.Config()
	if err != nil {
		return err
	}
	if repoConf.Internal.BackupBootstrapInterval != nil {
		cfg.BackupBootstrapInterval = repoConf.Internal.BackupBootstrapInterval.WithDefault(time.Hour)
	}

	n.Bootstrapper, err = bootstrap.Bootstrap(n.Identity, n.PeerHost, n.Routing, cfg)

	return err
}

var TempBootstrapPeersKey = datastore.NewKey("/local/temp_bootstrap_peers")

func (n *SynnefoNode) loadBootstrapPeers() ([]peer.AddrInfo, error) {
	data, err := n.Db.Get([]byte("BootstrapPeers"), nil)
	if err != nil {
		return nil, err
	}

	cfg, err := n.Repo.Config()
	if err != nil {
		return nil, err
	}

	return cfg.BootstrapPeers()
}

func (n *SynnefoNode) saveTempBootstrapPeers(ctx context.Context, peerList []peer.AddrInfo) error {
	ds := n.Repo.Datastore()
	bytes, err := json.Marshal(config.BootstrapPeerStrings(peerList))
	if err != nil {
		return err
	}

	if err := ds.Put(ctx, TempBootstrapPeersKey, bytes); err != nil {
		return err
	}
	return ds.Sync(ctx, TempBootstrapPeersKey)
}

func (n *SynnefoNode) loadTempBootstrapPeers(ctx context.Context) ([]peer.AddrInfo, error) {
	ds := n.Repo.Datastore()
	bytes, err := ds.Get(ctx, TempBootstrapPeersKey)
	if err != nil {
		return nil, err
	}

	var addrs []string
	if err := json.Unmarshal(bytes, &addrs); err != nil {
		return nil, err
	}
	return config.ParseBootstrapPeers(addrs)
}

type ConstructPeerHostOpts struct {
	AddrsFactory      p2pbhost.AddrsFactory
	DisableNatPortMap bool
	DisableRelay      bool
	EnableRelayHop    bool
	ConnectionManager connmgr.ConnManager
}
