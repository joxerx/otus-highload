package dialog

import (
	"net/http"
	"otus-highload-messages/db"
	"otus-highload-messages/models"
	"otus-highload-messages/utils"
	"strings"

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
	var shardKey string
	if authenticatedUserID < userID {
		shardKey = authenticatedUserID + ":" + userID
	} else {
		shardKey = userID + ":" + authenticatedUserID
	}

	query := `
		SELECT id, sender, recipient, content, created_at
		FROM messages
		WHERE shard_key = $1
		ORDER BY created_at ASC;
	`

	rows, err := db.ExecuteReadQuery(query, shardKey)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch messages"})
		return
	}
	defer rows.Close()

	// Collect messages
	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.Sender, &msg.Recipient, &msg.Content, &msg.CreatedAt); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to parse messages"})
			return
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error reading rows"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, messages)
}
