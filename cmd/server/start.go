/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package server

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start server",
	Long:  `Start server. Default RPC port is 8813`,
	Run: func(cmd *cobra.Command, args []string) {
		rpc, _ := cmd.Flags().GetInt16("rpc")
		r := gin.Default()
		// Set up CORS middleware
		config := cors.DefaultConfig()
		config.AllowAllOrigins = true
		r.Use(cors.New(config))
		r.Run(":" + strconv.Itoa(int(rpc)))
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().Int16P("rpc", "p", 8813, "Set RPC port to connect to client. Default port is 8813")
}
