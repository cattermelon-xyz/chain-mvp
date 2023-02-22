/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package client

import (
	"fmt"

	"github.com/spf13/cobra"
)

// voteCmd represents the vote command
var voteCmd = &cobra.Command{
	Use:   "vote",
	Short: "vote -i {id string} -o {option int}",
	Long:  `vote -i {id string} -o {option int}.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vote called")
	},
}

func init() {
	rootCmd.AddCommand(voteCmd)
}
