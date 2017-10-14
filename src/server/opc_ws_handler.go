package main

import (
	"log"
	"net/http"

	"github.com/kellydunn/go-opc"

	"github.com/gorilla/websocket"
)

func OpcWsHandler(bcast *OpcBroadcaster) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{} // use default options
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to Ugrade websocket connection\n")
			conn.Close()
			return
		}

		// Start listening to the OPC broadcasts.
		opcReceiver := OpcReceiver{
			opcMessages: make(chan *opc.Message),
		}
		bcast.Push(&opcReceiver)

		go func() {
			defer bcast.Pop(&opcReceiver)
		}()
	})
}
