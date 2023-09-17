/*
Copyright © 2023 NAME HERE yusufmirza55@hotmail.com
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/EtrusChain/synnefo/config"
	"github.com/EtrusChain/synnefo/core"
	"github.com/EtrusChain/synnefo/core/bootstrap"
	"github.com/EtrusChain/synnefo/core/node"
	"github.com/EtrusChain/synnefo/p2p"
	"github.com/EtrusChain/synnefo/peering"
	"github.com/EtrusChain/synnefo/repo"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
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
			"/ip4/192.168.0.11/tcp/5200/p2p/QmX7jAWE95GidPbrdwFof326TGbbg7nuDFFgzHJh7EmzKm",
			"/ip4/192.168.0.11/udp/5200/quic-v1/p2p/QmX7jAWE95GidPbrdwFof326TGbbg7nuDFFgzHJh7EmzKm",
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

		bootstrapPeers, err := config.DefaultBootstrapPeers()
		if err != nil {
			panic(err)
		}

		bootsrapConfig := bootstrap.BootstrapConfigWithPeers(bootstrapPeers)

		err = n.Bootstrap(bootsrapConfig)
		if err != nil {
			panic(err)
		}

		peering.NewPeeringService(node)

		bootstrapPeerss, err := sd.BootstrapPeers()
		if err != nil {
			return
		}

		peering := peering.NewPeeringService(node)
		peering.Start()
		peering.AddPeer(bootstrapPeers[0])
		listPeers := peering.ListPeers()
		fmt.Println(listPeers)

		// PUB/SUB
		gossipSub, err := pubsub.NewGossipSub(ctx, node)
		if err != nil {
			panic(err)
		}

		hostInfor := host.InfoFromHost(node)
		fmt.Println(hostInfor)

		mDNS := mdns.NewMdnsService(node, ServiceName, &discoveryNotifee{h: node})

		if err != nil {
			log.Fatal(err)
		}

		mDNS.Start()

		room := "/synnefo/daemon/" + node.ID().String()
		topic, err := gossipSub.Join(room)
		if err != nil {
			panic(err)
		}

		publish(ctx, topic)

		if bootstrapPeerss[0].ID != node.ID() {
			bootstrapRoom := "/synnefo/daemon/QmX7jAWE95GidPbrdwFof326TGbbg7nuDFFgzHJh7EmzKm"

			bootstapTopic, err := gossipSub.Join(bootstrapRoom)
			if err != nil {
				panic(err)
			}
			// subscribe to topic
			subscriber, err := bootstapTopic.Subscribe()
			if err != nil {
				panic(err)
			}

			subscribe(subscriber, ctx, node.ID())

			/*
				peerMA, err := multiaddr.NewMultiaddr("/ip4/192.168.0.11/tcp/5200/p2p/QmX7jAWE95GidPbrdwFof326TGbbg7nuDFFgzHJh7EmzKm")
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
				s, err := node.NewStream(context.Background(), bootstrapPeerss[0].ID, "/synnefo/1.0.0")
				if err != nil {
					panic(err)
				}

				go writeCounter(s)
				go readCounter(s)
			*/
		}

		sigCh := make(chan os.Signal)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		<-sigCh

		//select {}
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

func publish(ctx context.Context, topic *pubsub.Topic) {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Printf("enter message to publish: \n")

			msg := scanner.Text()
			if len(msg) != 0 {
				// publish message to topic
				bytes := []byte(msg)
				topic.Publish(ctx, bytes)
			}
		}
	}
}

func subscribe(subscriber *pubsub.Subscription, ctx context.Context, hostID peer.ID) {
	for {
		msg, err := subscriber.Next(ctx)
		if err != nil {
			panic(err)
		}

		// only consider messages delivered by other peers
		if msg.ReceivedFrom == hostID {
			continue
		}

		fmt.Printf("got message: %s, from: %s\n", string(msg.Data), msg.ReceivedFrom.Pretty())
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
