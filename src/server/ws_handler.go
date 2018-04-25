package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// clientMsg is a message is a message sent by clients.
type clientMsg struct {
	MessageType string `json:"message_type"`
}

// clientSelectMsg is sent by clients when selecting an animation
// from the list of available animations returned in servListMsg.
type clientSelectMsg struct {
	clientMsg
	Selected string `json:"selected"`
}

// serverMsg is a message returned over the websocket.
type serverMsg struct {
	MessageType string `json:"message_type"`
}

type serverControlPanelMsg struct {
	serverMsg
	Data ControlPanelMsg `json"data"`
}

// serverAnim represents a processing animation sketch.
type serverAnim struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// servListMsg is a message containing a list of processing Animations available to choose from.
type serverListMsg struct {
	serverMsg
	Anims []serverAnim `json:"anims,omitempty"`
}

// serverStatusMsg is a status message for any action a client performed.
type serverStatusMsg struct {
	serverMsg
	Success bool   `json:"success,omitempty"`
	Text    string `json:"text,omitempty"`
}

// getAnimsListMsg returns a list of available animation processes.
func getAnimsListMsg(opcManager *OpcProcessManager) (serverListMsg, error) {
	msg := serverListMsg{
		serverMsg: serverMsg{MessageType: "list"},
	}

	msg.Anims = make([]serverAnim, len(opcManager.Processes))
	i := 0
	for name, val := range opcManager.Processes {
		msg.Anims[i].Name = name
		msg.Anims[i].Description = val.Description
	}

	return msg, nil
}

// Unmarshals a client message and returns a server status.
func unmarshalClientMsg(data []byte, opcManager *OpcProcessManager) (string, error) {
	var msg clientMsg
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON '%v': %v", string(data), err)
	}

	log.Printf("Got Websocket client message:\n%v\n", msg)

	switch msg.MessageType {
	default:
		return "", fmt.Errorf("unexpected message type from client: %v", msg.MessageType)
	case "select":
		// TODO: Make into function.
		var selectMsg clientSelectMsg
		err := json.Unmarshal(data, &selectMsg) // TODO: Does not fail if "selected" is missing
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal JSON '%v':\n  %v", string(data), err)
		}

		if opcManager.IsControlPanelOwner() {
			log.Printf("Control panel owns animation selection, ignoring client request")
			return "", fmt.Errorf("The control panel owns animation selection")
		}

		// TODO: Selecting a new sketch should broadcast the selection to all ws clients

		if err := opcManager.StartAnim(selectMsg.Selected); err != nil {
			return "", err
		}

		return fmt.Sprintf("Started animation %v", selectMsg.Selected), nil
	}
}

// WsHandler is the websocket handler for "normal" websocket clients that are not the control panel.
func WsHandler(bcast *ControlPanelBroadcaster, opcManager *OpcProcessManager) http.HandlerFunc {
	// TODO: Break out go functions
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{} // use default options
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to Ugrade websocket connection\n")
			conn.Close()
			return
		}

		serverMessages := make(chan interface{})

		// Add this new client as a control panel broadcast listener.
		ctrlPanelClient := NewControlPanelReceiver()
		bcast.Push(ctrlPanelClient)

		// Clients needs to reply to Ping.
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error {
			log.Println("Pong! ", conn.RemoteAddr())
			conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		// Reader.
		go readOpcWsConn(conn, serverMessages, ctrlPanelClient, opcManager)

		// Writer.
		go func() {
			pingTicker := time.NewTicker(pingPeriod)

			defer func() {
				conn.Close()
				bcast.Pop(ctrlPanelClient)
				pingTicker.Stop()
			}()

			log.Printf("Websocket Client connected: %v\n", conn.RemoteAddr())

			// Start by sending a list of animations.
			anims, err := getAnimsListMsg(opcManager)
			if err != nil {
				log.Println("Failed to get list of animations")
				return
			}
			conn.WriteJSON(anims)

			for {
				select {
				case msg := <-serverMessages:
					err := conn.WriteJSON(msg)
					if err != nil {
						log.Printf("Failed to write to websocket client %v: %v\n", conn.RemoteAddr(), err)
						return
					}
				case msg := <-ctrlPanelClient.controlPanel:
					log.Println("Broadcasting control panel message: ", msg)
					serverControlPanelMsg := serverControlPanelMsg{
						serverMsg: serverMsg{
							MessageType: "control_panel",
						},
						Data: msg,
					}
					conn.WriteJSON(serverControlPanelMsg)
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

// readOpcWsConn reads incoming Websocket client messages.
func readOpcWsConn(conn *websocket.Conn, serverMessages chan interface{},
	ctrlPanelClient *ControlPanelReceiver, opcManager *OpcProcessManager) {

	defer func() {
		conn.Close()
		close(serverMessages)
		close(ctrlPanelClient.controlPanel)
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		statusText, err := unmarshalClientMsg(data, opcManager)
		if err != nil {
			statusText = err.Error()
		}

		serverMessages <- serverStatusMsg{
			serverMsg: serverMsg{MessageType: "status"},
			Success:   (err != nil),
			Text:      statusText,
		}
	}
}
