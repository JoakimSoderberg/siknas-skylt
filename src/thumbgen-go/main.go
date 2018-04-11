package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var rootCmd = &cobra.Command{Use: "thumbgen", Run: func(c *cobra.Command, args []string) {}}
	rootCmd.Flags().String("host", "localhost:8080", "OPC websocket server host including port")

	rootCmd.Flags().String("logo-svg", "siknas-skylt.svg", "Path to Siknäs logo")
	rootCmd.MarkFlagFilename("logo-svg", "svg")

	rootCmd.Flags().String("led-layout", "layout.json", "Path to the LED layout.json")
	rootCmd.MarkFlagFilename("led-layout", "json")

	rootCmd.Flags().Duration("capture-duration", 10*time.Second, "Duration of data we should capture (in seconds)")
	rootCmd.Flags().String("output", "output.svg", "Output filename") // TODO: change to directory
	rootCmd.Flags().String("output_path", "output/", "Output path where to place the animation SVGs")

	rootCmd.Flags().String("ws-opc-path", "/ws/opc", "OPC websocket path to connect to")
	rootCmd.Flags().String("ws-path", "/ws", "Websocket control path to connect to")
	rootCmd.Flags().Bool("force", false, "Force overwriting any existing ouput files. They are skipped by defaul")
	rootCmd.Flags().Bool("list-only", false, "Only list the available sketches on the server. Will not generate any SVGs")
	// TODO: Add option to fetch SVG from server.
	// TODO: Add option to only list the sketches from the server.

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}

	viper.BindPFlags(rootCmd.Flags())

	if viper.GetBool("help") {
		os.Exit(1)
	}

	// TODO: Connect to the "control websocket"
	// TODO: Get list of animations
	// TODO: Iterate over animations. Check if exists in output directory
	// TODO: Send a message to the server to switch to the current animation. Sleep for a while
	// TODO: Record animation for a given amount and save to disk.

	// TODO: Move this to a separate file.

	interrupt := registerSignalHandler()

	ConnectOPCWebsocket(interrupt, viper.GetString("output"))
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
