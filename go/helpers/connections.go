// handlers/connections.go
package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients    = make(map[*websocket.Conn]bool) // connected clients
	broadcast  = make(chan Message)             // broadcast channel
	unregister = make(chan *websocket.Conn)     // channel to unregister clients
)

// Message struct
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
	Type      string `json:"type,omitempty"` // Type of message: join, leave, regular
}

// HandleConnections handles WebSocket connections
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v\n", err)
		return
	}
	defer ws.Close()

	// Register new client
	clients[ws] = true

	// Send join message
	msg := Message{
		Sender:  "Server",
		Content: "A new user has joined",
		Type:    "join",
	}
	broadcast <- msg

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			break
		}
		// Send the received message to the broadcast channel
		broadcast <- msg
	}

	// Close WebSocket connection and unregister client
	delete(clients, ws)
	closeConnection(ws)
}

// closeConnection closes the WebSocket connection and sends leave message
func closeConnection(ws *websocket.Conn) {
	ws.Close()
	unregister <- ws

	// Send leave message
	leaveMsg := Message{
		Sender:  "Server",
		Content: "A user has left",
		Type:    "leave",
	}
	broadcast <- leaveMsg
}
