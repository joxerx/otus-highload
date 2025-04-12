package counter

import (
	"fmt"
	"net/http"
	"otus-highload-counter/redis"
	"otus-highload-counter/utils"

	"strings"

	"github.com/gorilla/mux"
)

func DialogCounterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dialogID := vars["dialogId"]
	if dialogID == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "User ID and Dialog ID are required"})
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

	key := fmt.Sprintf("unreadCount:%s:%s", dialogID, authenticatedUserID)
	count, err := redis.RDB.Get(redis.CTX(), key).Int()
	if err != nil {
		utils.RespondWithJSON(w, http.StatusOK, map[string]int{"count": 0})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]int{"count": count})
}
