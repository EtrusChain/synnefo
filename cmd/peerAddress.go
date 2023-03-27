/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/EtrusChain/synnefo/peering"
	"github.com/spf13/cobra"
)

// peerAddressCmd represents the peerAddress command
var peerAddressCmd = &cobra.Command{
	Use:   "peerAddress",
	Short: "Check you peer address",
	Long:  `Check you peer address`,
	Run: func(cmd *cobra.Command, args []string) {
		peering.GetPeer()
	},
}

func init() {
	rootCmd.AddCommand(peerAddressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// peerAddressCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// peerAddressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
