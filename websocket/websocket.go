package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"rest-api-gorilla/models"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.Message)
var upgrader = websocket.Upgrader{}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Failed to upgrade connection: %v", err)
	}

	defer ws.Close()

	clients[ws] = true

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		log.Println("Broadcast:", clients, broadcast)
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Fatalf("Error sending message: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func BroadcastMessage(msg models.Message) {
	broadcast <- msg
}
