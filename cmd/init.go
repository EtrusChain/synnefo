/*
Copyright © 2023 NAME HERE yusufmirza55@hotmail.com
*/
package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/EtrusChain/synnefo/config"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long:  `A longer description.`,
	Run: func(cmd *cobra.Command, args []string) {
		var identity config.Identity

		identity, err := config.CreateIdentity(os.Stdout, []config.KeyGenerateOption{
			config.Key.Size(-1),
			config.Key.Type("rsa"),
		})
		if err != nil {
			panic(err)
		}

		conf, err := config.InitWithIdentity(identity)
		if err != nil {
			panic(err)
		}

		jsonData, err := json.MarshalIndent(conf, "", " ")
		if err != nil {
			log.Fatal(err)
		}

		db, err := leveldb.OpenFile("user/db", nil)
		if err != nil {
			panic(err)
		}

		defer db.Close()

		db.Put([]byte("identity"), []byte(jsonData), nil)
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
