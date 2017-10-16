package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
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

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	// TODO: Host aurelia webpage
}

func main() {
	var rootCmd = &cobra.Command{Use: "siknas-skylt", Run: func(c *cobra.Command, args []string) {}}
	rootCmd.Flags().Int("port", 8080, "The port the webserver should listen on")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}

	// Allow viper to parse the command line flags also.
	viper.BindPFlags(rootCmd.Flags())
	//viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))

	viper.SetConfigName("siknas")
	viper.AddConfigPath("/etc/siknas/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	// Broadcast channel for control panel.
	controlPanelBroadcaster := NewControlPanelBroadcaster()
	opcBroadcaster := NewOpcBroadcaster()

	opcProcessManager, err := NewOpcProcessManager()
	if err != nil {
		log.Fatalln("Failed to create OPC process manager:", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/debug", WsDebugHandler)
	router.HandleFunc("/debug/opc", WsDebugOpcHandler)
	router.HandleFunc("/debug/control_panel", WsDebugControlPanelHandler)

	// Websocket handlers.
	router.HandleFunc("/ws", WsHandler(controlPanelBroadcaster, opcProcessManager))
	router.HandleFunc("/ws/opc", OpcWsHandler(opcBroadcaster))
	router.HandleFunc("/ws/control_panel", ControlPanelWsHandler(controlPanelBroadcaster))

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
	go RunOPCProxy("tcp", ":7890", opcBroadcaster)

	port := viper.GetInt("port")
	log.Printf("Starting Siknas-skylt webserver on %v...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}
