/*
Copyright © 2023 NAME HERE yusufmirza55@hotmail.com
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/EtrusChain/synnefo/config"
	"github.com/EtrusChain/synnefo/core"
	"github.com/EtrusChain/synnefo/core/bootstrap"
	"github.com/EtrusChain/synnefo/core/node"
	"github.com/EtrusChain/synnefo/p2p"
	"github.com/EtrusChain/synnefo/peering"
	"github.com/EtrusChain/synnefo/repo"
	"github.com/libp2p/go-libp2p/core/protocol"
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
		fmt.Println("daemon called")
		ctx := context.Background()
		node, err := node.NewNode(ctx)
		if err != nil {
			panic(err)
		}

		p2pHost := p2p.New(node.ID(), node, node.Peerstore())

		nodePeering := peering.NewPeeringService(node)

		a := []string{
			"/ip4/178.233.168.239/tcp/4001/p2p/QmX7jAWE95GidPbrdwFof326TGbbg7nuDFFgzHJh7EmzKm",         // mars.i.ipfs.io
			"/ip4/178.233.168.239/udp/4001/quic-v1/p2p/QmX7jAWE95GidPbrdwFof326TGbbg7nuDFFgzHJh7EmzKm", // mars.i.ipfs.io
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
		check := p2pHost.CheckProtoExists("/x/")
		fmt.Println(check)

		sfdsf := p2pHost.ListenersP2P
		defer sfdsf.Register(p2pHost.ListenersLocal.Listeners[protocol.ID("/x/")])

		fmt.Println(sfdsf)

		nodePeering.AddPeer(node.Peerstore().PeerInfo(node.ID()))

		nodePeering.Start()
		nodePeering.Stop()

		listPeers := nodePeering.ListPeers()

		fmt.Println(listPeers)
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
