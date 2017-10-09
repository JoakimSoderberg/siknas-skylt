package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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

// ControlPanelMsg represents the state of the control panel hardware.
type ControlPanelMsg struct {
	Program    int    `json:"program,omitempty"`
	Color      [3]int `json:"color,omitempty"`
	Brightness int    `json:"brightness,omitempty"`
}

type controlPanelClient struct {
	controlPanel chan ControlPanelMsg
}

type controlPanelBroadcaster struct {
	sync.Mutex
	clients []*controlPanelClient
}

// Push adds a new client as a broadcast listener to the control panel.
func (bcast *controlPanelBroadcaster) Push(c *controlPanelClient) {
	bcast.Lock()
	defer bcast.Unlock()

	bcast.clients = append(bcast.clients, c)

	log.Println("Added control panel broadcast listening client")
}

// Pop removes a client from the control panel broadcast.
func (bcast *controlPanelBroadcaster) Pop(c *controlPanelClient) {
	bcast.Lock()
	defer bcast.Unlock()

	i := -1
	for j, cur := range bcast.clients {
		if cur == c {
			i = j
			break
		}
	}

	if i < 0 {
		return
	}

	// TODO: Keeping clients in a slice might not be the best solution?
	copy(bcast.clients[i:], bcast.clients[i+1:])
	bcast.clients[len(bcast.clients)-1] = nil // or the zero value of T
	bcast.clients = bcast.clients[:len(bcast.clients)-1]

	log.Println("Removed control panel broadcast listening client")
}

// Broadcast will send an incoming message from the control panel to all listening channels.
func (bcast *controlPanelBroadcaster) Broadcast(routine func(*controlPanelClient)) {
	bcast.Lock()
	defer bcast.Unlock()

	// Broadcasts to all clients.
	for _, c := range bcast.clients {
		log.Println("Broadcasting to ", c)
		routine(c)
	}
}

var upgrader = websocket.Upgrader{} // use default options

// controlPanelWsListener listens on a websocket to messages from the control panel hardware.
func controlPanelWsListener(bcast *controlPanelBroadcaster) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

				bcast.Broadcast(func(c *controlPanelClient) {
					c.controlPanel <- jsonMsg
				})
			}
		}()
	})
}

type rawClientMsg struct {
	MessageType int
	Message     []byte
}

type clientMsg struct {
	MessageType string `json:"message_type,omitempty"`
}

type clientSelectMsg struct {
	clientMsg
	Selected string `json:"selected,omitempty"`
}

type serverMsg struct {
	MessageType string `json:"message_type,omitempty"`
}

type serverAnim struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type serverListMsg struct {
	serverMsg
	Anims []serverAnim `json:"anims,omitempty"`
}

type serverStatusMsg struct {
	serverMsg
	Success bool   `json:"success,omitempty"`
	Text    string `json:"text,omitempty"`
}

func getAnimsListMsg() (serverListMsg, error) {
	// TODO: Get a real list of files
	msg := serverListMsg{
		serverMsg: serverMsg{MessageType: "list"},
		Anims:     []serverAnim{{Name: "hej"}, {Name: "hopp"}, {Name: "arne"}},
	}

	return msg, nil
}

// wsListener is the websocket handler for "normal" websocket clients that are not the control panel.
func wsListener(bcast *controlPanelBroadcaster) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to Ugrade websocket connection\n")
			conn.Close()
			return
		}

		// TODO: Does it have to be this generic?
		serverMessages := make(chan interface{})

		// Add this new client as a control panel broadcast listener.
		ctrlPanelClient := controlPanelClient{
			controlPanel: make(chan ControlPanelMsg),
		}
		bcast.Push(&ctrlPanelClient)

		// Reader.
		go func() {
			defer conn.Close()
			defer close(serverMessages)

			// Clients needs to reply to Ping.
			conn.SetReadDeadline(time.Now().Add(pongWait))
			conn.SetPongHandler(func(string) error {
				log.Println("Pong! ", conn.RemoteAddr())
				conn.SetReadDeadline(time.Now().Add(pongWait))
				return nil
			})

			// TODO: Make sure channel is closed
			for {
				_, data, err := conn.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					break
				}

				// TODO: Breakout into function
				var msg clientMsg
				err = json.Unmarshal(data, &msg)
				if err != nil {
					log.Printf("Failed to unmarshal JSON '%v': %v\n", string(data), err)
					break
				}

				switch msg.MessageType {
				default:
					log.Println("Unexpected message type from client: ", msg.MessageType)
					return
				case "select":
					var selectMsg clientSelectMsg
					err := json.Unmarshal(data, &selectMsg)
					if err != nil {
						log.Printf("Failed to unmarshal JSON '%v':\n  %v\n", string(data), err)
						return
					}

					log.Println("Select: ", selectMsg.Selected)
					serverMessages <- serverStatusMsg{
						serverMsg: serverMsg{MessageType: "status"},
						Success:   true,
						Text:      fmt.Sprint("Selected ", selectMsg.Selected),
					}
				}

				log.Printf("recv: %s", msg)
			}
		}()

		// Writer.
		go func() {
			pingTicker := time.NewTicker(pingPeriod)

			defer func() {
				conn.Close()
				bcast.Pop(&ctrlPanelClient)
				pingTicker.Stop()
			}()

			log.Printf("Websocket Client connected: %v\n", conn.RemoteAddr())

			// Start by sending a list of animations.
			anims, err := getAnimsListMsg()
			if err != nil {
				log.Println("Failed to get list of animations")
				return
			}
			conn.WriteJSON(anims)

			for {
				select {
				case msg := <-serverMessages:
					conn.WriteJSON(msg)
				case msg := <-ctrlPanelClient.controlPanel:
					log.Println("Broadcasting: ", msg)
					conn.WriteJSON(msg)
				case <-pingTicker.C:
					log.Println("Ping! ", conn.RemoteAddr())
					conn.SetWriteDeadline(time.Now().Add(writeWait))
					if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
						return
					}
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
	//controlPanel := make(chan ControlPanelMsg)
	controlPanel := &controlPanelBroadcaster{
		clients: make([]*controlPanelClient, 0, 1),
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/panel", controlPanelClientHandler)

	// Websocket handlers.
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
