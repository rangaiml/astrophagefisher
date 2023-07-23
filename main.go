package projectHailMary

import (
	"encoding/json"
	"log"
	"net/http"
)

package main

import (
"encoding/json"
"log"
"net/http"

"github.com/gorilla/websocket"
)

type Player struct {
	ID        string  `json:"id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Angle     float64 `json:"angle"`
	IsFiring  bool    `json:"isFiring"`
	IsActive  bool    `json:"isActive"`
	WebSocket *websocket.Conn
}

var players = make(map[string]*Player)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow any origin for WebSocket connections
	},
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Generate a unique ID for the player
	playerID := "player_" + getRandomID()

	// Create a new player
	newPlayer := &Player{
		ID:        playerID,
		X:         50, // Starting position X
		Y:         50, // Starting position Y
		Angle:     0,  // Starting angle
		IsFiring:  false,
		IsActive:  true,
		WebSocket: conn,
	}

	// Add the player to the players map
	players[playerID] = newPlayer

	// Notify other clients about the new player
	broadcastPlayerUpdate(newPlayer)

	// Handle messages from the player's WebSocket
	handlePlayerMessages(newPlayer)

	// Remove the player when they disconnect
	delete(players, playerID)

	// Notify other clients about the disconnected player
	newPlayer.IsActive = false
	broadcastPlayerUpdate(newPlayer)
}

func broadcastPlayerUpdate(player *Player) {
	data, err := json.Marshal(player)
	if err != nil {
		log.Println("Error marshaling player data:", err)
		return
	}

	for _, p := range players {
		if p.IsActive && p.WebSocket != nil {
			p.WebSocket.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func handlePlayerMessages(player *Player) {
	for {
		_, message, err := player.WebSocket.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		// Handle messages from the client
		// For example, updating player position, firing bullets, etc.
		// You can implement the game logic here based on the client's messages.
		// In this example, we don't implement the game logic for simplicity.
	}
}

// Utility function to generate a random ID
func getRandomID() string {
	// Implementation of getRandomID is omitted for brevity.
	// You can use a UUID library or any other method to generate a unique ID.
	// For this example, we'll return a fixed ID for simplicity.
	return "123456789"
}

func main() {
	http.HandleFunc("/ws", handleWS)
	http.Handle("/", http.FileServer(http.Dir("public"))) // Serve the frontend files (index.html, p5.js, etc.)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

