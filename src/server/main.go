package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gopkg.in/alecthomas/kingpin.v2"
)

// ControlPanelMsg represents the state of the control panel hardware.
type ControlPanelMsg struct {
	Program    int    `json:"program,omitempty"`
	Color      [3]int `json:"color,omitempty"`
	Brightness int    `json:"brightness,omitempty"`
}

var upgrader = websocket.Upgrader{} // use default options

// controlPanelWsListener listens on a websocket to messages from the control panel hardware.
func controlPanelWsListener(controlPanel chan ControlPanelMsg) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to Ugrade websocket connection\n")
			conn.Close()
			return
		}

		// TODO: Only allow one client.

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

				controlPanel <- jsonMsg
			}
		}()
	})
}

type clientMsg struct {
	MessageType int
	Message     []byte
}

func wsListener(controlPanel chan ControlPanelMsg) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to Ugrade websocket connection\n")
			conn.Close()
			return
		}

		clientMessages := make(chan clientMsg)

		// Reader.
		go func() {
			defer conn.Close()

			for {
				mt, message, err := conn.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					break
				}
				log.Printf("recv: %s", message)
				clientMessages <- clientMsg{mt, message}
			}
		}()

		// Writer.
		go func() {
			defer conn.Close()

			log.Printf("Websocket Client connected: %v\n", conn.RemoteAddr())

			// Writer
			for {
				select {
				case msg := <-clientMessages:
					conn.WriteMessage(msg.MessageType, msg.Message)
				case msg := <-controlPanel:
					// TODO: This only sends it to one client...
					log.Println("Broadcasting: ", msg)
					conn.WriteJSON(msg)
				}
			}
		}()
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	// TODO: Host aurelia webpage
	homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
}

func controlPanelClient(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	// TODO: Host aurelia webpage
	homeTemplate.Execute(w, "ws://"+r.Host+"/ws/control_panel")
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
	router.HandleFunc("/panel", controlPanelClient)
	router.HandleFunc("/ws", wsListener(controlPanel))
	router.HandleFunc("/ws/control_panel", controlPanelWsListener(controlPanel))

	log.Println("Starting Siknas-skylt webserver...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), router))
}

var homeTemplate = template.Must(template.New("").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
	<meta charset="utf-8">
	<script>  
	window.addEventListener("load", function(evt) {
		var output = document.getElementById("output");
		var input = document.getElementById("input");
		var ws;
		var print = function(message) {
			var d = document.createElement("div");
			d.innerHTML = message;
			output.appendChild(d);
		};
		document.getElementById("open").onclick = function(evt) {
			if (ws) {
				return false;
			}
			ws = new WebSocket("{{.}}");
			ws.onopen = function(evt) {
				print("OPEN");
			}
			ws.onclose = function(evt) {
				print("CLOSE");
				ws = null;
			}
			ws.onmessage = function(evt) {
				print("RESPONSE: " + evt.data);
			}
			ws.onerror = function(evt) {
				print("ERROR: " + evt.data);
			}
			return false;
		};
		document.getElementById("send").onclick = function(evt) {
			if (!ws) {
				return false;
			}
			print("SEND: " + input.value);
			ws.send(input.value);
			return false;
		};
		document.getElementById("close").onclick = function(evt) {
			if (!ws) {
				return false;
			}
			ws.close();
			return false;
		};
	});
	</script>
	</head>
	<body>
	<table>
	<tr><td valign="top" width="50%">
	<p>Click "Open" to create a connection to the server, 
	"Send" to send a message to the server and "Close" to close the connection. 
	You can change the message and send multiple times.
	<p>
	<form>
	<button id="open">Open</button>
	<button id="close">Close</button>
	<p><input id="input" type="text" value="Hello world!">
	<button id="send">Send</button>
	</form>
	</td><td valign="top" width="50%">
	<div id="output"></div>
	</td></tr></table>
	</body>
	</html>
	`))
