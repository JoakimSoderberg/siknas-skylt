package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jacobsa/go-serial/serial"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	server     = kingpin.Arg("server", "Websocket server to connect to. (defaults to localhost)").Default("localhost").String()
	port       = kingpin.Flag("port", "The port to use.").Default("80").Int()
	serialPort = kingpin.Flag("serial_port", "The serial port to listen to.").String()
)

// ControlPanelMsg represents the state of the control panel hardware.
type ControlPanelMsg struct {
	Program    int    `json:"program,omitempty"`
	Color      [3]int `json:"color,omitempty"`
	Brightness int    `json:"brightness,omitempty"`
}

func (msg *ControlPanelMsg) String() string {
	return fmt.Sprintf("Program: %v Color: (%v, %v, %v) Brightness: %v", msg.Program, msg.Color[0], msg.Color[0], msg.Color[0], msg.Brightness)
}

// NewControlPanelMsg Creates a new ControlPanelMsg struct based on an byte array
// containing a row of data read from the serial port.
func NewControlPanelMsg(msgBytes []byte) (*ControlPanelMsg, error) {
	strs := strings.Split(string(msgBytes[:]), " ")
	if len(strs) < 5 {
		return nil, fmt.Errorf("Control message missing values. Expected %v but got %v", 5, len(strs))
	}

	for i := 0; i < 5; i++ {
		_, err := strconv.Atoi(strs[i])
		if err != nil {
			return nil, fmt.Errorf("Control message contains non-integer value: %v", strs[i])
		}
	}

	msg := new(ControlPanelMsg)
	msg.Program, _ = strconv.Atoi(strs[0])
	msg.Color[0], _ = strconv.Atoi(strs[1])
	msg.Color[1], _ = strconv.Atoi(strs[2])
	msg.Color[2], _ = strconv.Atoi(strs[3])
	msg.Brightness, _ = strconv.Atoi(strs[4])

	return msg, nil
}

// SerialPortListener listens for messages on the serial port.
func SerialPortListener(messages chan ControlPanelMsg) {
	options := serial.OpenOptions{
		PortName:        *serialPort,
		BaudRate:        19200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	// Make sure to close it later.
	defer port.Close()

	reader := bufio.NewReader(port)

	for {
		msgBytes, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatalf("Failed to read line from serial port: %v", err)
		}

		msg, err := NewControlPanelMsg(msgBytes)

		if err == nil {
			messages <- *msg
		}
	}
}

func main() {
	kingpin.UsageTemplate(kingpin.DefaultUsageTemplate).Version("1.0").Author("Joakim Soderberg")
	kingpin.CommandLine.Help = "Siknas-skylt Control Panel Websocket Client."
	kingpin.Parse()

	messages := make(chan ControlPanelMsg)

	go SerialPortListener(messages)

	// TODO: Connect to websocket

	for {
		msg := <-messages
		log.Println(msg)
		// TODO: Send to websocket
	}
}
