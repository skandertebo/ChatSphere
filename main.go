package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
    Conn *websocket.Conn
    Send chan []byte
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

var (
    connections []*Connection
    connMutex   sync.Mutex
)

func broadcastMessage(msg []byte, sender *Connection) {
    connMutex.Lock()
    defer connMutex.Unlock()

    for _, conn := range connections {
        if conn != sender {
            select {
            case conn.Send <- msg:
            default:
                continue
            }
        }
    }
}

func handleConnection(conn *websocket.Conn) {
    c := &Connection{
        Conn: conn,
        Send: make(chan []byte, 256),
    }

    connMutex.Lock()
    connections = append(connections, c)
    connMutex.Unlock()

    go func() {
        defer func() {
            conn.Close()
            connMutex.Lock()
            for i, connection := range connections {
                if connection == c {
                    connections = append(connections[:i], connections[i+1:]...)
                    break
                }
            }
            connMutex.Unlock()
        }()

        for {
            _, msg, err := conn.ReadMessage()
            if err != nil {
                return
            }
            broadcastMessage(msg, c)
        }
    }()

    for msg := range c.Send {
        if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            return
        }
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade failed:", err)
        return
    }
    handleConnection(conn)
}

func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}