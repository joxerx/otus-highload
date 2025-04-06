package dialog

import (
	"encoding/json"
	"net/http"
	"otus-highload-messages/models"
	"otus-highload-messages/redis"
	"otus-highload-messages/utils"
	"strings"

	"fmt"

	"github.com/gorilla/mux"
)

func ListDialogHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	if userID == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "User ID is required"})
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authorization token is required"})
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")

	authenticatedUserID, err := utils.ValidateToken(token)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
		return
	}

	// Consistent chat ID
	var chatID string
	if authenticatedUserID < userID {
		chatID = fmt.Sprintf("chat:%s:%s", authenticatedUserID, userID)
	} else {
		chatID = fmt.Sprintf("chat:%s:%s", userID, authenticatedUserID)
	}

	// Get all messages
	rawMessages, err := redis.RDB.LRange(redis.CTX(), chatID, 0, -1).Result()
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to read messages from Redis"})
		return
	}

	var messages []models.Message
	for _, raw := range rawMessages {
		var msg models.Message
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			continue // skip invalid messages
		}
		messages = append(messages, msg)
	}

	utils.RespondWithJSON(w, http.StatusOK, messages)
}
