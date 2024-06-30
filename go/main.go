package main

import (
	"log"
	"net/http"

	// Import gorilla/handlers for CORS
	"github.com/gorilla/handlers" // Import gorilla/handlers for CORS

	helpers "performance.go/helpers" // Update with your actual package path
)

func main() {
	// Configure websocket route
	http.HandleFunc("/ws", helpers.HandleConnections)

	// Start listening for incoming chat messages
	go helpers.HandleMessages()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// CORS middleware
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Change "*" to your specific allowed origins
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Create a new handler with CORS middleware
	handlerWithCORS := cors(http.DefaultServeMux)

	// Start the server on localhost:8080
	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", handlerWithCORS)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
