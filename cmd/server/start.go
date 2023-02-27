/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hectagon-finance/chain-mvp/types"
	"github.com/spf13/cobra"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start server",
	Long:  `Start server. Default API port is 8813`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt16("port")
		types.StartConsensus()
		go handleEvents()
		startListen(port)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().Int16P("port", "p", 8813, "Set API port to connect to client. Default port is 8813")
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleEvents() {
	for {
		ev := <-types.Broadcast
		fmt.Println("ev: ", ev)
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(ev)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
func startListen(port int16) {
	r := gin.Default()
	/**
	* Return the current Block
	 */
	r.GET("/block", func(c *gin.Context) {
		blockNo := types.GetCurrentBlockNumber()
		c.JSON(http.StatusOK, gin.H{
			"currentBlockNo": blockNo,
		})

		// HACK: mock event
		v, _ := json.Marshal(blockNo)
		_, id := types.CreateEvent("BlockCalled", v)
		go types.Emit(id)
	})
	r.GET("/ws", func(c *gin.Context) {
		// Upgrade initial GET request to a websocket
		// Configure the upgrader
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
		}
		// Make sure we close connection when the function return
		// defer ws.Close()
		clients[ws] = true
		fmt.Println("A client is registered")
	})
	// Set up CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.Run(":" + strconv.Itoa(int(port)))
}
