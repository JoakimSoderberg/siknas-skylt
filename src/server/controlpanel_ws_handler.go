package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// ControlPanelWsHandler listens on a websocket to messages from the control panel hardware.
func ControlPanelWsHandler(bcast *ControlPanelBroadcaster) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{} // use default options
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to Ugrade websocket connection\n")
			return
		}
		defer conn.Close()

		// TODO: Only allow one client.
		// TODO: When this client is connected ignore selection via webpage if not "custom program" is selected.
		// TODO: Enable connectin

		// Clients needs to reply to Ping.
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error {
			log.Println("OPC WS Pong! ", conn.RemoteAddr())
			conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		// Reader.
		go func() {
			log.Printf("Websocket Control Panel Client connected: %v\n", conn.RemoteAddr())

			for {
				jsonMsg := ControlPanelMsg{}

				err := conn.ReadJSON(&jsonMsg)
				if err != nil {
					log.Println("Failed to read control panel message: ", err)
					return
				}

				log.Println("Got control panel message: ", jsonMsg)

				bcast.Broadcast(func(c *ControlPanelReceiver) {
					c.controlPanel <- jsonMsg
				})
			}
		}()

		// Writer.
		pingTicker := time.NewTicker(pingPeriod)

		defer func() {
			pingTicker.Stop()
		}()

		for {
			select {
			case <-pingTicker.C:
				log.Println("Control Panel Ping!", conn.RemoteAddr())
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return
				}
			}
		}
	})
}
