package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8081",
		Path:   "/",
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()

	// Send a message to the server
	message := []byte("popcorn*****api")
	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Error sending message to server:", err)
	}

	// Receive messages from the server
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from server:", err)
			break
		}

		log.Println("Received message from server:", string(message))
	}
}
