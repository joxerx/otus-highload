package post

import (
	"log"

	"encoding/json"
	"net/http"
	"strings"

	"otus-highload/db"
	"otus-highload/models"
	"otus-highload/redis"
	"otus-highload/utils"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil || post.Text == "" {
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

	insertPostQuery := `
		INSERT INTO posts (user_id, text) 
		VALUES ($1, $2) 
		RETURNING id, text, user_id, created_at;
	`
	err = db.MasterDB.QueryRow(insertPostQuery, authenticatedUserID, post.Text).Scan(
		&post.ID,
		&post.Text,
		&post.UserID,
		&post.CreatedAt,
	)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error creating post"})
		return
	}

	subscribers, err := db.GetSubscribers(authenticatedUserID)
	if err != nil {
		log.Printf("Failed to retrieve subscribers: %v", err)
	}

	for _, subscriber := range subscribers {
		if err := redis.EnqueueTask(redis.NotificationsStreamName, subscriber, "notify_friend", post); err != nil {
			log.Printf("Failed to enqueue notification task: %v", err)
		}
		if err := redis.EnqueueTask(redis.FeedStreamName, subscriber, "append_post", post); err != nil {
			log.Printf("Failed to enqueue task: %v", err)
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, post)
}
