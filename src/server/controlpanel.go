package main

//go:generate go run broadcaster/gen.go ControlPanel broadcaster/broadcast.tmpl

// ControlPanelMsg represents the state of the control panel hardware.
type ControlPanelMsg struct {
	Program    int    `json:"program,omitempty"`
	Color      [3]int `json:"color,omitempty"`
	Brightness int    `json:"brightness,omitempty"`
}

// ControlPanelReceiver is a client that wants to listen to control panel messages.
type ControlPanelReceiver struct {
	controlPanel chan ControlPanelMsg
}
