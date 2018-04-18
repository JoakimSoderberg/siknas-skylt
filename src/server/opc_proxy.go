package main

//go:generate go run broadcaster/gen.go Opc broadcaster/broadcast.tmpl

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	// TODO: Replace with own version since lib is abandoned...
	"github.com/kellydunn/go-opc"
)

// OpcReceiver is the context used by the OpcBroadcaster, it contains the channel
// that the OPC messages is broadcasted to.
type OpcReceiver struct {
	opcMessages chan *opc.Message
}

// NewOpcReceiver creates a new OpcReceiver.
func NewOpcReceiver() *OpcReceiver {
	return &OpcReceiver{
		opcMessages: make(chan *opc.Message),
	}
}

// OpcSink is the interface that the OPC proxy uses to write an OPC message to some end point.
type OpcSink interface {
	Write(msg *opc.Message) error
}

// Write will broadcast the OPC message to all listeners.
func (b *OpcBroadcaster) Write(msg *opc.Message) error {
	b.Broadcast(func(c *OpcReceiver) {
		select {
		case c.opcMessages <- msg:
		default:
			// We should not block if a client stops receiving.
		}
	})

	return nil
}

// RunOPCProxy runs an Open Pixel Protocol proxy server that writes any incoming
// OPC messages on the listening port to a set of OpcSinks.
func RunOPCProxy(protocol string, port string, sink OpcSink) error {

	log.Printf("Starting Open Pixel Control proxy server on port %s...", port)

	// Channel used to pass on incoming OPC messages.
	messages := make(chan *opc.Message)

	// TODO: when no client is sending us messages we should send our own.
	// We can save the last message and fade out from that state to full black
	// Then on a new client connection we should fade in instead.

	// Listen for incoming OPC clients.
	go func() {
		listener, err := net.Listen(protocol, port)
		if err != nil {
			log.Fatalln("Failed to start OPC server: ", err)
		}

		log.Println("Listening for OPC connections...")

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Failed to accept client: ", err)
				continue
			}

			log.Println("[OPC incoming client] Connected:", conn.RemoteAddr())

			// Reads from the OPC messages into the channel.
			handleOpcCon(messages, conn)
		}
	}()

	// Process the OPC messages.
	go processOpc(messages, sink)

	return nil
}

// TODO: Move to separate file

// FadecandyColorCorrectionMsg represents a color correction message for the Fadecandy LED controller board.
type FadecandyColorCorrectionMsg struct {
	Gamma      float32   `json:"gamma"`
	Whitepoint []float32 `json:"whitepoint"`
}

func createFadecandyColorCorrectionPacket(gamma, red, green, blue float32) (*opc.Message, error) {
	msg := opc.NewMessage(0)

	contentMsg := FadecandyColorCorrectionMsg{Gamma: gamma, Whitepoint: []float32{red, green, blue}}
	contentBytes, err := json.Marshal(contentMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal color correction message: ", err)
	}

	data := []byte{0x00, 0x01} // Command ID for color correction.
	data = append(data, contentBytes...)

	msg.SystemExclusive(
		[]byte{0x00, 0x01}, // System ID for Fadecandy board.
		data)
	msg.SetLength(uint16(len(data) + 2)) // Include System ID 2 bytes

	return msg, nil
}

// handleOpcCon handles connections from OPC clients.
func handleOpcCon(messages chan *opc.Message, conn net.Conn) {
	defer func() {
		log.Println("[OPC incoming client] Disconnected:", conn.RemoteAddr())
		conn.Close()
	}()

	// TODO: Take into account what the current max brightness is set to.
	// We will ramp it up by interjecting color correction packets that raises the brightness gradually.
	brightness := float32(0.0)
	fadeTimeout := 3000
	fadeDoneTimer := time.NewTimer(time.Duration(fadeTimeout) * time.Millisecond)
	fadeTicker := time.NewTicker(time.Duration(fadeTimeout) * time.Millisecond / 1000)

	// TODO: Break out into separet function.
fadeLoop:
	for {
		msg, err := opc.ReadOpc(conn)
		if err != nil {
			log.Println("[OPC incoming client] Failed to read OPC:", conn.RemoteAddr())
			return
		}

		select {
		case <-fadeTicker.C:
			// Fade in the brightness interleaved.
			colorCorrMsg, err := createFadecandyColorCorrectionPacket(float32(2.5), brightness, brightness, brightness)
			if err != nil {
				log.Println("[OPC incoming client]: ", err)
				return
			}
			messages <- colorCorrMsg
			brightness += float32(1.0 / fadeTimeout)

		case <-fadeDoneTimer.C:
			// Fade completed (make sure we are att full brightness).
			colorCorrMsg, err := createFadecandyColorCorrectionPacket(float32(2.5), 1.0, 1.0, 1.0)
			if err != nil {
				log.Println("[OPC incoming client]: ", err)
				return
			}
			messages <- colorCorrMsg
			break fadeLoop

		default:
			// Send the regular OPC message.
			messages <- msg
		}
	}

	// Normal parsing of packets.
	for {
		msg, err := opc.ReadOpc(conn)
		if err != nil {
			log.Println("[OPC incoming client] Failed to read OPC:", conn.RemoteAddr())
			return
		}

		messages <- msg

		// TODO: Example of global brightness control in C++ https://github.com/scanlime/fadecandy/blob/9d09c62a4ce83f12d1f2aca0429c1cb18b9ec28a/examples/cpp/lib/brightness.h

		// TODO: Color correction code can be used for brightness:
		// https://github.com/scanlime/fadecandy/blob/3bace3a766e96bea0d23bf7846d26a230eee118b/examples/processing/grid24x8z_waves/OPC.pde#L153-L223
		// Example how it is used to set brightness:
		// https://github.com/scanlime/fadecandy/blob/3bace3a766e96bea0d23bf7846d26a230eee118b/examples/processing/grid24x8z_waves/grid24x8z_waves.pde#L38-L40

		// TODO: Make it possible to ignore incoming messages (don't forward them)

		// TODO: Create a standalone color correction Go-server that can be started and listens to
		// websocket connections. The webpage can then connect to that and perform color
		// correction. Like this:
		// https://github.com/scanlime/fadecandy/blob/686ab1f5570e563a287474424565bfbf8d8fe4a8/examples/python/color-correction-ui.py#L19-L30
	}
}

// processOpc receives the incoming OPC messages and dispatches them
func processOpc(messages chan *opc.Message, sink OpcSink) {
	for {
		msg := <-messages

		// This will broadcast the messages.
		sink.Write(msg)
	}
}
