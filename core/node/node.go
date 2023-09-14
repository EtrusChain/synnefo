package node

import (
	"context"
	"fmt"

	"github.com/EtrusChain/synnefo/config"
	"github.com/EtrusChain/synnefo/repo"
	"github.com/libp2p/go-libp2p"
	p2phost "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
)

func NewNode(ctx context.Context) (p2phost.Host, error) {
	db, err := repo.NewDatabaseHandler(repo.GetOs())
	if err != nil {
		return nil, err
	}

	defer db.Close()

	pubKey, err := db.GetValue("peerID")
	if err != nil {
		return nil, err
	}
	privKey, err := db.GetValue("privKey")
	if err != nil {
		return nil, err
	}

	sd := &config.Identity{
		PeerID:  string(pubKey),
		PrivKey: string(privKey),
	}

	key, err := sd.DecodePrivateKey(sd.PrivKey)
	if err != nil {
		return nil, err
	}

	// Create a list of libp2p options, including the DHT option
	opts := []libp2p.Option{
		//libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/5200"),
		//libp2p.Transport(tcp.NewTCPTransport),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/5200"),
		libp2p.EnableRelay(),
		libp2p.EnableNATService(), // Enable NAT service (optional)
		libp2p.DefaultTransports,  // Use default transports (optional)
		libp2p.NATPortMap(),
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Identity(key),
		libp2p.Ping(true),
		//libp2p.Security("/x/", ctx),
		/*
			libp2p.Routing(func(h p2phost.Host) (routing.PeerRouting, error) {
				dht, err := dht.New(ctx, h)

				return dht, err
			}),
		*/
		// libp2p.EnableAutoRelayWithPeerSource(),  // Enable auto relay (optional)
	}

	// Create the libp2p host with the DHT option
	host, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	/*
		// Attach the DHT to the host
		dht, err := dht.New(ctx, host)
		if err != nil {
			return nil, err
		}

		// Attach the DHT to the host as a routing option
		host.Peerstore().AddAddrs(host.ID(), host.Addrs(), peerstore.PermanentAddrTTL)
		err = dht.Bootstrap(ctx)
		if err != nil {
			return nil, err
		}
	*/

	host.Peerstore().AddAddrs(host.ID(), host.Addrs(), peerstore.PermanentAddrTTL)
	fmt.Println("Peer Address: ", host.ID().String())
	return host, nil
}
