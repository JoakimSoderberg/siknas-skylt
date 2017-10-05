package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/jacobsa/go-serial/serial"
	"golang.org/x/net/websocket"
	"gopkg.in/alecthomas/kingpin.v2"
)

// ControlPanelMsg represents the state of the control panel hardware.
type ControlPanelMsg struct {
	Program    int    `json:"program,omitempty"`
	Color      [3]int `json:"color,omitempty"`
	Brightness int    `json:"brightness,omitempty"`
}

func (msg ControlPanelMsg) String() string {
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
	// Make sure to close it later.
	defer port.Close()

	reader := bufio.NewReader(port)

	for {
		msgBytes, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatalf("Failed to read line from serial port: %v", err)
		}

		msg, err := NewControlPanelMsg(msgBytes)

		if err != nil {
			log.Printf("Failed to create msg: %v", err)
			continue
		}

		messages <- *msg
	}
}

// registerSignalHandler handles interrupt signals.
func registerSignalHandler(c *websocket.Conn, port io.ReadWriteCloser) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		for range ch {
			log.Println("\nReceived signal, exiting...")
			c.Close()
			os.Exit(0)
		}
	}()
}

func connectWebsocket(addr string) *websocket.Conn {
	log.Printf("Connecting to %s...", addr)

	conf, err := websocket.NewConfig(addr, addr)
	if err != nil {
		log.Fatal(err)
	}

	ws, err := websocket.DialConfig(conf)
	if err != nil {
		log.Fatal(err)
	}
	return ws
}

var (
	server = kingpin.Flag("server_url", "Websocket server url to connect to.").
		Default("ws://localhost/ws/control_panel/").String()
	serialPort = kingpin.Flag("serial_port", "The serial port to listen to.").String()
	debug      = kingpin.Flag("debug", "Enable debug output").Bool()
)

func main() {
	kingpin.UsageTemplate(kingpin.DefaultUsageTemplate).Version("1.0").Author("Joakim Soderberg")
	kingpin.CommandLine.Help = "Siknas-skylt Control Panel Listener."
	kingpin.Parse()

	log.Println("Starting Siknas-skylt Control Panel listener...")

	messages := make(chan ControlPanelMsg)

	ws := connectWebsocket(*server)
	port := openSerialPort(*serialPort)

	registerSignalHandler(ws, port)

	// Listen for serial port messages non-blocking.
	go serialPortListener(messages, port)

	for {
		select {
		case msg := <-messages:
			if *debug {
				log.Println(msg)
			}

			err := websocket.JSON.Send(ws, msg)
			if err != nil {
				log.Fatalf("Failed to write to websocket: %v\n", err)
			}
		default:

		}
	}
}
