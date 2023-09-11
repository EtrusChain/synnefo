/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/EtrusChain/synnefo/repo"
	"github.com/spf13/cobra"
)

// identityCmd represents the identity command
var identityCmd = &cobra.Command{
	Use:   "identity",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := repo.NewDatabaseHandler(repo.GetOs())
		if err != nil {
			panic(err)
		}

		defer db.Close()
		data, err := db.GetValue("identity")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(data))
	},
}

func init() {
	rootCmd.AddCommand(identityCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// identityCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// identityCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
