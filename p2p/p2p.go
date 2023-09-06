package p2p

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/syndtr/goleveldb/leveldb"
)

/*
type P2PNode struct {
	host host.Host
}
*/

func CreateLibp2pHost(ctx context.Context) string {
	db, err := leveldb.OpenFile("user/db", nil)
	if err != nil {
		return ""
	}

	data, err := db.Get([]byte("privKey"), nil)
	if err != nil {
		return ""
	}

	defer db.Close()
	privateKey, err := crypto.UnmarshalPrivateKey(data)
	if err != nil {
		return ""
	}
	// Create a list of libp2p options, including the DHT option
	opts := []libp2p.Option{
		libp2p.DisableRelay(),     // Disable relay (optional)
		libp2p.EnableNATService(), // Enable NAT service (optional)
		// libp2p.EnableAutoRelayWithPeerSource(),  // Enable auto relay (optional)
		libp2p.EnableNATService(), // Enable NAT port mapping (optional)
		libp2p.DefaultTransports,  // Use default transports (optional)
		libp2p.NATPortMap(),
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Security("synnefo", ctx),
		libp2p.Identity(privateKey),
		libp2p.Ping(true),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			dht, err := dht.New(ctx, h)

			return dht, err
		}),
	}

	// Create the libp2p host with the DHT option
	host, err := libp2p.New(opts...)
	if err != nil {
		return ""
	}

	fmt.Println("Hello world")
	fmt.Println(host.ID().String())
	fmt.Println(host.ID().Pretty())

	log.Fatalln(host.ID().Pretty())
	// Attach the DHT to the host
	dht, err := dht.New(ctx, host)
	if err != nil {
		return ""
	}

	// Attach the DHT to the host as a routing option
	host.Peerstore().AddAddrs(host.ID(), host.Addrs(), peerstore.PermanentAddrTTL)
	err = dht.Bootstrap(ctx)
	if err != nil {
		return ""
	}

	return host.ID().String()
}

/*
func privKeyToString(privKey crypto.PrivKey) (string, error) {
	// Serialize the private key to a DER format
	privKeyBytes, err := privKey.Raw()
	if err != nil {
		return "", err
	}

	// Create a PEM block
	privKeyPEM := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privKeyBytes,
	}

	// Encode the PEM block to a string
	privKeyStr := string(pem.EncodeToMemory(privKeyPEM))

	return privKeyStr, nil
}
*/
