package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	var rootCmd = &cobra.Command{Use: "siknas-skylt thumbnail generator", Run: func(c *cobra.Command, args []string) {}}
	rootCmd.Flags().String("logo-svg", "siknas-skylt.svg", "Path to Siknäs logo")
	rootCmd.Flags().String("led-layout", "layout.json", "Path to the LED layout.json")
	rootCmd.Flags().String("host", "ws://localhost:8080/ws/opc", "OPC websocket server host including port")
	rootCmd.Flags().Duration("capture-duration", 10, "Duration of data we should capture (in seconds)")
	rootCmd.Flags().String("output", "output.svg", "Output filename")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}

	viper.BindPFlags(rootCmd.Flags())

	loadSvg()

	host := viper.GetString("host")
	// TODO: Build URL
	captureDuration := viper.GetDuration("capture-duration")

	ws := connectWebsocket(host)
	defer ws.Close()

	interrupt := registerSignalHandler(ws)

	done := make(chan struct{})

	// Websocket reader.
	go websocketReader(ws, interrupt, done)

	// Websocket writer.
	go websocketWriter(ws, interrupt, done)

	time.Sleep(time.Second * captureDuration)
	log.Printf("Finished capturing after %s\n", time.Second*captureDuration)
	close(done)
	time.Sleep(1 * time.Second)
}

var signalCount int

// registerSignalHandler handles interrupt signals.
func registerSignalHandler(c *websocket.Conn) chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signalCount++

	go func() {
		for range sigChan {
			log.Println("\nReceived signal, exiting...")

			if signalCount > 1 {
				c.Close()
				os.Exit(0)
			}
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

// OpcMessageHeader defines the header of the Open Pixel Control (OPC) Protocol
type OpcMessageHeader struct {
	Channel byte
	Command byte
	Length  uint16
}

func websocketReader(ws *websocket.Conn, interrupt chan os.Signal, done chan struct{}) {
	var opcMsgHdr OpcMessageHeader

	log.Println("Starting websocket reader")

	for {
		select {
		default:
			messageType, messageData, err := ws.ReadMessage()
			if err != nil {
				log.Println("Failed to read: ", err)
				return
			}

			if messageType != websocket.BinaryMessage {
				log.Println("ERROR: Got a Text message on the OPC Websocket, expected Binary")
				break
			}

			buf := bytes.NewBuffer(messageData[0:binary.Size(opcMsgHdr)])
			err = binary.Read(buf, binary.BigEndian, &opcMsgHdr)
			if err != nil {
				log.Println("ERROR: Failed to read OPC message: ", err)
				break
			}

			realMsgLength := uint16(len(messageData) - binary.Size(opcMsgHdr))

			if opcMsgHdr.Length != realMsgLength {
				log.Printf("ERROR: Got a %d byte invalid OPC message. Header says %d, got %d bytes\n", opcMsgHdr.Length, opcMsgHdr.Length, realMsgLength)
				break
			}

			log.Printf("Msg %d bytes\n", opcMsgHdr.Length)
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
	//defer close(done)

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

func loadSvg() {
	svgLogoPath := viper.GetString("logo-svg")
	ledLayoutPath := viper.GetString("led-layout")
	outputPath := viper.GetString("output")

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

	for i, pos := range ledPositions {
		circleNode := ledGroupNode.CreateNode("circle")
		circleNode.SetAttributeValue("id", fmt.Sprintf("led%d", i))
		circleNode.SetAttributeValue("cx", fmt.Sprintf("%f", (pos.Point[0]*width*0.81)+(width*0.05)))
		circleNode.SetAttributeValue("cy", fmt.Sprintf("%f", (pos.Point[1]*height*0.55)+(height*0.20)))
		circleNode.SetAttributeValue("r", "10")

		//addLedAnimation(i, circleNode)

		//circleNode.SetAttributeValue("style", "fill: rgb(255,255,255)")
	}

	ioutil.WriteFile(outputPath, []byte(doc.XMLPretty()), 0644)
}

func addLedAnimation(i int, circleNode *xmldom.Node) {
	animNode := circleNode.CreateNode("animate")
	animNode.SetAttributeValue("attributeName", "fill")
	animNode.SetAttributeValue("dur", "0.1s") // TODO: Get avreage time between frames
	animNode.SetAttributeValue("repeatCount", "indefinite")

	colors := make([]string, 2)
	for i := 0; i < 2; i++ {
		colors[i] = fmt.Sprintf("rgb(%d,%d,%d)", r.Int()%255, r.Int()%255, r.Int()%255)
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
