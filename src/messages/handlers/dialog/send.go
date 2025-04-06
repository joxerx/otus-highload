package dialog

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math/rand"
	"net/http"
	"otus-highload-messages/models"
	"otus-highload-messages/redis"
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

	var chatID string
	if authenticatedUserID < userID {
		chatID = fmt.Sprintf("chat:%s:%s", authenticatedUserID, userID)
	} else {
		chatID = fmt.Sprintf("chat:%s:%s", userID, authenticatedUserID)
	}

	createdAt := time.Now()

	// Generate a unique message ID
	now := time.Now().Unix()         // current Unix time in seconds
	randPart := rand.Intn(1_000_000) // random 6-digit number
	h := fnv.New64()
	h.Write([]byte(authenticatedUserID))
	hashInt := int64(h.Sum64())
	senderNum := hashInt % 1_000_000
	messageID := int((now%1_000_000_000)*1_000_000+int64(randPart)) + int(senderNum)

	msg := models.Message{
		ID:        messageID,
		Sender:    authenticatedUserID,
		Recipient: userID,
		Content:   msgReq.Text,
		CreatedAt: createdAt.String(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to encode message"})
		return
	}

	if err := redis.RDB.RPush(redis.CTX(), chatID, data).Err(); err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to store message in Redis"})
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message":    "Message sent successfully",
		"message_id": messageID,
	})
}
