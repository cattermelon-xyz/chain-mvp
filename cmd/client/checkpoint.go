/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package client

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nodeCmd represents the node command
var checkpointCmd = &cobra.Command{
	Use:   "checkpoint",
	Short: "To setup CheckPoint",
	Long:  `To setup CheckPoint with subcommands: new, attach, detach`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteId, _ := cmd.Flags().GetString("delete")
		if deleteId != "" {
			fmt.Println("Try to delete ", deleteId)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkpointCmd)
	checkpointCmd.PersistentFlags().StringP("id", "i", "", "The ID of the CheckPoint")

	checkpointCmd.Flags().StringP("delete", "d", "", "Use without any command. The ID of CheckPoint to delete")
}
