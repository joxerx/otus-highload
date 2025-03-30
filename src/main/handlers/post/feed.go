package post

import (
	"net/http"
	"otus-highload/models"
	"otus-highload/redis"
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
	feed := []models.Post{}

	if offset < 1000 {
		redisPosts, err := redis.GetCachedPosts(authenticatedUserID, offset, limit)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve feed from cache"})
			return
		}
		feed = append(feed, redisPosts...)
	}

	if len(feed) < limit {
		missing := limit - len(feed)
		dbPosts, err := redis.FetchOlderPostsFromDB(authenticatedUserID, offset+len(feed), missing)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve feed from DB"})
			return
		}
		feed = append(feed, dbPosts...)
	}

	utils.RespondWithJSON(w, http.StatusOK, feed)
}
