/*
Copyright © 2023 NAME HERE yusufmirza55@hotmail.com
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/EtrusChain/synnefo/config"
	"github.com/EtrusChain/synnefo/repo"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long:  `A longer description.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := repo.NewDatabaseHandler(repo.GetOs())
		if err != nil {
			panic(err)
		}

		defer db.Close()

		data, err := db.CheckValue("identity")
		if err != nil {
			panic(err)
		}
		if data {
			fmt.Println("Peer identity already exists.")
			return
		}

		var identity config.Identity

		identity, err = config.CreateIdentity(os.Stdout, []config.KeyGenerateOption{
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

		db.SetValue("identity", []byte(jsonData))
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
