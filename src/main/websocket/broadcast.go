package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func SendToUser(userID string, message map[string]string) {
	// fmt.Printf("Event %s added to stream: %s with ID: %s\n", taskType, streamName, msgID)
	conn, ok := WebSocketHub.GetClient(userID)
	if !ok {
		log.Printf("User %s not connected", userID)
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to serialize message for %s: %v", userID, err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Failed to send message to %s: %v", userID, err)
		conn.Close()
		WebSocketHub.Unregister(userID)
	}
}
