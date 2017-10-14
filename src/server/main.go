package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	// TODO: Host aurelia webpage
}

var (
	port = kingpin.Flag("port", "The port the webserver should listen on").Default("8080").Int()
)

func main() {
	kingpin.UsageTemplate(kingpin.DefaultUsageTemplate).Version("1.0").Author("Joakim Soderberg")
	kingpin.CommandLine.Help = "Siknas-skylt Webserver"
	kingpin.Parse()

	// Broadcast channel for control panel.
	controlPanelBroadcaster := NewControlPanelBroadcaster()
	opcBroadcaster := NewOpcBroadcaster()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/debug", WsDebugHandler)
	router.HandleFunc("/debug/control_panel", WsDebugControlPanelHandler)

	// Websocket handlers.
	router.HandleFunc("/ws", WsHandler(controlPanelBroadcaster))
	router.HandleFunc("/ws/opc", OpcWsHandler(opcBroadcaster))
	router.HandleFunc("/ws/control_panel", ControlPanelWsHandler(controlPanelBroadcaster))

	log.Println("Starting Siknas-skylt webserver...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), router))
}
