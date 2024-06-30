// handlers/messages.go
package handlers

import (
	"log"
)

// HandleMessages handles broadcasting messages to clients
func HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		// Send it out to every client that is currently connected
		for client := range clients {
			// Check if the client is not the sender
			if client.RemoteAddr().String() != msg.Sender {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
