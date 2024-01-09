package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
    connections []*websocket.Conn
    connMutex   sync.Mutex
    upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }
)

func listenToMessages() {
    for {
        connMutex.Lock()
        for _, conn := range connections {
            _, msg, err := conn.ReadMessage()
            if err != nil {
                fmt.Println(err)
                continue
            }
            fmt.Println(string(msg))
            // distribute message to all connections except the sender
            for _, connInner := range connections {
                if conn != connInner {
                    if err = connInner.WriteMessage(websocket.TextMessage, msg); err != nil {
                        fmt.Println(err)
                        continue
                    }
                }
            }
        }
        connMutex.Unlock()
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    connections = append(connections, conn)
}

func main() {
    http.HandleFunc("/", handler)
    go listenToMessages()
    log.Fatal(http.ListenAndServe(":8080", nil))
}