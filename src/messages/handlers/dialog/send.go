package dialog

import (
	"encoding/json"
	"net/http"
	"otus-highload-messages/db"
	"otus-highload-messages/models"
	"otus-highload-messages/utils"

	"strings"
	"time"

	"github.com/gorilla/mux"
)

func SendDialogHandler(w http.ResponseWriter, r *http.Request) {
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

	var msgReq models.MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&msgReq); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if msgReq.Text == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Text field is required"})
		return
	}

	var shardKey string
	if authenticatedUserID < userID {
		shardKey = authenticatedUserID + ":" + userID
	} else {
		shardKey = userID + ":" + authenticatedUserID
	}

	query := `
		INSERT INTO messages (sender, recipient, content, created_at, shard_key)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	var messageID string
	messageID, err = db.ExecuteInsertQuery(query, authenticatedUserID, userID, msgReq.Text, time.Now(), shardKey)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to send message"})
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message":    "Message sent successfully",
		"message_id": messageID,
	})
}
