package friend

import (
	"log"
	"net/http"
	"otus-highload/db"
	"otus-highload/redis"
	"otus-highload/utils"
	"strings"
)

func DeleteFriendHandler(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimPrefix(r.URL.Path, "/friend/delete/")
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

	deleteFriendQuery := "DELETE FROM friends WHERE user_id = $1 AND friend_id = $2"
	rowsAffected, err := db.ExecuteUpdateQuery(deleteFriendQuery, authenticatedUserID, userID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error deleting friend"})
		return
	}

	if rowsAffected == 0 {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "No friend record found to delete"})
		return
	}

	if err := redis.EnqueueTask(authenticatedUserID, "create_feed", nil); err != nil {
		log.Printf("Failed to enqueue task: %v", err)
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Friend successfully deleted"})
}
