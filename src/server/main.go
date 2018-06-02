package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO: Add more constants for hard coded things
const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func main() {
	var rootCmd = &cobra.Command{Use: "siknas-skylt", Run: func(c *cobra.Command, args []string) {}}
	rootCmd.Flags().Int("port", 8080, "The port the webserver should listen on")
	rootCmd.Flags().Int("opc-listen-port", 7890, "The port to listen for OpenPixelControl protocol data on")
	rootCmd.Flags().String("static-path", "static/siknas-skylt", "Path to the static files")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}

	if viper.GetBool("help") {
		os.Exit(1)
	}

	// Allow viper to parse the command line flags also.
	viper.BindPFlags(rootCmd.Flags())
	//viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))

	viper.SetConfigName("siknas")
	viper.AddConfigPath("/etc/siknas/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Printf("No config file: %s\n", err)
	} else {
		log.Printf("Found config file: %s\n", viper.ConfigFileUsed())
	}

	// TODO: Move the config parsing here and pass as arguments.

	// Broadcasts control panel messages to all connected websocket clients.
	controlPanelBroadcaster := NewControlPanelBroadcaster()

	// Broadcasts all OPC messages coming from the animation process
	// to both Websocket clients and the display(s) we are connecting to.
	opcBroadcaster := NewOpcBroadcaster()

	// Used to process the current state of what animation is playing
	// to all connected websocket clients.
	opcProcessManagerBroadcaster := NewOpcProcessManagerBroadcaster()

	// Handles the animation processes that produces the incoming OPC network traffic.
	opcProcessManager, err := NewOpcProcessManager(opcProcessManagerBroadcaster)
	if err != nil {
		log.Fatalln("Failed to create OPC process manager:", err)
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/debug", WsDebugHandler)
	router.HandleFunc("/debug/opc", WsDebugOpcHandler)
	router.HandleFunc("/debug/control_panel", WsDebugControlPanelHandler)

	// Websocket handlers.
	router.HandleFunc("/ws", WsHandler(controlPanelBroadcaster, opcProcessManager, opcBroadcaster))
	router.HandleFunc("/ws/opc", OpcWsHandler(opcBroadcaster))
	router.HandleFunc("/ws/control_panel", ControlPanelWsHandler(controlPanelBroadcaster, opcProcessManager))

	// Must be last so we don't shadow the other routes.
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(viper.GetString("static-path"))))

	// TODO: Move to function
	// Add OPC servers we should send to.
	opcServers := viper.GetStringMap("opc-servers")

	for name := range opcServers {
		// TODO: Unmarshal into struct instead. And pass on to ws handler so clients can list these
		opcHost := viper.GetString(fmt.Sprintf("opc-servers.%v.host", name))
		opcPort := viper.GetString(fmt.Sprintf("opc-servers.%v.port", name))
		// TODO: Enable setting connection backoff in config

		opcAddr := fmt.Sprintf("%v:%v", opcHost, opcPort)
		log.Println(opcAddr)
		go RunOpcClient("tcp", opcAddr, 30*time.Second, opcBroadcaster)
	}

	// OPC Proxy.
	opcListenPort := viper.GetInt("opc-listen-port")
	go RunOPCProxy("tcp", fmt.Sprintf(":%v", opcListenPort), opcBroadcaster)

	port := viper.GetInt("port")
	log.Printf("Starting Siknas-skylt webserver on %v...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}
