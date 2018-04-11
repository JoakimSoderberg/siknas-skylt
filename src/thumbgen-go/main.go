package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/gorilla/websocket"
	xmldom "github.com/subchen/go-xmldom"
)

// OpcMessageHeader defines the header of the Open Pixel Control (OPC) Protocol.
type OpcMessageHeader struct {
	Channel byte
	Command byte
	Length  uint16
}

// OpcMessage defines a OPC message including header.
type OpcMessage struct {
	Header OpcMessageHeader
	Data   []byte
}

// RGB returns the Red, Green and Blue values between 0-255 for a given pixel.
func (m *OpcMessage) RGB(ledIndex int) (uint8, uint8, uint8) {
	i := 3 * ledIndex
	return m.Data[i], m.Data[i+1], m.Data[i+2]
}

var opcMessages []OpcMessage

func main() {
	var rootCmd = &cobra.Command{Use: "thumbgen", Run: func(c *cobra.Command, args []string) {}}
	rootCmd.Flags().String("logo-svg", "siknas-skylt.svg", "Path to Siknäs logo")
	rootCmd.Flags().String("led-layout", "layout.json", "Path to the LED layout.json")
	rootCmd.Flags().String("host", "localhost:8080", "OPC websocket server host including port")
	rootCmd.Flags().String("ws-opc-path", "/ws/opc", "OPC websocket path to connect to")
	rootCmd.Flags().String("ws-path", "/ws", "Websocket control path to connect to")
	rootCmd.Flags().Duration("capture-duration", 10*time.Second, "Duration of data we should capture (in seconds)")
	rootCmd.Flags().String("output", "output.svg", "Output filename") // TODO: change to directory
	// TODO: Add option to fetch SVG from server.
	// TODO: Add option to force overwrite existing thumbnails (otherwise skip them).
	// TODO: Add option to only list the sketches from the server.

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}

	viper.BindPFlags(rootCmd.Flags())

	if viper.IsSet("help") {
		os.Exit(0)
	}

	// TODO: Connect to the "control websocket"
	// TODO: Get list of animations
	// TODO: Iterate over animations. Check if exists in output directory
	// TODO: Send a message to the server to switch to the current animation. Sleep for a while
	// TODO: Record animation for a given amount and save to disk.

	// TODO: Move this to a separate file.
	url := url.URL{Scheme: "ws", Host: viper.GetString("host"), Path: viper.GetString("ws-opc-path")}
	captureDuration := viper.GetDuration("capture-duration")

	ws := connectWebsocket(url.String())
	defer ws.Close()

	interrupt := registerSignalHandler(ws)

	done := make(chan struct{})

	// Websocket reader.
	go websocketReader(ws, interrupt, done)

	// Websocket writer.
	go websocketWriter(ws, interrupt, done)

	// TODO: Fix closing down nicely.
	time.Sleep(captureDuration)
	log.Printf("Finished capturing after %s\n", captureDuration)
	close(done)

	createOutputSVG()
}

// registerSignalHandler handles interrupt signals.
func registerSignalHandler(c *websocket.Conn) chan os.Signal {
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

func websocketReader(ws *websocket.Conn, interrupt chan os.Signal, done chan struct{}) {
	var opcMsg OpcMessage

	log.Println("Starting websocket reader")

	started := false

	for {
		select {
		default:
			messageType, messageData, err := ws.ReadMessage()
			if err != nil {
				log.Println("Failed to read: ", err)
				return
			}

			if !started {
				started = true
				log.Printf("Started capturing %s of animation\n", viper.GetDuration("capture-duration"))
			}

			if messageType != websocket.BinaryMessage {
				log.Println("ERROR: Got a Text message on the OPC Websocket, expected Binary")
				break
			}

			buf := bytes.NewBuffer(messageData[0:binary.Size(opcMsg.Header)])
			err = binary.Read(buf, binary.BigEndian, &opcMsg.Header)
			if err != nil {
				log.Println("ERROR: Failed to read OPC message: ", err)
				break
			}

			realMsgLength := uint16(len(messageData) - binary.Size(opcMsg.Header))

			if opcMsg.Header.Length != realMsgLength {
				log.Printf("ERROR: Got a %d byte invalid OPC message. Header says %d, got %d bytes\n", opcMsg.Header.Length, opcMsg.Header.Length, realMsgLength)
				break
			}

			// Note we don't really need the OPC Length here, since this is Websockets
			// and we already have a known message length.
			opcMsg.Data = messageData[binary.Size(opcMsg.Header):]

			opcMessages = append(opcMessages, opcMsg)

		case <-interrupt:
			// TODO: Fix this
			log.Println("Reader got interrupted...")
			return
		case <-done:
			log.Println("Reader done...")
			return
		}
	}
}

// websocketWriter receives control panel messages and forwards them to the websocket server.
func websocketWriter(ws *websocket.Conn, interrupt chan os.Signal, done chan struct{}) {
	defer ws.Close()

	pingTicker := time.NewTicker(time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-done:
		case <-interrupt:
			// User wants to close.
			log.Println("Writer got interrupted, attempting clean Websocket close...")

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
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

// Creates the output SVG including the animation recorded.
func createOutputSVG() {
	svgLogoPath := viper.GetString("logo-svg")
	ledLayoutPath := viper.GetString("led-layout")
	outputPath := viper.GetString("output") // TODO: REmove and pass as argument instead.

	doc := xmldom.Must(xmldom.ParseFile(svgLogoPath))
	svg := doc.Root

	width, err := strconv.ParseFloat(svg.GetAttributeValue("width"), 32)
	if err != nil {
		log.Fatalln("Failed to parse SVG width: ", width)
	}

	height, err := strconv.ParseFloat(svg.GetAttributeValue("height"), 32)
	if err != nil {
		log.Fatalln("Failed to parse SVG height: ", height)
	}

	// Make the text black as a background to the LEDs.
	svg.QueryEach("//g[@id = 'Siknas']//path",
		func(i int, node *xmldom.Node) {
			node.SetAttributeValue("style", "fill: black")
		})

	// Create a group for the LEDs
	ledGroupNode := svg.CreateNode("g")

	// A copy of the "Siknäs" text paths exists in the SVG defined as a clip-path
	// we want the LEDs to be stay inside of this. So we clip the group using it.
	ledGroupNode.SetAttributeValue("clip-path", "url(#SiknasClipPath)")

	ledPositions, err := readLEDLayout(ledLayoutPath)
	if err != nil {
		log.Fatalln(err)
	}

	// Create the circles that represents the LED:s
	for i, pos := range ledPositions {
		circleNode := ledGroupNode.CreateNode("circle")
		circleNode.SetAttributeValue("id", fmt.Sprintf("led%d", i))
		circleNode.SetAttributeValue("cx", fmt.Sprintf("%f", (pos.Point[0]*width*0.81)+(width*0.05)))
		circleNode.SetAttributeValue("cy", fmt.Sprintf("%f", (pos.Point[1]*height*0.55)+(height*0.20)))
		circleNode.SetAttributeValue("r", "10")

		addLedAnimation(i, circleNode)
	}

	ioutil.WriteFile(outputPath, []byte(doc.XMLPretty()), 0644)
}

func addLedAnimation(ledIndex int, circleNode *xmldom.Node) {
	animNode := circleNode.CreateNode("animate")
	animNode.SetAttributeValue("attributeName", "fill")
	animNode.SetAttributeValue("dur", fmt.Sprintf("%.2fs", viper.GetDuration("capture-duration").Seconds()))
	animNode.SetAttributeValue("repeatCount", "indefinite")

	colors := make([]string, len(opcMessages))

	for i := 0; i < len(opcMessages); i++ {
		r, g, b := opcMessages[i].RGB(ledIndex)
		colors[i] = fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
	}
	animNode.SetAttributeValue("values", strings.Join(colors, ";"))
}

// LedPosition is a point for a LED.
type LedPosition struct {
	Point [3]float64 `json:"point"`
}

func readLEDLayout(path string) ([]LedPosition, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %s", path, err)
	}

	ledPositions := make([]LedPosition, 0)
	err = json.Unmarshal(bytes, &ledPositions)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON from %s: %s", path, err)
	}

	return ledPositions, nil
}
