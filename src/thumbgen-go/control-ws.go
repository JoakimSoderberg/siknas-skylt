package main

import (
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

// Animation represents a single animation.
type Animation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AnimationListMessage represents a list of animations sent by the server
// when the client connects to the control websocket.
type AnimationListMessage struct {
	Anims []Animation `json:"anims"`
}

// ConnectControlWebsocket connects to the OPC Websocket.
func ConnectControlWebsocket(done chan struct{}, interrupt chan os.Signal) (*websocket.Conn, []Animation, chan string, error) {
	url := url.URL{Scheme: "ws", Host: viper.GetString("host"), Path: viper.GetString("ws-path")}

	log.Printf("Connecting to control websocket %v...", url.String())
	ws, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect to websocket server: ", err)
	}

	animationListChan := make(chan []Animation)
	ctrlMessages := make(chan string)

	// Websocket reader.
	go websocketControlReader(ws, animationListChan, interrupt, done)

	// Websocket writer.
	go websocketControlWriter(ws, ctrlMessages, interrupt, done)

	log.Printf("Waiting for Animation list...")
	var animationList []Animation
	select {
	case animationList = <-animationListChan:
	case <-time.After(10 * time.Second):
		log.Fatalln("Control websocket timed out after 10s waiting for the animation list")
	}

	return ws, animationList, ctrlMessages, nil
}

func websocketControlReader(ws *websocket.Conn, animationListChan chan []Animation, interrupt chan os.Signal, done chan struct{}) {

	log.Println("Starting control websocket reader")

	// The first message the server sends is the animation list.
	var animationsMsg AnimationListMessage

	err := ws.ReadJSON(&animationsMsg)
	if err != nil {
		log.Println("Failed to read Animation List message: ", err)
		return
	}

	// Send the Animation list for anyone waiting.
	animationListChan <- animationsMsg.Anims

	for {
		select {
		default:
			// TODO: Read messages and pass over channel here.
		case <-interrupt:
			log.Println("Control Websocket Reader got interrupted...")
			return
		case <-done:
			log.Println("Control Websocket reader done...")
			return
		}
	}
}

func websocketControlWriter(ws *websocket.Conn, ctrlMessages chan string, interrupt chan os.Signal, done chan struct{}) {
	defer ws.Close()

	pingTicker := time.NewTicker(time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case msg := <-ctrlMessages:
			// What to select
			err := ws.WriteJSON(struct {
				MessageType string `json:"message_type"`
				Selected    string `json:"selected"`
			}{
				MessageType: "select",
				Selected:    msg,
			})

			if err != nil {
				log.Println("Control Websocket failed to select Animation over control websocket: ", err)
				return
			}
		case <-done:
			return
		case <-interrupt:
			// User wants to close.
			log.Println("Control Websocket writer got interrupted, attempting clean Websocket close...")

			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := ws.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Failed to write:", err)
				return
			}

			// Wait for the reader to be done or a timeout.
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
