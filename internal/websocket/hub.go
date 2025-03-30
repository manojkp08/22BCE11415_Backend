package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
}

var Clients = make(map[string]*Client)

func BroadcastToUser(userID string, message interface{}) {
	if client, ok := Clients[userID]; ok {
		if err := client.Conn.WriteJSON(message); err != nil {
			log.Printf("WebSocket error: %v", err)
			client.Conn.Close()
			delete(Clients, userID)
		}
	}
}
