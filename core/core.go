package core

import (
	"context"
	"fmt"
	"io"

	"github.com/EtrusChain/synnefo/p2p"
	"github.com/EtrusChain/synnefo/peering"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	ddht "github.com/libp2p/go-libp2p-kad-dht/dual"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	psrouter "github.com/libp2p/go-libp2p-pubsub-router"
	"github.com/libp2p/go-libp2p/core/crypto"
	ic "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/metrics"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	madns "github.com/multiformats/go-multiaddr-dns"
	"github.com/syndtr/goleveldb/leveldb"
)

// synnefoNode is IPFS Core module. It represents an IPFS instance.
type SynnefoNode struct {
	// Self
	Identity peer.ID // the local node's identity

	// Services
	Reporter *metrics.BandwidthCounter `optional:"true"`
	// Discovery            mdns.Service              `optional:"true"`
	PrivateKey ic.PrivKey `optional:"true"` // the local node's private Key

	// Online
	Peering         *peering.PeeringService `optional:"true"`
	Bootstrapper    io.Closer               `optional:"true"` // the periodic bootstrapper
	DNSResolver     *madns.Resolver         // the DNS resolver
	ResourceManager network.ResourceManager `optional:"true"`

	PubSub   *pubsub.PubSub             `optional:"true"`
	PSRouter *psrouter.PubsubValueStore `optional:"true"`

	DHT       *ddht.DHT       `optional:"true"`
	DHTClient routing.Routing `name:"dhtc" optional:"true"`

	P2P *p2p.P2P `optional:"true"`

	ctx context.Context

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

	privKey, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
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
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
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

	privateKeyBytes, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return "", err
	}

	db.Put([]byte("privKey"), privateKeyBytes, nil)
	db.Put([]byte("pubKey"), []byte(host.ID().String()), nil)

	return host.ID().String(), nil
}
