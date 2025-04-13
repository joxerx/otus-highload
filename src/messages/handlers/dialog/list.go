package dialog

import (
	"encoding/json"
	"net/http"
	"otus-highload-messages/metrics"
	"otus-highload-messages/models"
	"otus-highload-messages/redis"
	"otus-highload-messages/utils"
	"strings"
	"time"

	"fmt"

	"github.com/gorilla/mux"
)

func ListDialogHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	handlerName := "listDialog"
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues(handlerName, "GET").Observe(duration)
	}()

	vars := mux.Vars(r)
	userID := vars["userId"]
	if userID == "" {
		metrics.RequestCounter.WithLabelValues(handlerName, "GET", "400").Inc()
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "User ID is required"})
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		metrics.RequestCounter.WithLabelValues(handlerName, "GET", "401").Inc()
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authorization token is required"})
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")

	authenticatedUserID, err := utils.ValidateToken(token)
	if err != nil {
		metrics.RequestCounter.WithLabelValues(handlerName, "GET", "401").Inc()
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
		return
	}

	var chatID string
	if authenticatedUserID < userID {
		chatID = fmt.Sprintf("chat:%s:%s", authenticatedUserID, userID)
	} else {
		chatID = fmt.Sprintf("chat:%s:%s", userID, authenticatedUserID)
	}

	rawMessages, err := redis.RDB.LRange(redis.CTX(), chatID, 0, -1).Result()
	if err != nil {
		metrics.RequestCounter.WithLabelValues(handlerName, "GET", "500").Inc()
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to read messages from Redis"})
		return
	}

	var messages []models.Message
	for _, raw := range rawMessages {
		var msg models.Message
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	utils.RespondWithJSON(w, http.StatusOK, messages)
	metrics.RequestCounter.WithLabelValues(handlerName, "GET", "200").Inc()
}
