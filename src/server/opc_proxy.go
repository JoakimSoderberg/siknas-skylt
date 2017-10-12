package main

import (
	"log"
	"net"

	"github.com/kellydunn/go-opc"
)

// TODO: Make the broadcaster in control_panel.go generic so we can use it here also
// for both websocket clients and other OPC receivers.

type OpcReceiver interface {
	Write(msg *opc.Message, channel int)
}

// RunOPCServer runs an Open Pixel Protocol server.
func RunOPCServer(protocol string, port string) error {
	// TODO: Once

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
	//go processOpc(messages) // TODO: Fix

	return nil
}

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
func processOpc(messages chan *opc.Message, receivers OpcReceiver) {
	defer func() {
		close(messages)
	}()

	for {
		//msg := <-messages

		// TODO: Send to websocket clients
		// TODO: Send to other OPC servers
	}
}
