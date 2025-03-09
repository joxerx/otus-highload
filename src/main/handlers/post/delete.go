package post

import (
	"log"
	"net/http"
	"otus-highload/db"
	"otus-highload/redis"
	"otus-highload/utils"
	"strings"
)

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := strings.TrimPrefix(r.URL.Path, "/post/delete/")
	if postID == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Post ID is required"})
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

	query := `
		UPDATE posts 
		SET is_deleted = true, updated_at = clock_timestamp() 
		WHERE id = $1 AND 
		user_id = $2 AND 
		is_deleted = false
	`
	rowsAffected, err := db.ExecuteUpdateQuery(query, postID, authenticatedUserID)

	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error deleting post"})
		return
	}
	if rowsAffected == 0 {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "No matching post found or already deleted"})
		return
	}

	subscribers, err := db.GetSubscribers(authenticatedUserID)
	if err != nil {
		log.Fatalf("Failed to retrieve subscribers: %v", err)
	}
	for _, subscriber := range subscribers {
		if err := redis.EnqueueTask(subscriber, "delete_post", map[string]string{"postID": postID}); err != nil {
			log.Printf("Failed to enqueue task: %v", err)
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Post successfully deleted"})

}
