package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/mux"
	"gopkg.in/alecthomas/kingpin.v2"
)

// ControlPanelMsg represents the state of the control panel hardware.
type ControlPanelMsg struct {
	Program    int    `json:"program,omitempty"`
	Color      [3]int `json:"color,omitempty"`
	Brightness int    `json:"brightness,omitempty"`
}

// controlPanelWsListener listens on a websocket to messages from the control panel hardware.
func controlPanelWsListener(controlPanel chan ControlPanelMsg) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
				msg, err := wsutil.ReadServerText(conn)
				if err != nil {
					log.Println("Failed to read control panel message: ", err)
					continue
				}

				jsonMsg := ControlPanelMsg{}

				err = json.Unmarshal(msg, &jsonMsg)
				if err != nil {
					log.Println("Failed to unpack control panel message: ", err)
					continue
				}

				controlPanel <- jsonMsg
			}
		}()
	})
}

func wsListener(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w, nil)
	if err != nil {
		fmt.Printf("Failed to Ugrade websocket connection\n")
		conn.Close()
		return
	}

	go func() {
		defer conn.Close()

		fmt.Printf("Websocket Client connected: %v\n", conn.RemoteAddr())

		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				// handle error
			}

			switch op {
			case ws.OpClose:
				log.Printf("Closing Websocket connection from: %v\n", conn.RemoteAddr())
				return
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

var (
	port = kingpin.Flag("port", "The port the webserver should listen on").Default("8080").Int()
)

func main() {
	kingpin.UsageTemplate(kingpin.DefaultUsageTemplate).Version("1.0").Author("Joakim Soderberg")
	kingpin.CommandLine.Help = "Siknas-skylt Webserver"
	kingpin.Parse()

	controlPanel := make(chan ControlPanelMsg)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/ws", wsListener)
	router.HandleFunc("/ws/control_panel", controlPanelWsListener(controlPanel))

	log.Println("Starting Siknas-skylt webserver...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), router))
}
