package main

//go:generate go run broadcaster/gen.go Opc broadcaster/broadcast.tmpl

import (
	"log"
	"net"
	//"time"

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

// handleOpcCon handles connections from OPC clients.
func handleOpcCon(messages chan *opc.Message, conn net.Conn) {
	defer func() {
		log.Println("[OPC incoming client] Disconnected:", conn.RemoteAddr())
		conn.Close()
	}()

	/*
			// TODO: Take into account what the current max brightness is set to.
			// We will ramp it up by interjecting color correction packets that raises the brightness gradually.
			brightness := float32(0.0)
			fadeTimeout := 3000
			fadeDoneTimer := time.NewTimer(time.Duration(fadeTimeout) * time.Millisecond)
			defer fadeDoneTimer.Stop()
			fadeTicker := time.NewTicker(time.Duration(fadeTimeout) * time.Millisecond / 100)
			defer fadeTicker.Stop()


			log.Println("Start fading")

			// TODO: Break out into separate function.
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
					colorCorrMsg, err := CreateFadecandyColorCorrectionPacket(float32(2.5), brightness, brightness, brightness)
					if err != nil {
						log.Println("[OPC incoming client]: ", err)
						return
					}
					messages <- colorCorrMsg
					brightness += float32(1.0 / 10.0) // TODO: Fix this to ramp correctly.

				case <-fadeDoneTimer.C:
					// Fade completed (make sure we are att full brightness).
					colorCorrMsg, err := CreateFadecandyColorCorrectionPacket(float32(2.5), 1.0, 1.0, 1.0)
					if err != nil {
						log.Println("[OPC incoming client]: ", err)
						return
					}
					messages <- colorCorrMsg
					log.Println("Done fading")
					break fadeLoop

				default:
					// Send the regular OPC message.
					messages <- msg
				}
			}
	*/

	// Normal parsing of packets.
	for {
		msg, err := opc.ReadOpc(conn)
		if err != nil {
			log.Println("[OPC incoming client] Failed to read OPC:", conn.RemoteAddr())
			return
		}

		messages <- msg

		// TODO: Make it possible to ignore incoming messages (don't forward them)
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
