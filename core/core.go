package core

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/EtrusChain/synnefo/config"
	"github.com/EtrusChain/synnefo/core/bootstrap"
	"github.com/EtrusChain/synnefo/p2p"
	"github.com/EtrusChain/synnefo/peering"
	"github.com/EtrusChain/synnefo/repo"
	irouting "github.com/EtrusChain/synnefo/routing"
	goprocess "github.com/jbenet/goprocess"
	ddht "github.com/libp2p/go-libp2p-kad-dht/dual"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	psrouter "github.com/libp2p/go-libp2p-pubsub-router"
	"github.com/libp2p/go-libp2p/core/connmgr"
	ic "github.com/libp2p/go-libp2p/core/crypto"
	p2phost "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/metrics"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	p2pbhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	madns "github.com/multiformats/go-multiaddr-dns"
)

// synnefoNode is Synnefo Core module. It represents an Synnefo instance.
type SynnefoNode struct {
	// Self
	Identity peer.ID // the local node's identity

	Repo repo.Repo

	// Services
	Reporter  *metrics.BandwidthCounter `optional:"true"`
	Discovery mdns.Service              `optional:"true"`

	// Discovery            mdns.Service              `optional:"true"`
	PubKey     ic.PubKey
	PrivateKey ic.PrivKey `optional:"true"` // the local node's private Key

	// Online
	PeerHost        p2phost.Host               `optional:"true"` // the network host (server+client)
	Peering         *peering.PeeringService    `optional:"true"`
	Bootstrapper    io.Closer                  `optional:"true"` // the periodic bootstrapper
	DNSResolver     *madns.Resolver            // the DNS resolver
	ResourceManager network.ResourceManager    `optional:"true"`
	Routing         irouting.ProvideManyRouter `optional:"true"` // the routing system. recommend synnefo-dht

	PubSub   *pubsub.PubSub             `optional:"true"`
	PSRouter *psrouter.PubsubValueStore `optional:"true"`

	DHT       *ddht.DHT       `optional:"true"`
	DHTClient routing.Routing `name:"dhtc" optional:"true"`

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

func (n *SynnefoNode) Bootstrap(cfg bootstrap.BootstrapConfig) error {
	var err error
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

func (n *SynnefoNode) loadBootstrapPeers() ([]peer.AddrInfo, error) {
	cfg, err := n.Repo.Config()
	if err != nil {
		return nil, err
	}

	return cfg.BootstrapPeers()
}

func (n *SynnefoNode) saveTempBootstrapPeers(ctx context.Context, peerList []peer.AddrInfo) error {
	db, err := repo.NewDatabaseHandler(repo.GetOs())
	if err != nil {
		panic(err)
	}

	defer db.Close()

	bytes, err := json.Marshal(config.BootstrapPeerStrings(peerList))
	if err != nil {
		return err
	}

	if err := db.SetValue("TempBootstrapPeersKey", bytes); err != nil {
		return err
	}
	return err
}

func (n *SynnefoNode) loadTempBootstrapPeers(ctx context.Context) ([]peer.AddrInfo, error) {
	db, err := repo.NewDatabaseHandler(repo.GetOs())
	if err != nil {
		panic(err)
	}

	defer db.Close()

	bytes, err := db.GetValue("TempBootstrapPeersKey")
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
