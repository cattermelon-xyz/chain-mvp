/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package client

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "htg-client",
	Short: "Client program to connect Hectagon Network",
	Long: `Client program to connect Hectagon Network.
	Default option is localhost:8813`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringP("network", "n", "localhost:8813", "Decide what network:port this client would connect to. Default is localhost:8813")
}
