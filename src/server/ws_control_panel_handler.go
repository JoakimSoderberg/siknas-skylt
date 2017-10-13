package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ControlPanelWsListener listens on a websocket to messages from the control panel hardware.
func ControlPanelWsListener(bcast *ControlPanelBroadcaster) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{} // use default options
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to Ugrade websocket connection\n")
			conn.Close()
			return
		}

		// TODO: Only allow one client.

		// Reader.
		go func() {
			defer conn.Close()

			log.Printf("Websocket Control Panel Client connected: %v\n", conn.RemoteAddr())

			for {
				jsonMsg := ControlPanelMsg{}

				err := conn.ReadJSON(&jsonMsg)
				if err != nil {
					log.Println("Failed to read control panel message: ", err)
					continue
				}

				log.Println("Got control panel message: ", jsonMsg)

				bcast.Broadcast(func(c *ControlPanelReceiver) {
					c.controlPanel <- jsonMsg
				})
			}
		}()
	})
}
