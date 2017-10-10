package main

import (
	"log"
	"sync"
)

// ControlPanelMsg represents the state of the control panel hardware.
type ControlPanelMsg struct {
	Program    int    `json:"program,omitempty"`
	Color      [3]int `json:"color,omitempty"`
	Brightness int    `json:"brightness,omitempty"`
}

type ControlPanelClient struct {
	controlPanel chan ControlPanelMsg
}

type ControlPanelBroadcaster struct {
	sync.Mutex
	clients []*ControlPanelClient
}

// Push adds a new client as a broadcast listener to the control panel.
func (bcast *ControlPanelBroadcaster) Push(c *ControlPanelClient) {
	bcast.Lock()
	defer bcast.Unlock()

	bcast.clients = append(bcast.clients, c)

	log.Println("Added control panel broadcast listening client")
}

// Pop removes a client from the control panel broadcast.
func (bcast *ControlPanelBroadcaster) Pop(c *ControlPanelClient) {
	bcast.Lock()
	defer bcast.Unlock()

	i := -1
	for j, cur := range bcast.clients {
		if cur == c {
			i = j
			break
		}
	}

	if i < 0 {
		return
	}

	// TODO: Keeping clients in a slice might not be the best solution?
	copy(bcast.clients[i:], bcast.clients[i+1:])
	bcast.clients[len(bcast.clients)-1] = nil // or the zero value of T
	bcast.clients = bcast.clients[:len(bcast.clients)-1]

	log.Println("Removed control panel broadcast listening client")
}

// Broadcast will send an incoming message from the control panel to all listening channels.
func (bcast *ControlPanelBroadcaster) Broadcast(routine func(*ControlPanelClient)) {
	bcast.Lock()
	defer bcast.Unlock()

	// Broadcasts to all clients.
	for _, c := range bcast.clients {
		log.Println("Broadcasting to ", c)
		routine(c)
	}
}
