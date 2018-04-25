package main

//go:generate go run broadcaster/gen.go ControlPanel broadcaster/broadcast.tmpl

// ControlPanelProgram is a control panel program choice.
type ControlPanelProgram int

const (
	// CustomProgram means that the web clients can choose animations instead of just the control panel.
	CustomProgram ControlPanelProgram = 4
)

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

// NewControlPanelReceiver creates a new ControlPanelReceiver
func NewControlPanelReceiver() *ControlPanelReceiver {
	return &ControlPanelReceiver{controlPanel: make(chan ControlPanelMsg)}
}
