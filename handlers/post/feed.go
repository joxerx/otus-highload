package post

import (
	"net/http"
	"otus-highload/db"
	"otus-highload/models"
	"otus-highload/utils"
	"strconv"
	"strings"
)

func FeedHandler(w http.ResponseWriter, r *http.Request) {
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

	queryParams := r.URL.Query()
	limit := 1000
	offset := 0

	if l := queryParams.Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if o := queryParams.Get("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	query := `
		SELECT p.id, p.text
		FROM posts p
		JOIN friends f ON p.user_id = f.friend_id
		WHERE f.user_id = $1
		  AND p.visibility = 'public'
		  AND p.is_deleted = false
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.ExecuteReadQuery(query, authenticatedUserID, limit, offset)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving feed"})
		return
	}
	defer rows.Close()

	var feed []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Text); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error processing feed"})
			return
		}
		feed = append(feed, post)
	}

	utils.RespondWithJSON(w, http.StatusOK, feed)
}
