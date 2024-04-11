package main

import (
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conns := make(map[*websocket.Conn]bool)
	defer func() {
		for conn, _ := range conns {
			conn.Close()
		}
	}()
	for i := 0; i < 1; i++ {
		conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
		if err != nil {
			log.Fatal(err)
		}
		conns[conn] = true
	}
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down ...")
}
