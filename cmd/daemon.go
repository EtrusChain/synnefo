/*
Copyright © 2023 NAME HERE yusufmirza55@hotmail.com
*/
package cmd

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/EtrusChain/synnefo/config"
	"github.com/EtrusChain/synnefo/core"
	"github.com/EtrusChain/synnefo/core/bootstrap"
	"github.com/EtrusChain/synnefo/core/node"
	"github.com/EtrusChain/synnefo/p2p"
	"github.com/EtrusChain/synnefo/peering"
	"github.com/EtrusChain/synnefo/repo"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		const (
			ServiceName = "_p2p._udp"
		)

		fmt.Println("daemon called")
		ctx := context.Background()
		node, err := node.NewNode(ctx)
		if err != nil {
			panic(err)
		}

		log.Println(node.Addrs())
		p2pHost := p2p.New(node.ID(), node, node.Peerstore())

		nodePeering := peering.NewPeeringService(node)

		a := []string{
			"/dnsaddr/bootstrap.libp2p.io/p2p/QmX7jAWE95GidPbrdwFof326TGbbg7nuDFFgzHJh7EmzKm",
			"/ip4/178.233.168.239/tcp/5200/p2p/QmX7jAWE95GidPbrdwFof326TGbbg7nuDFFgzHJh7EmzKm",
			"/ip4/178.233.168.239/udp/5200/quic-v1/p2p/QmX7jAWE95GidPbrdwFof326TGbbg7nuDFFgzHJh7EmzKm",
		}

		sd := config.Config{
			Identity: config.Identity{
				PeerID:  string(node.ID()),
				PrivKey: string(""),
			},
			Datastore: config.Datastore{},
			Addresses: config.Addresses{
				Swarm: []string{"/ip4/0.0.0.0/tcp/4001", "/ip4/0.0.0.0/udp/4001/quic-v1"},
			},
			Discovery: config.Discovery{
				MDNS: config.MDNS{
					Enabled: true,
				},
			},
			Bootstrap: a,

			Peering: config.Peering{
				Peers: nodePeering.ListPeers(),
			},
			DNS: config.DNS{
				Resolvers:   map[string]string{},
				MaxCacheTTL: &config.OptionalDuration{},
			},

			Internal: config.Internal{
				Bitswap: &config.InternalBitswap{
					TaskWorkerCount:             config.OptionalInteger{},
					EngineBlockstoreWorkerCount: config.OptionalInteger{},
					EngineTaskWorkerCount:       config.OptionalInteger{},
					MaxOutstandingBytesPerPeer:  config.OptionalInteger{},
					ProviderSearchDelay:         config.OptionalDuration{},
				},
				UnixFSShardingSizeThreshold: &config.OptionalString{},
				Libp2pForceReachability:     &config.OptionalString{},
				BackupBootstrapInterval:     &config.OptionalDuration{},
			},
		}

		repo := &repo.Mock{
			C: sd,
		}

		n := &core.SynnefoNode{
			Identity: node.ID(),

			Repo: repo,

			PubKey:     node.Peerstore().PubKey(node.ID()),
			PrivateKey: node.Peerstore().PrivKey(node.ID()),

			PeerHost: node,
			Peering:  &peering.PeeringService{},

			P2P: p2pHost,

			IsOnline: true,
			IsDaemon: true,
		}

		fmt.Println(n)

		bootstrapPeers, err := config.DefaultBootstrapPeers()
		if err != nil {
			panic(err)
		}

		bootsrapConfig := bootstrap.BootstrapConfigWithPeers(bootstrapPeers)

		err = n.Bootstrap(bootsrapConfig)
		if err != nil {
			panic(err)
		}

		fmt.Println(bootstrapPeers)
		fmt.Println(bootsrapConfig)
		peering.NewPeeringService(node)

		/*
			listeners := &p2p.Listeners{
				Listeners: map[protocol.ID]p2p.Listener{},
			}

			name := listeners.Register(listeners.Listeners["/x/"])
			errorr := name.Error()
			fmt.Println(errorr)
		*/

		/*
			check := p2pHost.CheckProtoExists("/x/")
			fmt.Println(check)

			nodeRegister := p2pHost.ListenersP2P
			defer nodeRegister.Register(p2pHost.ListenersLocal.Listeners[protocol.ID("/synnefo/1.0.0")])

			fmt.Println(nodeRegister)

			nodePeering.AddPeer(node.Peerstore().PeerInfo(node.ID()))
			defer peering.NewPeeringService(node)
			defer nodePeering.Start()
			defer nodePeering.Stop()

			defer node.Connect(ctx, bootstrapPeerss[0])
			listPeers := nodePeering.ListPeers()

			fmt.Println(listPeers)
		*/

		bootstrapPeerss, err := sd.BootstrapPeers()
		if err != nil {
			return
		}

		peering := peering.NewPeeringService(node)
		peering.Start()
		peering.AddPeer(bootstrapPeers[0])
		listPeers := peering.ListPeers()
		fmt.Println(listPeers)

		host.Host.SetStreamHandler(node, "/libp2p/circuit/relay/0.2.0/hop", func(s network.Stream) {
			go writeCounter(s)
			go readCounter(s)

			s.Close()
		})

		hostInfor := host.InfoFromHost(node)
		fmt.Println(hostInfor)

		mDNS := mdns.NewMdnsService(node, ServiceName, &discoveryNotifee{h: node})

		if err != nil {
			log.Fatal(err)
		}

		defer mDNS.Close()

		if bootstrapPeerss[0].ID != node.ID() {
			fmt.Println("Non Bootstrap Peer")

			peerMA, err := multiaddr.NewMultiaddr("/ip4/178.233.168.239/tcp/5200/p2p/" + node.ID().String() + "/p2p-circuit/p2p/" + "QmX7jAWE95GidPbrdwFof326TGbbg7nuDFFgzHJh7EmzKm")
			if err != nil {
				panic(err)
			}
			peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
			if err != nil {
				panic(err)
			}

			if err := node.Connect(context.Background(), *peerAddrInfo); err != nil {
				panic(err)
			}

			fmt.Println("Connected to", peerAddrInfo.String())
			s, err := node.NewStream(context.Background(), bootstrapPeerss[0].ID, "/libp2p/circuit/relay/0.2.0/hop")
			if err != nil {
				panic(err)
			}

			defer s.Close()

			go writeCounter(s)
			go readCounter(s)
		}

		select {}
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// daemonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daemonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func writeCounter(s network.Stream) {
	var counter uint64

	for {
		<-time.After(time.Second)
		counter++

		err := binary.Write(s, binary.BigEndian, counter)
		if err != nil {
			panic(err)
		}
	}

}

func readCounter(s network.Stream) {
	for {
		var counter uint64

		err := binary.Read(s, binary.BigEndian, &counter)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Received %d from %s\n", counter, s.ID())
	}

}

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", pi.ID.Pretty())
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
	}
}
