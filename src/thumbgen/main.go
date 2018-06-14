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
		Long: "\nThis progam is used to generate animated SVGs that are used as thumbnails for the Siknas skylt web page",
		Run:  func(c *cobra.Command, args []string) {},
	}
	rootCmd.Flags().String("host", "localhost:8080", "OPC websocket server host including port")

	rootCmd.Flags().String("logo-svg", "", "Path to Siknäs logo")
	rootCmd.MarkFlagFilename("logo-svg", "svg")

	rootCmd.Flags().String("led-layout", "layout.json", "Path to the LED layout.json")
	rootCmd.MarkFlagFilename("led-layout", "json")

	rootCmd.Flags().Duration("capture-duration", 10*time.Second, "Duration of data we should capture (in seconds). See also --max-frames")
	rootCmd.Flags().String("output", "output/", "Output path where to place the animation SVGs. A subdirectory is created for each animation")
	rootCmd.Flags().Bool("inline-svg-animation", false, "This produces an SVG with inline animation (this becomes very slow in a browser though). Normal behavior is to output each frame as a separate SVG.")

	rootCmd.Flags().String("ws-opc-path", "/ws/opc", "(Advanced) Websocket OPC path to connect to")
	rootCmd.Flags().String("ws-path", "/ws", "(Advanced) Websocket control path to connect to")
	rootCmd.Flags().BoolP("force", "f", false, "Force overwriting any existing ouput files. They are skipped by default")
	rootCmd.Flags().BoolP("list-only", "l", false, "Only list the available sketches on the server. Will not generate any SVGs")
	rootCmd.Flags().Int("max-frames", 0, "Capture up to this amount of frames. capture-duration sets max time. Set 0 to allow any frame count")
	rootCmd.Flags().Duration("frame-timeout", 3*time.Second, "The time between frames (Example: 3s, 3000ms)")
	rootCmd.Flags().StringSlice("filter", nil, "Filter to only generate thumbnails for these animations (Use --list-only to get a list of all animations)")
	rootCmd.Flags().String("siknas-background-color", "", "Force the 'Siknas' text background to be this color (should be a valid CSS style color)")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}

	viper.BindPFlags(rootCmd.Flags())

	if viper.GetBool("help") {
		os.Exit(1)
	}

	interrupt := registerSignalHandler()
	doneCtrl := make(chan struct{})

	// Read the layout file. We expect each OPC message from the server to
	// be as long as the number of LEDs in this layout.
	ledPositions, err := readLEDLayout(viper.GetString("led-layout"))
	if err != nil {
		log.Fatalln(err)
	}
	expectedMsgLen := uint16(len(ledPositions) * 3)

	wsCtrl, animationList, ctrlMessagesChan, err := ConnectControlWebsocket(doneCtrl, interrupt)
	if err != nil {
		log.Fatalln(err)
	}
	defer wsCtrl.Close()

	fmt.Println("\nAnimations:")
	for _, animation := range animationList {
		fmt.Printf("%v - %v\n", animation.Name, animation.Description)
	}
	fmt.Println()

	if viper.GetBool("list-only") {
		return
	}

	captureDuration := viper.GetDuration("capture-duration")
	outputPath := viper.GetString("output")

	err = os.MkdirAll(outputPath, 0644)
	if err != nil {
		log.Fatalln("Failed to create output dir: ", err)
	}

	filteredAnimations := viper.GetStringSlice("filter")

	/*
		Recurring bug:
		2018/06/14 18:22:35 Connecting to OPC websocket ws://192.168.1.71:8080/ws/opc...
		2018/06/14 18:22:35 Playing 'flames' for capture...
		2018/06/14 18:22:35 Capturing...
		2018/06/14 18:22:35 OPC Websocket started capturing 10s of animation (max frames: 0)
		2018/06/14 18:22:45 Finished capturing after 10s
		2018/06/14 18:22:45 OPC Websocket failed to read:  read tcp 172.17.0.2:42998->192.168.1.71:8080: use of closed network connection
		exit status 1
	*/

AnimationLoop:
	for _, animation := range animationList {
		// Skip any not in the filter list.
		if len(filteredAnimations) > 0 {
			var found = false
			for _, filter := range filteredAnimations {
				if filter == animation.Name {
					found = true
					break
				}
			}

			if !found {
				continue AnimationLoop
			}
		}
		log.Println("====================")
		log.Println(animation.Name)
		log.Println("====================")

		var targetPath string

		if viper.GetBool("inline-svg-animation") {
			targetPath = path.Join(outputPath, fmt.Sprintf("%s.%s", animation.Name, "svg"))
		} else {
			// Each frame will be a separate file, so create a directory.
			targetPath = path.Join(outputPath, animation.Name)
			// TODO: If --force remove files in the directory first.
			err = os.MkdirAll(targetPath, 0644)
			if err != nil {
				log.Fatalf("Failed to create directory '%v': %v", targetPath, err)
			}
		}

		if !viper.GetBool("force") {
			if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
				log.Printf("Skipping existing animation '%v', use --force to overwrite\n", targetPath)
				continue
			}
		}

		// stopOpcChan is used to close the connection
		// opcDoneChan returns the read messages at closing time.
		stopOpcChan := make(chan struct{})
		opcDoneChan := make(chan []OpcMessage)
		ConnectOPCWebsocket(stopOpcChan, opcDoneChan, expectedMsgLen, interrupt)

		// Selects the animation via the Control Websocket.
		log.Printf("Playing '%s' for capture...\n", animation.Name)
		ctrlMessagesChan <- animation.Name

		log.Printf("Capturing...")

	captureLoop:
		for {
			select {
			case <-stopOpcChan:
				// The websocket closed because it has captured max frames.
				break captureLoop
			case <-time.After(captureDuration):
				// The OPC Websocket will close when the stopOpcChan chan is closed.
				close(stopOpcChan)
				break captureLoop
			}
		}
		log.Printf("Finished capturing after %v\n", captureDuration)

		opcMessages := <-opcDoneChan

		if viper.GetBool("inline-svg-animation") {
			log.Printf("Saving inline SVG animation for '%s' to %v\n", animation.Name, targetPath)
			createAnimatedSVG(targetPath, opcMessages, ledPositions)
		} else {
			log.Printf("Saving SVG %d frames for '%s' to %v", len(opcMessages), animation.Name, targetPath)
			createSVGFrames(targetPath, animation.Name, opcMessages, ledPositions)
		}
	}

	// Stop any animation still running.
	ctrlMessagesChan <- ""
	time.Sleep(1 * time.Second)
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
