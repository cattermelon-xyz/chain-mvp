/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package client

import (
	"log"

	"github.com/spf13/cobra"
)

// nodeCmd represents the node command
var initiativeCmd = &cobra.Command{
	Use:   "initiative",
	Short: "To setup Initiative",
	Long: `To setup Initiative with option [build, start, stop, clean]
Note: option "build" will clear the root checkpoint and all of its children`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		option, _ := cmd.Flags().GetString("option")
		if option != "" {
			log.Println("Try ", option)
			if id != "" && option == "build" {
				log.Println("Invalid parameters")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initiativeCmd)
	initiativeCmd.Flags().StringP("id", "i", "", "The ID of the Initiative")
	initiativeCmd.Flags().StringP("option", "o", "", "The option with [build, start, stop, clean, display]")
	initiativeCmd.Flags().StringP("from", "f", "", "The Id of root Checkpoint")
}
