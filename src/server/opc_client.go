package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/cenkalti/backoff"
)

// RunOpcClient will start an OPC client that continues to reconnect with an exponential
// backoff to the given server.
func RunOpcClient(protocol string, host string, maxBackoff time.Duration, bcast *OpcBroadcaster) error {

	opcConnect := func() error {
		log.Printf("OPC proxy client connecting to %v...", host)

		conn, err := net.DialTimeout(protocol, host, 3*time.Second)
		if err != nil {
			return fmt.Errorf("failed proxy connect to OPC server %v: %v", host, err)
		}

		// We listen to incoming OPC messages being broadcasted and then
		// forward them to the server we are connected to.
		opcReceiver := NewOpcReceiver()
		bcast.Push(opcReceiver)
		defer bcast.Pop(opcReceiver)
		defer conn.Close()

		for {
			select {
			case msg := <-opcReceiver.opcMessages:
				_, err := conn.Write(msg.ByteArray())
				if err != nil {
					opcErr := fmt.Errorf("failed to send to OPC server %v: %v", host, err)
					log.Println(opcErr)
					return opcErr
				}

				conn.SetReadDeadline(time.Now())
				if _, err := conn.Read([]byte{}); err == io.EOF {
					return err
				}
			}
		}
	}

	// Keep reconnecting forever.
	for {
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = maxBackoff
		b.MaxInterval = 30 * time.Second
		err := backoff.Retry(opcConnect, b)
		if err != nil {
			log.Printf("Reconnecting: %v\n", err)
		}
	}
}
