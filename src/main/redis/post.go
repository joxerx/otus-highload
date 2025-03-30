package redis

import (
	"encoding/json"
	"log"
	"otus-highload/db"
	"otus-highload/models"
	"time"

	"github.com/redis/go-redis/v9"
)

const feedTTLHours = 1

func AppendPostToFeed(userID string, post models.Post) error {
	key := "feed:" + userID

	postJSON, err := json.Marshal(post)
	if err != nil {
		return err
	}

	score := float64(post.CreatedAt.Unix())
	if err := RDB.ZAdd(CTX(), key, redis.Z{
		Score:  score,
		Member: postJSON,
	}).Err(); err != nil {
		return err
	}

	if err := RDB.ZRemRangeByRank(CTX(), key, 0, -1001).Err(); err != nil {
		return err
	}
	if err := RDB.Expire(CTX(), key, feedTTLHours*time.Hour).Err(); err != nil {
		return err
	}
	return nil
}

func CreateWholeFeed(userID string) error {
	posts, err := FetchOlderPostsFromDB(userID, 0, 1000)
	if err != nil {
		return err
	}

	key := "feed:" + userID

	pipe := RDB.TxPipeline()

	pipe.Del(CTX(), key)

	for _, post := range posts {
		postJSON, err := json.Marshal(post)
		if err != nil {
			return err
		}
		score := float64(post.CreatedAt.Unix())
		pipe.ZAdd(CTX(), key, redis.Z{
			Score:  score,
			Member: postJSON,
		})
	}

	_, err = pipe.Exec(CTX())
	if err := RDB.Expire(CTX(), key, feedTTLHours*time.Hour).Err(); err != nil {
		return err
	}
	return err
}

func DeletePostFromFeed(userID string, postID string) error {
	key := "feed:" + userID

	postsJSON, err := RDB.ZRange(CTX(), key, 0, -1).Result()
	if err != nil {
		return err
	}

	pipe := RDB.TxPipeline()

	for _, postJSON := range postsJSON {
		var post models.Post
		if err := json.Unmarshal([]byte(postJSON), &post); err != nil {
			continue
		}

		if post.ID == postID {
			pipe.ZRem(CTX(), key, postJSON)
			break
		}
	}
	if err := RDB.Expire(CTX(), key, feedTTLHours*time.Hour).Err(); err != nil {
		return err
	}
	_, err = pipe.Exec(CTX())
	return err
}

func UpdatePostInFeed(userID string, updatedPost models.Post) error {
	key := "feed:" + userID

	postsJSON, err := RDB.ZRange(CTX(), key, 0, -1).Result()
	if err != nil {
		return err
	}

	pipe := RDB.TxPipeline()

	for _, postJSON := range postsJSON {
		var post models.Post
		if err := json.Unmarshal([]byte(postJSON), &post); err != nil {
			continue
		}

		if post.ID == updatedPost.ID {

			pipe.ZRem(CTX(), key, postJSON)

			updatedPostJSON, err := json.Marshal(updatedPost)
			if err != nil {
				return err
			}
			score := float64(updatedPost.CreatedAt.Unix())
			pipe.ZAdd(CTX(), key, redis.Z{
				Score:  score,
				Member: updatedPostJSON,
			})
			break
		}
	}
	if err := RDB.Expire(CTX(), key, feedTTLHours*time.Hour).Err(); err != nil {
		return err
	}
	_, err = pipe.Exec(CTX())
	return err
}

func GetCachedPosts(userID string, offset, limit int) ([]models.Post, error) {
	key := "feed:" + userID
	posts := []models.Post{}

	postJSONs, err := RDB.ZRevRange(CTX(), key, 0, 999).Result()
	if err != nil {
		return nil, err
	}

	if len(postJSONs) > offset {
		for _, postJSON := range postJSONs[offset:min(len(postJSONs), offset+limit, 1000)] {
			var post models.Post
			if err := json.Unmarshal([]byte(postJSON), &post); err != nil {
				continue
			}
			posts = append(posts, post)
		}
	}

	if len(posts) == 0 {
		if err := EnqueueTask(FeedStreamName, userID, "create_feed", nil); err != nil {
			log.Printf("Failed to enqueue task: %v", err)
		}
	}
	return posts, nil
}

func FetchOlderPostsFromDB(userID string, offset, limit int) ([]models.Post, error) {
	query := `
		SELECT id, text, user_id, created_at
		FROM posts
		WHERE user_id IN (
			SELECT friend_id FROM friends WHERE user_id = $1
		)
		AND visibility = 'public'
		AND is_deleted = false
		ORDER BY created_at DESC
		OFFSET $2 LIMIT $3
	`

	rows, err := db.ExecuteReadQuery(query, userID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Text, &post.UserID, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
