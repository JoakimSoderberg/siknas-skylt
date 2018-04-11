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
func ConnectOPCWebsocket(done chan struct{}, interrupt chan os.Signal, filename string) {
	url := url.URL{Scheme: "ws", Host: viper.GetString("host"), Path: viper.GetString("ws-opc-path")}
	captureDuration := viper.GetDuration("capture-duration")

	ws := connectWebsocket(url.String())
	defer ws.Close()

	// Websocket reader.
	go websocketReader(ws, interrupt, done)

	// Websocket writer.
	go websocketWriter(ws, interrupt, done)

	// TODO: Fix closing down nicely.
	time.Sleep(captureDuration)
	log.Printf("Finished capturing after %s\n", captureDuration)
	close(done)

	createOutputSVG(filename)
}

func connectWebsocket(addr string) *websocket.Conn {
	log.Printf("Connecting to OPC websocket %s...", addr)

	ws, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Fatal("Failed to connect to websocket server: ", err)
	}

	return ws
}

func websocketReader(ws *websocket.Conn, interrupt chan os.Signal, done chan struct{}) {
	var opcMsg OpcMessage

	log.Println("Starting websocket reader")

	started := false

	for {
		select {
		default:
			messageType, messageData, err := ws.ReadMessage()
			if err != nil {
				log.Println("Failed to read: ", err)
				return
			}

			if !started {
				started = true
				log.Printf("Started capturing %s of animation\n", viper.GetDuration("capture-duration"))
			}

			if messageType != websocket.BinaryMessage {
				log.Println("ERROR: Got a Text message on the OPC Websocket, expected Binary")
				break
			}

			buf := bytes.NewBuffer(messageData[0:binary.Size(opcMsg.Header)])
			err = binary.Read(buf, binary.BigEndian, &opcMsg.Header)
			if err != nil {
				log.Println("ERROR: Failed to read OPC message: ", err)
				break
			}

			realMsgLength := uint16(len(messageData) - binary.Size(opcMsg.Header))

			if opcMsg.Header.Length != realMsgLength {
				log.Printf("ERROR: Got a %d byte invalid OPC message. Header says %d, got %d bytes\n", opcMsg.Header.Length, opcMsg.Header.Length, realMsgLength)
				break
			}

			// Note we don't really need the OPC Length here, since this is Websockets
			// and we already have a known message length.
			opcMsg.Data = messageData[binary.Size(opcMsg.Header):]

			opcMessages = append(opcMessages, opcMsg)

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

func websocketWriter(ws *websocket.Conn, interrupt chan os.Signal, done chan struct{}) {
	defer ws.Close()

	pingTicker := time.NewTicker(time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-done:
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
