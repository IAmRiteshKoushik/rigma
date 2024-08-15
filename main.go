package main

import (
    "time"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
    // Alternately you can maintain an array of strings which can maintain 
    // connection strings to each server
    conns map[*websocket.Conn]bool
}

func NewServer() *Server {
    return &Server{
        conns: make(map[*websocket.Conn]bool),
    }
}

func (s *Server) handleWSOrderbook(ws *websocket.Conn) {
    fmt.Println("New incoming connection from client to orderbook feed:", ws.RemoteAddr())
    for {
        // simulate LIVE data from orderbook
        payload := fmt.Sprintf("Orderbook data -> %d\n", time.Now().UnixNano())
        ws.Write([]byte(payload))
        time.Sleep(time.Second * 2)
    }
}

func (s *Server) handleWS(ws *websocket.Conn){
    fmt.Println("New incoming connection from client:", ws.RemoteAddr())

    // maps in golang are not concurrent safe, so you should use some mutex
    // to avoid race-conditions
    s.conns[ws] = true
    s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
    buf := make([]byte, 1024)
    for {
        n, err := ws.Read(buf)
        if err != nil {
            // The connection on the other side has closed itself
            // then you can break the connection
            if err == io.EOF {
                break
            }
            // Otherwise you can persist
            fmt.Println("Read error:", err)
            continue
        }
        msg := buf[:n]
        // fmt.Println(string(msg))

        // Whatever message is received by the websocket-server, it is 
        // immediately broadcasted to all other client connections
        s.broadcast(msg)
    }
}

func (s *Server) broadcast(b []byte) {
    // Looping through all the active connections
    for ws := range s.conns {
        go func(ws *websocket.Conn) {
            if _, err := ws.Write(b); err != nil {
                fmt.Println("Write error:", err)
            }
        }(ws)
    }
}

func main() {
    server := NewServer()
    http.Handle("/ws", websocket.Handler(server.handleWS))
    http.Handle("/orderbook", websocket.Handler(server.handleWSOrderbook))
    http.ListenAndServe(":3000", nil)
}
