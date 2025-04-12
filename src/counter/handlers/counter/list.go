package counter

import (
	"fmt"
	"net/http"
	"otus-highload-counter/redis"
	"otus-highload-counter/utils"

	"strings"
)

func ListCounterHandler(w http.ResponseWriter, r *http.Request) {
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

	pattern := fmt.Sprintf("unreadCount:*:%s", authenticatedUserID)
	keys, err := redis.RDB.Keys(redis.CTX(), pattern).Result()
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to read counters"})
		return
	}

	result := make(map[string]int)
	for _, key := range keys {
		count, err := redis.RDB.Get(redis.CTX(), key).Int()
		if err == nil {
			parts := strings.Split(key, ":")
			if len(parts) == 3 {
				dialogID := parts[1]
				result[dialogID] = count
			}
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, result)
}
