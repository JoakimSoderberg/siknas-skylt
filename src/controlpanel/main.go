package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jacobsa/go-serial/serial"
	"gopkg.in/alecthomas/kingpin.v2"
)

// ControlPanelMsg represents the state of the control panel hardware.
type ControlPanelMsg struct {
	Program    int    `json:"program,omitempty"`
	Color      [3]int `json:"color,omitempty"`
	Brightness int    `json:"brightness,omitempty"`
}

func (msg ControlPanelMsg) String() string {
	return fmt.Sprintf("Program: %v Color: (%v, %v, %v) Brightness: %v",
		msg.Program, msg.Color[0], msg.Color[1], msg.Color[2], msg.Brightness)
}

// NewControlPanelMsg Creates a new ControlPanelMsg struct based on an byte array
// containing a row of data read from the serial port.
// The expected data format is:
// 	"<Program> <Red> <Green> <Blue> <Brightness>\n"
// So for example, program 1 with full bright redness:
// 	1 255 0 0 255
func NewControlPanelMsg(msgBytes []byte) (*ControlPanelMsg, error) {
	// Get rid of any spaces or null characters first.
	msgStr := strings.TrimSpace(string(bytes.Trim(msgBytes[:], "\x00")))

	strs := strings.Split(msgStr, " ")
	if len(strs) < 5 {
		return nil, fmt.Errorf("Control message missing values. Expected %v but got %v", 5, len(strs))
	}

	for i := 0; i < 5; i++ {
		_, err := strconv.Atoi(strs[i])
		if err != nil {
			return nil, fmt.Errorf("Control message '%v' contains non-integer value at index %v: %v", strs, i, strs[i])
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

func openSerialPort(serialPort string) io.ReadWriteCloser {
	options := serial.OpenOptions{
		PortName:        serialPort,
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	log.Printf("Opening serial port: %v\n", options.PortName)

	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}

	return port
}

// serialPortListener listens for messages on the serial port.
func serialPortListener(messages chan ControlPanelMsg, port io.ReadWriteCloser) {
	defer port.Close()

	reader := bufio.NewReader(port)

	for {
		msgBytes, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatalf("Failed to read line from serial port: %v", err)
		}
		//log.Printf("Msgbytes: '%v'\n", string(msgBytes))
		msg, err := NewControlPanelMsg(msgBytes)

		if err != nil {
			log.Printf("Failed to create msg: %v", err)
			continue
		}

		messages <- *msg
	}
}

// registerSignalHandler handles interrupt signals.
func registerSignalHandler(c *websocket.Conn, port io.ReadWriteCloser) chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		for range sigChan {
			log.Println("\nReceived signal, exiting...")
			c.Close()
			os.Exit(0)
		}
	}()

	return sigChan
}

func connectWebsocket(addr string) *websocket.Conn {
	log.Printf("Connecting to %s...", addr)

	ws, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Fatal("Failed to connect to websocket server: ", err)
	}

	return ws
}

func websocketReader(ws *websocket.Conn, readDone chan struct{}) {
	defer ws.Close()
	defer close(readDone)

	defer func() {
		log.Println("Websocket reader ended")
	}()

	if *debug {
		log.Println("Websocket reader started")
	}

	// The reader will responed to things like PING even though
	// we don't care about any messages.
	for {
		if _, _, err := ws.NextReader(); err != nil {
			log.Println("Failed to read from websocket: ", err)
			break
		}
	}
}

// websocketWriter receives control panel messages and forwards them to the websocket server.
func websocketWriter(ws *websocket.Conn,
	messages chan ControlPanelMsg, interrupt chan os.Signal, readDone chan struct{}) {
	defer ws.Close()
	defer func() {
		log.Println("Websocket writer ended")
	}()

	if *debug {
		log.Println("Websocket writer started")
	}

	pingTicker := time.NewTicker(time.Second)
	defer pingTicker.Stop()

	for {
		select {
		// TODO: For some reason we get double messages here at times
		case msg := <-messages:
			// Receive a control panel message and forward it to the websocket.
			if *debug {
				log.Printf("Writing: '%v'\n", msg)
			}

			err := websocket.WriteJSON(ws, msg)
			if err != nil {
				log.Fatalf("Failed to write to websocket: %v\n", err)
			}
		case <-interrupt:
			// User wants to close.

			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := ws.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Failed to write:", err)
				return
			}

			// Wait for the reader to be done or a timeout.
			select {
			case <-readDone:
			case <-time.After(time.Second):
			}
			ws.Close()
			return
		}
	}
}

func Min(x, y int16) int16 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int16) int16 {
	if x > y {
		return x
	}
	return y
}

// TODO: Replace kingping with cobra!
var (
	server = kingpin.Flag("server_url", "Websocket server url to connect to.").
		Default("ws://localhost/ws/control_panel/").String()
	serialPort  = kingpin.Flag("serial_port", "The serial port to listen to.").String()
	debug       = kingpin.Flag("debug", "Enable debug output").Bool()
	dummy       = kingpin.Flag("dummy", "Use a dummy serial port").Bool()
	spaceToSend = kingpin.Flag("space_to_send", "Require pressing Space before sending a message in the dummy mode (--dummy --interactive must be used)").Bool()
)

func main() {
	kingpin.UsageTemplate(kingpin.DefaultUsageTemplate).Version("1.0").Author("Joakim Soderberg")
	kingpin.CommandLine.Help = "Siknas-skylt Control Panel Listener."
	kingpin.Parse()

	log.Println("Starting Siknas-skylt Control Panel listener...")

	// Channel receiving control panel messages via serial port.
	messages := make(chan ControlPanelMsg, 0)

	var port io.ReadWriteCloser
	if dummy != nil && *dummy {
		// For testing without the real hardware.
		dummyMsgs := make(chan string, 0)
		port = NewDummySerialPort(dummyMsgs)

		go DummyInteractive(dummyMsgs)
	} else {
		port = openSerialPort(*serialPort)
	}
	defer port.Close()

	// Listen for serial port messages.
	go serialPortListener(messages, port)

	ws := connectWebsocket(*server)
	defer ws.Close()

	interrupt := registerSignalHandler(ws, port)

	// When this is closed we are done reading from the websocket.
	readDone := make(chan struct{})

	// Websocket reader.
	go websocketReader(ws, readDone)

	// Websocket writer.
	go websocketWriter(ws, messages, interrupt, readDone)

	<-readDone
	log.Println("Exiting")
}
