/*
Copyright © 2023 NAME HERE <yusufmirza55@hotmail.com>
*/
package cmd

import (
	"github.com/EtrusChain/synnefo/node"
	"github.com/spf13/cobra"
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
		node.Bootstrap()
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
