package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// OpcWsHandler handles websockets clients that wants to listen to OPC messages
// passed on by the OPC proxy.
func OpcWsHandler(bcast *OpcBroadcaster) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{} // use default options
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to Ugrade websocket connection\n")
			return
		}

		log.Println("OPC WS Client connected: ", conn.RemoteAddr())

		// Start listening to the OPC broadcasts.
		// TODO: Race condition, use a channel here
		opcReceiver := NewOpcReceiver()
		bcast.Push(opcReceiver)

		// Clients needs to reply to Ping.
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error {
			log.Println("OPC WS Pong! ", conn.RemoteAddr())
			conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		go func() {
			pingTicker := time.NewTicker(pingPeriod)

			defer func() {
				conn.Close()
				bcast.Pop(opcReceiver)
				pingTicker.Stop()
				close(opcReceiver.opcMessages)
			}()

			//timeSinceLastMsg := time.Now()

			for {
				select {
				case msg := <-opcReceiver.opcMessages:
					// TODO: Make a setting to tweak this throttling.
					/*if time.Since(timeSinceLastMsg) < 100*time.Millisecond {
						continue
					}
					timeSinceLastMsg = time.Now()*/

					conn.SetWriteDeadline(time.Now().Add(writeWait))
					err := conn.WriteMessage(websocket.BinaryMessage, msg.ByteArray())
					if err != nil {
						log.Printf("Failed to write to websocket client %v: %v\n", conn.RemoteAddr(), err)
						return
					}
				case <-pingTicker.C:
					log.Println("OPC WS Ping!", conn.RemoteAddr())
					conn.SetWriteDeadline(time.Now().Add(writeWait))
					if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
						return
					}
				}
			}
		}()
	})
}
