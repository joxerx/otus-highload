package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[string]*websocket.Conn
	mu      sync.Mutex
}

var WebSocketHub = &Hub{
	clients: make(map[string]*websocket.Conn),
}

// Register adds a new client connection
func (h *Hub) Register(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = conn
	log.Printf("User %s connected", userID)
}

// Unregister removes a client connection
func (h *Hub) Unregister(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conn, ok := h.clients[userID]; ok {
		conn.Close()
		delete(h.clients, userID)
		log.Printf("User %s disconnected", userID)
	}
}

// GetClient returns the WebSocket connection for a given user
func (h *Hub) GetClient(userID string) (*websocket.Conn, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	conn, ok := h.clients[userID]
	return conn, ok
}
