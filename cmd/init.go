/*
Copyright © 2023 NAME HERE yusufmirza55@hotmail.com
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/EtrusChain/synnefo/core"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long:  `A longer description.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")

		ctx := context.Background()
		data, err := core.CreateLibp2pHost(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println("user Data n:", data)
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
