/*
Copyright © 2023 NAME HERE <yusufmirza55@hotmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
)

// rootNodeCmd represents the rootNode command
var rootNodeCmd = &cobra.Command{
	Use:   "rootNode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		node, err := libp2p.New(
			libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/3423"),
			libp2p.Ping(false),
		)
		if err != nil {
			panic(err)
		}

		// configure our own ping protocol
		pingService := &ping.PingService{Host: node}
		node.SetStreamHandler(ping.ID, pingService.PingHandler)

		// wait for a SIGINT or SIGTERM signal
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("Received signal, shutting down...")

		// shut the node down
		if err := node.Close(); err != nil {
			panic(err)
		}

		peerInfo := peerstore.AddrInfo{
			ID:    node.ID(),
			Addrs: node.Addrs(),
		}
		addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
		if err != nil {
			return
		}

		fmt.Println("libp2p node address:", addrs[0])
	},
}

func init() {
	rootCmd.AddCommand(rootNodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rootNodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rootNodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
