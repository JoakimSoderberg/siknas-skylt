package main

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/kellydunn/go-opc"
)

// RunOpcClient will start an OPC client that continues to reconnect with an exponential
// backoff to the given server.
func RunOpcClient(protocol string, host string, maxBackoff time.Duration, bcast *OpcBroadcaster) error {
	c := opc.NewClient()

	opcConnect := func() error {
		err := c.Connect(protocol, host)
		if err != nil {
			return fmt.Errorf("failed to connect to OPC server %v: %v", host, err)
		}

		// We listen to incoming OPC messages being broadcasted and then
		// forward them to the server we are connected to.
		opcReceiver := NewOpcReceiver()
		bcast.Push(opcReceiver)
		defer bcast.Pop(opcReceiver)

		for {
			select {
			case msg := <-opcReceiver.opcMessages:
				err := c.Send(msg)
				if err != nil {
					return fmt.Errorf("failed to send to OPC server %v: %v", host, err)
				}
			}
		}
	}

	// Keep reconnecting forever.
	for {
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = 30 * time.Second
		b.MaxInterval = 15 * time.Second
		backoff.Retry(opcConnect, b)
	}
}
