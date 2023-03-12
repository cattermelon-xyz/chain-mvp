/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package client

import (
	"log"

	"github.com/spf13/cobra"
)

// nodeCmd represents the node command
var checkpointDetachCmd = &cobra.Command{
	Use:   "detach",
	Short: "To detach a CheckPoint from an existed one",
	Long:  `To detach a CheckPoint from an existed one`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		option, _ := cmd.Flags().GetInt("option")
		log.Printf("Attach a CheckPoint %s to %d", id, option)
	},
}

func init() {
	checkpointCmd.AddCommand(checkpointDetachCmd)
	checkpointCmd.Flags().IntP("option", "o", 0, "Use with `detach`. Index of option to detach")
}
