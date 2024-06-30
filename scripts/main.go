package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	numClients        int
	messagesPerClient int
	port              int
	server            string
	wg                sync.WaitGroup
	mutex             sync.Mutex
	results           = make(map[string]*Matrix)
)

type Message struct {
	Sender  string    `json:"sender"`
	ID      string    `json:"id"`
	Content string    `json:"content"`
	IAT     time.Time `json:"iat"`
	Type    string    `json:"type"`
}

type Timing struct {
	MessageID string  `json:"message_id"`
	Time      float64 `json:"time"` // Store time in milliseconds as float64
	ClientID  string  `json:"client_id"`
}

type Matrix struct {
	ClientID        string   `json:"client_id"`
	MessagesTimings []Timing `json:"messages_timings"`
	Avg             float64  `json:"avg"`
	Min             float64  `json:"min"`
	Max             float64  `json:"max"`
}

func init() {
	flag.IntVar(&numClients, "clients", 10, "Number of clients to simulate")
	flag.IntVar(&messagesPerClient, "messages", 100, "Number of messages per client")
	flag.IntVar(&port, "port", 8080, "WebSocket server's port")
	flag.StringVar(&server, "server", "typescript", "WebSocket server's port")
	flag.Parse()
}

func main() {
	wg.Add(numClients)

	for i := 0; i < numClients; i++ {
		clientID := fmt.Sprintf("Client_%d", i+1)
		simulateClient(clientID)
	}

	wg.Wait()

	calculateOverallAverages()
}

func simulateClient(clientID string) {
	defer wg.Done()

	url := fmt.Sprintf("ws://localhost:%d/ws", port)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Printf("Client %s: Error connecting to WebSocket: %v\n", clientID, err)
		return
	}
	defer conn.Close()

	log.Printf("Client %s connected.\n", clientID)

	for j := 0; j < messagesPerClient; j++ {
		message := Message{
			Sender:  clientID,
			ID:      uuid.New().String(),
			Content: fmt.Sprintf("Message %d from %s", j+1, clientID),
			IAT:     time.Now(),
			Type:    "message",
		}

		// Retry sending the message up to 3 times on failure
		for retry := 0; retry < 3; retry++ {
			err := conn.WriteJSON(message)
			if err != nil {
				log.Printf("Client %s: Error sending message (retry %d): %v\n", clientID, retry+1, err)

				// Attempt to reconnect if the connection is closed unexpectedly
				conn.Close()
				conn, _, err = websocket.DefaultDialer.Dial(url, nil)
				if err != nil {
					log.Printf("Client %s: Error reconnecting: %v\n", clientID, err)
					time.Sleep(1 * time.Second) // Wait before retrying
					continue
				}
				continue
			}

			// Message sent successfully
			break
		}

		// Handle message sending failures after retries
		if err != nil {
			log.Printf("Client %s: Failed to send message after retries: %v\n", clientID, err)
			continue // Skip recording timings if message sending fails
		}

		// Wait for response message
		var response Message
		err = conn.ReadJSON(&response)
		if err != nil {
			log.Printf("Client %s: Error reading response: %v\n", clientID, err)

			// Attempt to reconnect if the connection is closed unexpectedly
			conn.Close()
			conn, _, err = websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				log.Printf("Client %s: Error reconnecting: %v\n", clientID, err)
				continue
			}
			continue
		}

		// Calculate elapsed time in milliseconds
		timeTaken := float64(time.Since(message.IAT).Nanoseconds()) / float64(time.Millisecond)

		// Record timing information
		timing := Timing{
			MessageID: response.ID,
			Time:      timeTaken,
			ClientID:  clientID,
		}

		mutex.Lock()
		if _, ok := results[clientID]; !ok {
			results[clientID] = &Matrix{
				ClientID:        clientID,
				MessagesTimings: make([]Timing, 0, messagesPerClient),
			}
		}
		results[clientID].MessagesTimings = append(results[clientID].MessagesTimings, timing)
		mutex.Unlock()
	}

	log.Printf("Client %s: All messages sent and received.\n", clientID)
}

func (m *Matrix) MarshalJSON() ([]byte, error) {
	type Alias Matrix
	return json.Marshal(&struct {
		Avg             float64  `json:"avg"`
		Min             float64  `json:"min"`
		Max             float64  `json:"max"`
		MessagesTimings []Timing `json:"messages_timings"`
		*Alias
	}{
		Avg:             m.Avg,
		Min:             m.Min,
		Max:             m.Max,
		MessagesTimings: m.MessagesTimings,
		Alias:           (*Alias)(m),
	})
}

func calculateOverallAverages() {
	log.Println("Calculating overall averages...")

	for clientID, matrix := range results {
		var sum float64
		min := float64(^uint64(0) >> 1)
		max := -1.0

		for _, timing := range matrix.MessagesTimings {
			sum += timing.Time
			if timing.Time < min {
				min = timing.Time
			}
			if timing.Time > max {
				max = timing.Time
			}
		}

		if len(matrix.MessagesTimings) > 0 {
			avg := sum / float64(len(matrix.MessagesTimings))
			matrix.Avg = avg
			matrix.Min = min
			matrix.Max = max
		} else {
			matrix.Avg = 0
			matrix.Min = 0
			matrix.Max = 0
		}

		results[clientID] = matrix
	}

	outputJSON, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling results to JSON: %v\n", err)
	}

	outputDir := "metrics"
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Error creating output directory: %v\n", err)
	}

	outputFile := fmt.Sprintf("results-client-%s-%d-messages-%d.json", server, numClients, messagesPerClient)
	outputPath := path.Join(outputDir, outputFile)

	err = os.WriteFile(outputPath, outputJSON, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON to file: %v\n", err)
	}

	log.Printf("Results saved to %s\n", outputPath)
}
