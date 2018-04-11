package main

import (
	"fmt"
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
func ConnectControlWebsocket(done chan struct{}, interrupt chan os.Signal) ([]Animation, error) {
	url := url.URL{Scheme: "ws", Host: viper.GetString("host"), Path: viper.GetString("ws-path")}

	ws := connectControlWebsocket(url.String())
	defer ws.Close()

	animationListChan := make(chan []Animation)
	ctrlMessages := make(chan string)

	// Websocket reader.
	go websocketControlReader(ws, animationListChan, interrupt, done)

	// Websocket writer.
	go websocketControlWriter(ws, ctrlMessages, interrupt, done)

	log.Printf("Waiting for Animation list...")
	var animationList []Animation
	for {
		select {
		case animationList = <-animationListChan:
			return animationList, nil
		case <-time.After(10 * time.Second):
			return nil, fmt.Errorf("timed out after 10s waiting for the animation list")
		}
	}
}

func connectControlWebsocket(addr string) *websocket.Conn {
	log.Printf("Connecting to control websocket %s...", addr)

	ws, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Fatal("Failed to connect to websocket server: ", err)
	}

	return ws
}

func websocketControlReader(ws *websocket.Conn, animationListChan chan []Animation, interrupt chan os.Signal, done chan struct{}) {

	log.Println("Starting control websocket reader")

	var receivedAnimList = false

	for {
		select {
		default:
			if !receivedAnimList {
				var animationsMsg AnimationListMessage

				err := ws.ReadJSON(&animationsMsg)
				if err != nil {
					log.Println("Failed to read Animation List message: ", err)
					return
				}

				// Send the Animation list for anyone waiting.
				animationListChan <- animationsMsg.Anims
			}
		case <-interrupt:
			// TODO: Fix this
			log.Println("Reader got interrupted...")
			return
		case <-done:
			log.Println("Reader done...")
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
				Selected    string `json:"name"`
			}{
				MessageType: "select",
				Selected:    msg,
			})

			if err != nil {
				log.Println("Failed to select Animation over control websocket: ", err)
				return
			}
		case <-done:
			return
		case <-interrupt:
			// User wants to close.
			log.Println("Writer got interrupted, attempting clean Websocket close...")

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
