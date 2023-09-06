/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("profile called")

		db, err := leveldb.OpenFile("user/db", nil)
		if err != nil {
			panic(err)
		}

		defer db.Close()

		pub, err := db.Get([]byte("pubKey"), nil)
		if err != nil {
			return
		}

		priv, err := db.Get([]byte("privKey"), nil)
		if err != nil {
			return
		}

		parsedPrivateKey, err := crypto.UnmarshalPrivateKey(priv)
		if err != nil {
			fmt.Println("Error generating private key:", err)
			return
		}

		fmt.Println("Database userData: ", string(pub), string(priv))
		host, err := libp2p.New(libp2p.Identity(parsedPrivateKey))
		if err != nil {
			return
		}

		fmt.Println("Peer id: ", host.ID().String())

	},
}

func init() {
	rootCmd.AddCommand(profileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// profileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// profileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
