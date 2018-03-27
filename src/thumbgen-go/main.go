package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	xmldom "github.com/subchen/go-xmldom"
)

func main() {
	var rootCmd = &cobra.Command{Use: "siknas-skylt thumbnail generator", Run: func(c *cobra.Command, args []string) {}}
	rootCmd.Flags().String("logo-svg", "", "Path to Siknäs logo")
	rootCmd.MarkFlagRequired("logo-svg")
	rootCmd.Flags().String("led-layout", "", "Path to the LED layout.json")
	rootCmd.MarkFlagRequired("led-layout")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}

	viper.BindPFlags(rootCmd.Flags())

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

	cssNode := svg.QueryOne("//defs//style")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var buffer bytes.Buffer
	buffer.WriteString(cssNode.Text)

	for i, _ := range ledPositions {
		buffer.WriteString(fmt.Sprintf("\n#led%d { fill: rgb(%d,%d,%d) }", i, r.Int()%255, r.Int()%255, r.Int()%255))
	}

	cssNode.Text = buffer.String()

	for i, pos := range ledPositions {
		circleNode := ledGroupNode.CreateNode("circle")
		circleNode.SetAttributeValue("id", fmt.Sprintf("led%d", i))
		circleNode.SetAttributeValue("cx", fmt.Sprintf("%f", (pos.Point[0]*width*0.81)+(width*0.05)))
		circleNode.SetAttributeValue("cy", fmt.Sprintf("%f", (pos.Point[1]*height*0.55)+(height*0.20)))
		circleNode.SetAttributeValue("r", "10")
		//circleNode.SetAttributeValue("style", "fill: rgb(255,255,255)")
	}

	ioutil.WriteFile("changed.svg", []byte(doc.XML()), 0644)
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

	//var ledPositions LedPositions
	ledPositions := make([]LedPosition, 0)
	err = json.Unmarshal(bytes, &ledPositions)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON from %s: %s", path, err)
	}

	return ledPositions, nil
}
