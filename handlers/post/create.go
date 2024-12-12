package post

import (
	"encoding/json"
	"net/http"
	"otus-highload/db"
	"otus-highload/models"
	"otus-highload/utils"
	"strings"
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

	insertPostQuery := "INSERT INTO posts (user_id, text) VALUES ($1, $2) RETURNING id"
	var postID string
	err = db.MasterDB.QueryRow(insertPostQuery, authenticatedUserID, post.Text).Scan(&postID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error creating post"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"ID": postID})
}
