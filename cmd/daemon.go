/*
Copyright © 2023 NAME HERE <Yusuf Mirza Pıçakcı>
*/
package cmd

import (
	"github.com/EtrusChain/synnefo/node"
	"github.com/spf13/cobra"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start a long-running daemon process",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		peerAddress, _ := cmd.Flags().GetString("peerAddress")
		node.RunNode(peerAddress)
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	rootCmd.PersistentFlags().String("rootNode", "false", "A search term for a rootNode.")
	rootCmd.PersistentFlags().String("peerAddress", "", "Connect for a rootNode.")
	rootCmd.PersistentFlags().String("maxByte", "", "A search term for a rootNode.")
}
