package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	xmldom "github.com/subchen/go-xmldom"
)

// Creates the output SVG including the animation recorded.
func createOutputSVG(outputPath string, opcMessages []OpcMessage) {
	svgLogoPath := viper.GetString("logo-svg")
	ledLayoutPath := viper.GetString("led-layout")

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

		addLedAnimation(opcMessages, i, circleNode)
	}

	ioutil.WriteFile(outputPath, []byte(doc.XMLPretty()), 0644)
}

func addLedAnimation(opcMessages []OpcMessage, ledIndex int, circleNode *xmldom.Node) {

	colors := make([]string, len(opcMessages))

	for i := 0; i < len(opcMessages); i++ {
		msg := opcMessages[i]

		// TODO: Figure out what's wrong here :(
		if uint16((ledIndex*3)+2) >= msg.Header.Length {
			return
		}
		r, g, b := msg.RGB(ledIndex)
		colors[i] = fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
	}

	animNode := circleNode.CreateNode("animate")
	animNode.SetAttributeValue("attributeName", "fill")
	animNode.SetAttributeValue("dur", fmt.Sprintf("%.2fs", viper.GetDuration("capture-duration").Seconds()))
	animNode.SetAttributeValue("repeatCount", "indefinite")
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
