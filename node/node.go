package node

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func RunNode(peerAddress string) {

	rootNode := [2]string{
		"/ip4/192.168.0.15/tcp/2454/p2p/12D3KooWGf23wxWLHEoypJCrvooWW2okKXeKJsMvtGRmyFgbywa9",
		"/ip4/192.168.0.15/tcp/3434/p2p/12D3KooWCVzjCLKFALBwKwkHKXhYsMbt1iAfhV86tvvXguYtYQ2r",
	}

	node, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/2454"),
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
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("libp2p node address:", addrs[0])
	fmt.Println("rootNode Address :", rootNode[0])

	if len(peerAddress) > 1 {
		addr, err := multiaddr.NewMultiaddr(peerAddress)
		if err != nil {
			panic(err)
		}
		peer, err := peerstore.AddrInfoFromP2pAddr(addr)
		if err != nil {
			panic(err)
		}
		if err := node.Connect(context.Background(), *peer); err != nil {
			panic(err)
		}
		fmt.Println("sending 5 ping messages to", addr)
		ch := pingService.Ping(context.Background(), peer.ID)
		for i := 0; i < 5; i++ {
			res := <-ch
			fmt.Println("pinged", addr, "in", res.RTT)
		}
	} else {
		// wait for a SIGINT or SIGTERM signal
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("Received signal, shutting down...")
	}

	// shut the node down
	if err := node.Close(); err != nil {
		panic(err)
	}
}
