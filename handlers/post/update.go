package post

import (
	"encoding/json"
	"net/http"
	"otus-highload/db"
	"otus-highload/models"
	"otus-highload/utils"
	"strings"
)

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if post.ID == "" || post.Text == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Post ID and text are required"})
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
		SET text = $1, updated_at = clock_timestamp() 
		WHERE id = $2 AND 
		user_id = $3 AND 
		is_deleted = false
	`
	rowsAffected, err := db.ExecuteUpdateQuery(query, post.Text, post.ID, authenticatedUserID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error updating post"})
		return
	}

	if rowsAffected == 0 {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "No matching post found or already deleted"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Post successfully updated"})
}
