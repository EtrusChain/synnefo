/*
Copyright © 2023 NAME HERE <yusufmirza55@hotmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/EtrusChain/synnefo/p2p"
	"github.com/libp2p/go-libp2p"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create your account for run your node in etrusChain Network",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")

		db, err := leveldb.OpenFile("db/userid", nil)
		if err != nil {
			return
		}
		defer db.Close()

		node, err := libp2p.New()
		if err != nil {
			panic(err)
		}

		// print the node's listening addresses
		fmt.Println("Listen addresses:", node.ID().Pretty())

		data, err := db.Get([]byte("userID"), nil)
		fmt.Println(data)
		if err != nil {
			return
		}

		userData := p2p.New(
			node.ID(),
			node,
			node.Peerstore(),
		)

		err = db.Put([]byte("userID"), []byte(node.ID().Pretty()), nil)
		if err != nil {
			return
		}

		fmt.Println(userData.ListenersP2P, userData)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
