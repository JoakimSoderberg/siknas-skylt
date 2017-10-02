package main

import (
    "fmt"
    "net/http"

    "github.com/gobwas/ws"
    "github.com/gobwas/ws/wsutil"
)

func ws_listener() {
}

func main() {
    fmt.Printf("Starting")
    http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        conn, _, _, err := ws.UpgradeHTTP(r, w, nil)
        if err != nil {
            fmt.Printf("Failed to Ugrade websocket connection")
            conn.Close()
            return
        }

        go func() {
            defer conn.Close()

            fmt.Printf("New Websocket Client\n")

            for {
                msg, op, err := wsutil.ReadClientData(conn)
                if err != nil {
                    // handle error
                }
                err = wsutil.WriteServerMessage(conn, op, msg)
                if err != nil {
                    // handle error
                }
            }
        }()
    }))
}

