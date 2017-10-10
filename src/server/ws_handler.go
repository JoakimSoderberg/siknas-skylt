package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

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

// Unmarshals a client message and returns a server status.
func unmarshalClientMsg(data []byte) (string, error) {
	var msg clientMsg
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON '%v': %v", string(data), err)
	}

	switch msg.MessageType {
	default:
		return "", fmt.Errorf("unexpected message type from client: %v", msg.MessageType)
	case "select":
		var selectMsg clientSelectMsg
		err := json.Unmarshal(data, &selectMsg)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal JSON '%v':\n  %v", string(data), err)
		}

		// TODO: Selecting a new sketch should broadcast the selection to all ws clients

		log.Println("Select: ", selectMsg.Selected)
		// TODO: Close any existing process
		// TODO: Start the selected process
		return fmt.Sprint("Selected ", selectMsg.Selected), nil
	}
}

// WsListener is the websocket handler for "normal" websocket clients that are not the control panel.
func WsListener(bcast *ControlPanelBroadcaster) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{} // use default options
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to Ugrade websocket connection\n")
			conn.Close()
			return
		}

		// TODO: Does it have to be this generic?
		serverMessages := make(chan interface{})

		// Add this new client as a control panel broadcast listener.
		ctrlPanelClient := ControlPanelClient{
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

				statusText, err := unmarshalClientMsg(data)
				if err != nil {
					statusText = err.Error()
				}

				serverMessages <- serverStatusMsg{
					serverMsg: serverMsg{MessageType: "status"},
					Success:   (err != nil),
					Text:      statusText,
				}
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
