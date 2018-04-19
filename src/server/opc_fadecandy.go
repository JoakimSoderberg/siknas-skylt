package main

import (
	"encoding/json"
	"fmt"

	opc "github.com/kellydunn/go-opc"
)

// FadecandyColorCorrectionMsg represents a color correction message for the Fadecandy LED controller board.
type FadecandyColorCorrectionMsg struct {
	Gamma      float32   `json:"gamma"`
	Whitepoint []float32 `json:"whitepoint"`
}

// CreateFadecandyColorCorrectionPacket creates a color correction packet for the Fadecandy LED controller.
func CreateFadecandyColorCorrectionPacket(gamma, red, green, blue float32) (*opc.Message, error) {
	msg := opc.NewMessage(0)

	contentMsg := FadecandyColorCorrectionMsg{Gamma: gamma, Whitepoint: []float32{red, green, blue}}
	contentBytes, err := json.Marshal(contentMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal color correction message: ", err)
	}

	data := []byte{0x00, 0x01} // Command ID for color correction.
	data = append(data, contentBytes...)

	msg.SystemExclusive(
		[]byte{0x00, 0x01}, // System ID for Fadecandy board.
		data)
	msg.SetLength(uint16(len(data) + 2)) // Include System ID 2 bytes

	return msg, nil
}
