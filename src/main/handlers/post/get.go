package post

import (
	"net/http"
	"otus-highload/db"
	"otus-highload/models"
	"otus-highload/utils"
	"strings"
)

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := strings.TrimPrefix(r.URL.Path, "/post/get/")
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
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	query := `
		SELECT id, text 
		FROM posts 
		WHERE id = $1 AND 
		(
			(user_id = $2 AND is_deleted = false) OR 
			(visibility = 'public' AND is_deleted = false)
		)
	`

	var post models.Post
	rows, err := db.ExecuteReadQuery(query, postID, authenticatedUserID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error executing read query"})
		return
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error reading post data"})
			return
		}
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		return
	}

	err = rows.Scan(&post.ID, &post.Text)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error parsing post data"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, post)
}
