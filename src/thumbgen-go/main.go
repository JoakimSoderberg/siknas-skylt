package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:  "thumbgen",
		Long: "This progam is used to generate animated SVGs that are used as thumbnails for the web page",
		Run:  func(c *cobra.Command, args []string) {},
	}
	rootCmd.Flags().String("host", "localhost:8080", "OPC websocket server host including port")

	rootCmd.Flags().String("logo-svg", "siknas-skylt.svg", "Path to Siknäs logo")
	rootCmd.MarkFlagFilename("logo-svg", "svg")

	rootCmd.Flags().String("led-layout", "layout.json", "Path to the LED layout.json")
	rootCmd.MarkFlagFilename("led-layout", "json")

	rootCmd.Flags().Duration("capture-duration", 10*time.Second, "Duration of data we should capture (in seconds)")
	//rootCmd.Flags().String("output", "output.svg", "Output filename") // TODO: change to directory
	rootCmd.Flags().String("output", "output/", "Output path where to place the animation SVGs")

	rootCmd.Flags().String("ws-opc-path", "/ws/opc", "Websocket OPC path to connect to")
	rootCmd.Flags().String("ws-path", "/ws", "Websocket control path to connect to")
	rootCmd.Flags().BoolP("force", "f", false, "Force overwriting any existing ouput files. They are skipped by defaul")
	rootCmd.Flags().BoolP("list-only", "l", false, "Only list the available sketches on the server. Will not generate any SVGs")
	// TODO: Add option to fetch SVG from server.

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}

	viper.BindPFlags(rootCmd.Flags())

	if viper.GetBool("help") {
		os.Exit(1)
	}

	// TODO: Iterate over animations. Check if exists in output directory
	// TODO: Send a message to the server to switch to the current animation. Sleep for a while
	// TODO: Record animation for a given amount and save to disk.

	interrupt := registerSignalHandler()
	doneCtrl := make(chan struct{})

	wsCtrl, animationList, ctrlMessagesChan, err := ConnectControlWebsocket(doneCtrl, interrupt)
	if err != nil {
		log.Fatalln(err)
	}
	defer wsCtrl.Close()

	if viper.GetBool("list-only") {
		fmt.Println("\nAnimations:")
		for _, animation := range animationList {
			fmt.Printf("%v - %v\n", animation.Name, animation.Description)
		}
		return
	}

	captureDuration := viper.GetDuration("capture-duration")
	outputPath := viper.GetString("output")

	err = os.MkdirAll(outputPath, 0644)
	if err != nil {
		log.Fatalln("Failed to create output dir: ", err)
	}

	for _, animation := range animationList {
		doneOPC := make(chan struct{})

		ws := ConnectOPCWebsocket(doneOPC, interrupt)

		// Selects the animation via the Control Websocket.
		log.Printf("Selecting '%s' for capture...\n", animation.Name)
		time.Sleep(3 * time.Second)
		ctrlMessagesChan <- animation.Name
		time.Sleep(3 * time.Second)
		log.Printf("Capturing...")

		time.Sleep(captureDuration)
		log.Printf("Finished capturing after %v\n", captureDuration)
		close(doneOPC)
		ws.Close()

		filePath := path.Join(outputPath, fmt.Sprintf("%s.%s", animation.Name, "svg"))

		if !viper.GetBool("force") {
			if _, err := os.Stat(filePath); err != nil {
				if !os.IsNotExist(err) {
					log.Println("Skipping existing file: ", filePath)
					continue
				}
			}
		}

		createOutputSVG(filePath)
	}
}

// registerSignalHandler handles interrupt signals.
func registerSignalHandler() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// TODO: Take a channel as argument that we signal close on.

	go func() {
		for range sigChan {
			log.Println("\nReceived signal, exiting...")
			os.Exit(0)
		}
	}()

	return sigChan
}
