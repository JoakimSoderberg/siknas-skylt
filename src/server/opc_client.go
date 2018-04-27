package main

import (
	"log"
	"net"
	"time"

	"github.com/cenkalti/backoff"
)

// RunOpcClient will start an OPC client that continues to reconnect with an exponential
// backoff to the given server.
func RunOpcClient(protocol string, host string, maxBackoff time.Duration, bcast *OpcBroadcaster) {
	// TODO: accept a OPC server struct instead of just host.
	opcConnect := func() error {
		conn, err := net.Dial("tcp", host)
		if err != nil {
			return err
		}

		log.Println("[OPC outgoing client] Connected:", host)

		// We listen to incoming OPC messages being broadcasted and then
		// forward them to the server we are connected to.
		opcReceiver := NewOpcReceiver()
		bcast.Push(opcReceiver)
		log.Println("[OPC outgoing client] Added as OPC broadcast receiver:", host)
		defer func() {
			bcast.Pop(opcReceiver)
			log.Println("[OPC outgoing client] Removed as OPC broadcast receiver:", host)
			conn.Close()
			log.Println("[OPC outgoing client] Disconnected:", host)
		}()

		for {
			select {
			case msg := <-opcReceiver.opcMessages:
				conn.SetWriteDeadline(time.Now().Add(time.Second))
				_, err := conn.Write(msg.ByteArray())
				if err != nil {
					//log.Println("[OPC outgoing client] Failed to write:", err)
					return err
				}
				//log.Printf("[OPC outgoing client] sent %v bytes\n", n)
			}
		}
	}

	// Keep reconnecting forever.
	for {
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = maxBackoff
		b.MaxInterval = 30 * time.Second
		err := backoff.RetryNotify(opcConnect, b, func(err error, duration time.Duration) {
			log.Printf("[OPC outgoing client] Reconnect in %.2fs: %v\n", duration.Seconds(), err) // TODO: Round this duration
		})
		if err != nil {
			log.Printf("[OPC outgoing client] Failed retry while reconnecting: %v\n", err)
		}
	}
}
