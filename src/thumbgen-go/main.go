package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

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

	// Make the text black by default.
	svg.QueryEach("//g[@id = 'Siknas']//path",
		func(i int, node *xmldom.Node) {
			node.SetAttributeValue("style", "fill: black")
		})

	// TODO: Check error
	width, err := strconv.ParseFloat(svg.GetAttributeValue("width"), 32)
	height, err := strconv.ParseFloat(svg.GetAttributeValue("height"), 32)

	// Create a group for the LEDs
	ledGroupNode := svg.CreateNode("g")

	// A copy of the "Siknäs" text paths exists in the SVG defined as a clip-path
	// we want the LEDs to be stay inside of this. So we clip the group using it.
	ledGroupNode.SetAttributeValue("clip-path", "url(#SiknasClipPath)")

	ledPositions, err := readLEDLayout(ledLayoutPath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, pos := range ledPositions {
		//log.Printf("x, y, z = %f, %f, %f\n", pos.Point[0], pos.Point[1], pos.Point[2])

		circleNode := ledGroupNode.CreateNode("circle")
		circleNode.SetAttributeValue("cx", fmt.Sprintf("%f", (pos.Point[0]*width*0.81)+(width*0.05)))
		circleNode.SetAttributeValue("cy", fmt.Sprintf("%f", (pos.Point[1]*height*0.55)+(height*0.20)))
		circleNode.SetAttributeValue("r", "10")
		circleNode.SetAttributeValue("style", "fill: rgb(255,255,255)")
	}

	//print(doc.XMLPretty())
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
