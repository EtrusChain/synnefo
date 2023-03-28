package node

import (
	"fmt"

	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/syndtr/goleveldb/leveldb"
)

func RunNode() {
	db, err := leveldb.OpenFile("db/userid", nil)
	if err != nil {
		return
	}

	defer db.Close()

	userID, err := db.Get([]byte("userID"), nil)
	if err != nil {
		return
	}
	fmt.Println(string(userID))

	// start a libp2p node that listens on a random local TCP port,
	// but without running the built-in ping protocol
	node, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/3423/p2p/"),
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}

	// configure our own ping protocol
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	// print the node's PeerInfo in multiaddr format
	peerInfo := peerstore.AddrInfo{
		ID:    peerstore.ID(string(userID)),
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println("libp2p node address:", addrs[0])
}
