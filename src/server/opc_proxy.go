package main

//go:generate go run broadcaster/gen.go Opc broadcaster/broadcast.tmpl

import (
	"log"
	"net"

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

	log.Println("Starting Open Pixel Control proxy server...")

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

		// TODO: Keep track of current OPC client, once it disconnects fade out

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

// handleOpcCon handles connections from OPC clients.
func handleOpcCon(messages chan *opc.Message, conn net.Conn) {
	defer func() {
		log.Println("[OPC incoming client] Disconnected:", conn.RemoteAddr())
		conn.Close()
	}()

	for {
		msg, err := opc.ReadOpc(conn)
		if err != nil {
			log.Println("[OPC incoming client] Failed to read OPC:", conn.RemoteAddr())
			return
		}

		// TODO: When a new OPC client connects, fade in the brightness.
		// TODO: Make it possible to ignore incoming messages (don't forward them)

		// TODO: Create a standalone color correction Go-server that can be started and listens to
		// websocket connections. The webpage can then connect to that and perform color
		// correction. Like this:
		// https://github.com/scanlime/fadecandy/blob/686ab1f5570e563a287474424565bfbf8d8fe4a8/examples/python/color-correction-ui.py#L19-L30

		messages <- msg
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
