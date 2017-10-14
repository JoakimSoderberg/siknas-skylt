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

// OpcSink is the interface that the OPC proxy uses to write an OPC message to some end point.
type OpcSink interface {
	Write(msg *opc.Message) error
}

// OpcBroadcastSink implements the OpcSink interface to broadcast all OPC messages on an OpcBroadcaster.
type OpcBroadcastSink struct {
	broadcaster *OpcBroadcaster
}

// Write broadcasts the OPC messages to all connected broadcast receivers.
func (o *OpcBroadcastSink) Write(msg *opc.Message) error {
	o.broadcaster.Broadcast(func(c *OpcReceiver) {
		c.opcMessages <- msg
	})

	return nil
}

// RunOPCProxy runs an Open Pixel Protocol proxy server that writes any incoming
// OPC messages on the listening port to a set of OpcSinks.
func RunOPCProxy(protocol string, port string, sinks []OpcSink) error {

	// Channel used to pass on incoming OPC messages.
	messages := make(chan *opc.Message)

	// Reads OPC messages.
	go func() {
		listener, err := net.Listen(protocol, port)
		if err != nil {
			log.Fatalln("Failed to start OPC server: ", err)
		}

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Failed to accept client: ", err)
				continue
			}

			go handleOpcCon(messages, conn)
		}
	}()

	// Process the OPC messages.
	go processOpc(messages, sinks)

	return nil
}

// handleOpcCon handles connections from OPC clients.
func handleOpcCon(messages chan *opc.Message, conn net.Conn) {
	defer conn.Close()
	defer close(messages)

	for {
		msg, err := opc.ReadOpc(conn)
		if err != nil {
			// If we encounter an error reading from the connection,
			// "break" out of the loop and stop reading.
			break
		}

		messages <- msg
	}
}

// processOpc receives the incoming OPC messages and dispatches them
func processOpc(messages chan *opc.Message, sinks []OpcSink) {
	defer func() {
		close(messages)
	}()

	for {
		msg := <-messages

		for _, s := range sinks {
			s.Write(msg)
		}
	}
}
