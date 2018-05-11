package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// LastKnownControlPanelState is the last state of the Control Panel so it can be broadcasted to websocket clients.
var LastKnownControlPanelState *ControlPanelMsg

// ControlPanelWsHandler listens on a websocket to messages from the control panel hardware.
func ControlPanelWsHandler(bcast *ControlPanelBroadcaster, opcManager *OpcProcessManager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{} // use default options
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to Ugrade websocket connection\n")
			return
		}
		defer conn.Close()

		// TODO: Only allow one client.

		// Clients needs to reply to Ping.
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error {
			log.Println("OPC WS Pong! ", conn.RemoteAddr())
			conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		defer func() {
			opcManager.SetControlPanelIsOwner(false)
		}()

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

				// Save so we can broadcast to new clients.
				LastKnownControlPanelState = &jsonMsg

				log.Println("Got control panel message: ", jsonMsg)

				// TODO: Map program choices from control panel map to animation sketches in config

				// If the control panel has anything but the "Custom" program selected,
				// it owns the animation selection, so the web clients don't get a say.
				if jsonMsg.Program != int(CustomProgram) {
					log.Println("Control panel owns the animation selection")
					opcManager.SetControlPanelIsOwner(true)
				} else {
					log.Println("Control panel no longer owner of animation selection")
					opcManager.SetControlPanelIsOwner(false)
				}

				// Broadcast so we can show state to web clients.
				bcast.Broadcast(func(c *ControlPanelReceiver) {
					c.controlPanel <- jsonMsg
				})
			}
		}()

		// Writer.
		pingTicker := time.NewTicker(pingPeriod)
		defer pingTicker.Stop()

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
