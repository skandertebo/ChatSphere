package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func handler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            fmt.Println(err)
            return
        }
        if err := conn.WriteMessage(messageType, p); err != nil {
            fmt.Println(err)
            return
        }
    }
}


func main() {
    http.HandleFunc("/", handler)

    log.Fatal(http.ListenAndServe(":8080", nil))
}