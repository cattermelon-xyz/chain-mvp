/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package client

import (
	"encoding/json"
	"log"

	"github.com/hectagon-finance/chain-mvp/third_party/utils"
	"github.com/spf13/cobra"
)

// nodeCmd represents the node command
var checkpointNewCmd = &cobra.Command{
	Use:   "new",
	Short: "To create a new CheckPoint",
	Long:  `To create a new CheckPoint`,
	Run: func(cmd *cobra.Command, args []string) {
		filepath, _ := cmd.Flags().GetString("path")
		if filepath == "" {
			log.Println("Cannot input an empty file. Let's preproces it before pushing to the server")
		}
	},
}

func init() {
	checkpointCmd.AddCommand(checkpointNewCmd)
	checkpointCmd.Flags().StringP("filepath", "f", "", "Use with `new`. The path to json configuration file of the CheckPoint")
}

func machineFromCommand(filepath string, machineType string) {
	b := utils.ReadFile(filepath)
	if b == nil {
		return
	}
	var info map[string]interface{}
	err := json.Unmarshal(b, &info)
	if err != nil {
		log.Println(err)
		return
	}

}
