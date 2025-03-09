package friend

import (
	"log"
	"net/http"
	"otus-highload/db"
	"otus-highload/redis"
	"otus-highload/utils"
	"strings"
)

func SetFriendHandler(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimPrefix(r.URL.Path, "/friend/set/")
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
	if authenticatedUserID == userID {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "You're trying to add yourself"})
		return
	}
	friendExistsQuery := "SELECT id FROM users WHERE id = $1"
	rows, err := db.MasterDB.Query(friendExistsQuery, userID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error while checking friend ID"})
		return
	}
	if !rows.Next() {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Friend user ID not found"})
		return
	}
	rows.Close()

	addFriendQuery := "INSERT INTO friends (user_id, friend_id) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	if err := db.ExecuteWriteQuery(addFriendQuery, authenticatedUserID, userID); err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error adding friend"})
		return
	}
	if err := redis.EnqueueTask(authenticatedUserID, "create_feed", nil); err != nil {
		log.Printf("Failed to enqueue task: %v", err)
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Friend successfully added"})
}
