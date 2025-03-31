package websocket

import (
	"log"
	"net/http"
	"strings"

	"otus-highload/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")

	userID, err := utils.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	WebSocketHub.Register(userID, conn)

	go func() {
		defer WebSocketHub.Unregister(userID)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("User %s disconnected", userID)
				return
			}
		}
	}()
}
