package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/JoakimSoderberg/go-opc"
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
	rootCmd.Flags().String("host", "localhost:7890", "OPC server host including port")
	rootCmd.Flags().Duration("capture-duration", 10, "Duration of data we should capture (in seconds)")
	rootCmd.Flags().String("output", "output.svg", "Output filename")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}

	viper.BindPFlags(rootCmd.Flags())

	loadSvg()

	host := viper.GetString("host")
	captureDuration := viper.GetDuration("capture-duration")

	log.Println("Connecting to ", host)

	c := opc.NewClient()

	conn, err := net.DialTimeout("tcp", host, time.Second*10)
	if err != nil {
		log.Fatalln("Failed to connect: ", err)
	}
	defer conn.Close()

	bufReader := bufio.NewReader(conn)

	timeoutTicker := time.NewTicker(time.Second * 10)
	defer timeoutTicker.Stop()

	var msg opc.Message

	for {
		select {
		case <-time.After(time.Second * captureDuration):
			break
		default:
			headerBytes := make([]byte, 4)
			n, err := bufReader.Read(headerBytes)
			if err != nil {
				break
			}

			buffer := bytes.NewBuffer(headerBytes)
			err = binary.Read(buffer, binary.BigEndian, &msg)
			if err != nil {
				log.Fatal("binary.Read failed", err)
			}

			dataBytes := make([]byte, msg.Length())
			n, err = bufReader.Read(dataBytes)
			log.Println("Read ", n)
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
