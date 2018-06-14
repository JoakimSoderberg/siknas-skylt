package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	xmldom "github.com/subchen/go-xmldom"
)

var logoSvgData []byte

func getLogoSvgReader() *bytes.Reader {
	if len(logoSvgData) > 0 {
		return bytes.NewReader(logoSvgData)
	}

	if viper.GetString("logo-svg") != "" {
		// Read logo from local file.
		var err error
		svgLogoPath := viper.GetString("logo-svg")
		logoSvgData, err = ioutil.ReadFile(svgLogoPath)
		if err != nil {
			log.Fatalf("Failed to read logo SVG %s: %s\n", svgLogoPath, err)
		}
	} else {
		// Download logo from server is none was specified.
		svgLogoURL := fmt.Sprintf("http://%s/images/siknas-skylt.svg", viper.GetString("host"))
		response, err := http.Get(svgLogoURL)
		if err != nil {
			log.Fatalf("No logo was specified using --logo-svg and failed to get '%s': %s\n", svgLogoURL, err)
		}
		defer response.Body.Close()

		logoSvgData, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("Failed to read logo from server '%s': %s\n", svgLogoURL, err)
		}
	}

	return bytes.NewReader(logoSvgData)
}

func createBaseSVG() (*xmldom.Document, *xmldom.Node, float64, float64) {
	logoReader := getLogoSvgReader()
	doc := xmldom.Must(xmldom.Parse(logoReader))
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

	return doc, ledGroupNode, width, height
}

func createCircleNode(ledGroupNode *xmldom.Node, pos *LedPosition, width, height float64, ledIndex int) *xmldom.Node {
	circleNode := ledGroupNode.CreateNode("circle")
	circleNode.SetAttributeValue("id", fmt.Sprintf("led%d", ledIndex))
	circleNode.SetAttributeValue("cx", fmt.Sprintf("%f", (pos.Point[0]*width*0.81)+(width*0.05)))
	circleNode.SetAttributeValue("cy", fmt.Sprintf("%f", (pos.Point[1]*height*0.55)+(height*0.20)))
	circleNode.SetAttributeValue("r", "10")

	return circleNode
}

func createSVGFrames(outputPath string, name string, opcMessages []OpcMessage, ledPositions []LedPosition) {
	// We re-create the directory so we don't have old frames laying about.
	os.RemoveAll(outputPath)
	err := os.MkdirAll(outputPath, 0644)
	if err != nil {
		log.Fatalf("Failed to create directory '%v': %v", outputPath, err)
	}

	for i := 0; i < len(opcMessages); i++ {
		outputFilePath := path.Join(outputPath, fmt.Sprintf("%s%0*d.svg", name, 6, i))

		doc, ledGroupNode, width, height := createBaseSVG()

		for ledIndex, pos := range ledPositions {
			circleNode := createCircleNode(ledGroupNode, &pos, width, height, ledIndex)
			r, g, b := opcMessages[i].RGB(ledIndex)
			circleNode.SetAttributeValue("style", fmt.Sprintf("fill: rgb(%d,%d,%d)", r, g, b))
		}

		ioutil.WriteFile(outputFilePath, []byte(doc.XMLPretty()), 0644)
	}
}

// Creates the output SVG including the animation recorded.
func createAnimatedSVG(outputPath string, opcMessages []OpcMessage, ledPositions []LedPosition) {
	doc, ledGroupNode, width, height := createBaseSVG()

	// Create the circles that represents the LED:s
	for i, pos := range ledPositions {
		circleNode := createCircleNode(ledGroupNode, &pos, width, height, i)
		addLedAnimation(opcMessages, i, circleNode)
	}

	ioutil.WriteFile(outputPath, []byte(doc.XMLPretty()), 0644)
}

func addLedAnimation(opcMessages []OpcMessage, ledIndex int, circleNode *xmldom.Node) {

	colors := make([]string, len(opcMessages))

	for i := 0; i < len(opcMessages); i++ {
		r, g, b := opcMessages[i].RGB(ledIndex)
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
