package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/mux"
)

func control_panel_ws_listener(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w, nil)
	if err != nil {
		fmt.Printf("Failed to Ugrade websocket connection\n")
		conn.Close()
		return
	}

	// TODO: Only allow one client.

	go func() {
		defer conn.Close()

		fmt.Printf("New Websocket Client\n")

		for {
			_, _, err := wsutil.ReadClientData(conn)
			if err != nil {
				// handle error
			}
			// TODO: Broadcast control panel messages to ws_listener
			/*
				err = wsutil.WriteServerMessage(conn, op, msg)
				if err != nil {
					// handle error
				}*/
		}
	}()
}

func ws_listener(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w, nil)
	if err != nil {
		fmt.Printf("Failed to Ugrade websocket connection\n")
		conn.Close()
		return
	}

	go func() {
		defer conn.Close()

		fmt.Printf("New Websocket Client\n")

		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				// handle error
			}
			err = wsutil.WriteServerMessage(conn, op, msg)
			if err != nil {
				// handle error
			}
		}
	}()
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	// TODO: Host aurelia webpage
}

func main() {
	fmt.Printf("Siknas-skylt Webserver\n\n")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/ws", ws_listener)
	router.HandleFunc("/ws/control_panel", control_panel_ws_listener)

	log.Fatal(http.ListenAndServe(":8080", router))
}
