package post

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"otus-highload/db"
	"otus-highload/models"
	"otus-highload/redis"
	"otus-highload/utils"
	"strings"
)

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil || post.ID == "" || post.Text == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
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

	updateQuery := `
		UPDATE posts 
		SET text = $1, updated_at = clock_timestamp() 
		WHERE id = $2 AND user_id = $3 AND is_deleted = false
		RETURNING id, text, user_id, created_at;
	`

	var updatedPost models.Post
	err = db.MasterDB.QueryRow(updateQuery, post.Text, post.ID, authenticatedUserID).Scan(
		&updatedPost.ID,
		&updatedPost.Text,
		&updatedPost.UserID,
		&updatedPost.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "No matching post found or already deleted"})
		} else {
			log.Printf("Failed to update post: %v", err)
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error updating post"})
		}
		return
	}

	subscribers, err := db.GetSubscribers(authenticatedUserID)
	if err != nil {
		log.Printf("Failed to retrieve subscribers: %v", err)
	}

	for _, subscriber := range subscribers {
		if err := redis.EnqueueTask(subscriber, "update_post", updatedPost); err != nil {
			log.Printf("Failed to enqueue task: %v", err)
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedPost)
}
