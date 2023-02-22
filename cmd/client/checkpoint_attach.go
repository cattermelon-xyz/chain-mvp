/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package client

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nodeCmd represents the node command
var checkpointAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "To attach a CheckPoint to an existed one",
	Long:  `To attach a CheckPoint to an existed one`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		to, _ := cmd.Flags().GetString("to")
		fmt.Printf("Attach a CheckPoint %s to %s", id, to)
	},
}

func init() {
	checkpointCmd.AddCommand(checkpointAttachCmd)
	checkpointCmd.Flags().StringP("to", "t", "", "Use with `attach`. The ID of CheckPoint to attach to")
}
