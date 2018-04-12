package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

// OpcMessageHeader defines the header of the Open Pixel Control (OPC) Protocol.
type OpcMessageHeader struct {
	Channel byte
	Command byte
	Length  uint16
}

// OpcMessage defines a OPC message including header.
type OpcMessage struct {
	Header OpcMessageHeader
	Data   []byte
}

// RGB returns the Red, Green and Blue values between 0-255 for a given pixel.
func (m *OpcMessage) RGB(ledIndex int) (uint8, uint8, uint8) {
	i := 3 * ledIndex
	return m.Data[i], m.Data[i+1], m.Data[i+2]
}

var opcMessages []OpcMessage

// ConnectOPCWebsocket connects to the OPC Websocket.
func ConnectOPCWebsocket(done chan struct{}, interrupt chan os.Signal) *websocket.Conn {
	url := url.URL{Scheme: "ws", Host: viper.GetString("host"), Path: viper.GetString("ws-opc-path")}

	log.Printf("Connecting to OPC websocket %v...", url)
	ws := connectWebsocket(url.String())

	// Websocket reader.
	go websocketReader(ws, interrupt, done)

	// Websocket writer.
	go websocketWriter(ws, interrupt, done)

	return ws
}

func connectWebsocket(addr string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Fatalln("Failed to connect to websocket server: ", err)
	}

	return ws
}

func websocketReader(ws *websocket.Conn, interrupt chan os.Signal, done chan struct{}) {
	var opcMsg OpcMessage

	log.Println("Starting OPC websocket reader")

	started := false

	for {
		select {
		default:
			messageType, messageData, err := ws.ReadMessage()
			if err != nil {
				log.Println("OPC Websocket failed to read: ", err)
				return
			}

			if !started {
				started = true
				log.Printf("OPC Websocket started capturing %s of animation\n", viper.GetDuration("capture-duration"))
			}

			if messageType != websocket.BinaryMessage {
				log.Fatalln("ERROR: Got a Text message on the OPC Websocket, expected Binary")
			}

			buf := bytes.NewBuffer(messageData[0:binary.Size(opcMsg.Header)])
			err = binary.Read(buf, binary.BigEndian, &opcMsg.Header)
			if err != nil {
				log.Fatalln("OPC Websocket Failed to read OPC message: ", err)
			}

			realMsgLength := uint16(len(messageData) - binary.Size(opcMsg.Header))

			if opcMsg.Header.Length != realMsgLength {
				log.Fatalf("ERROR: Got a %d byte invalid OPC message. Header says %d, got %d bytes\n", opcMsg.Header.Length, opcMsg.Header.Length, realMsgLength)
			}

			// Note we don't really need the OPC Length here, since this is Websockets
			// and we already have a known message length.
			opcMsg.Data = messageData[binary.Size(opcMsg.Header):]

			// TODO: Only start appending once we are signaled to.
			opcMessages = append(opcMessages, opcMsg)

		case <-interrupt:
			// TODO: Fix this
			log.Println("OPC Websocket Reader got interrupted...")
			return
		case <-done:
			log.Println("OPC Websocket Reader done...")
			return
		}
	}
}

func websocketWriter(ws *websocket.Conn, interrupt chan os.Signal, done chan struct{}) {
	pingTicker := time.NewTicker(time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-done:
		case <-interrupt:
			// User wants to close.
			log.Println("OPC Websocket Writer got interrupted, attempting clean Websocket close...")

			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := ws.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("OPC Websocket Failed to write:", err)
				return
			}

			// Wait for the reader to be done or a timeout.
			select {
			case <-done:
				log.Println("OPC Websocket done!")
			case <-time.After(time.Second):
				log.Println("OPC Websocket timed out")
			}
			return
		}
	}
}
