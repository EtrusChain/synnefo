package node

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	ic "github.com/libp2p/go-libp2p/core/crypto"
	p2phost "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/syndtr/goleveldb/leveldb"
)

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

	fmt.Println("Peer Address: ", host.ID().String())
	return host.ID().String(), nil
}

func LoadNode(ctx context.Context) (string, error) {
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
