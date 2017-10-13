package main

import (
	"fmt"
	"html/template"
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
	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	// TODO: Host aurelia webpage
	homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
}

func controlPanelClientHandler(w http.ResponseWriter, r *http.Request) {
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

	// Broadcast channel for control panel.
	// TODO: Create a new function instead
	controlPanel := &ControlPanelBroadcaster{
		clients: make([]*ControlPanelReceiver, 0, 1),
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/panel", controlPanelClientHandler)

	// Websocket handlers.
	router.HandleFunc("/ws", WsListener(controlPanel))
	router.HandleFunc("/ws/control_panel", ControlPanelWsListener(controlPanel))

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
