package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"otus-highload/models"

	"github.com/redis/go-redis/v9"
)

func EnqueueTask(streamName string, userID string, taskType string, data interface{}) error {
	var eventData map[string]interface{}

	switch v := data.(type) {
	case models.Post:
		postData, err := json.Marshal(v)
		if err != nil {
			log.Printf("Failed to serialize post: %v", err)
			return err
		}
		eventData = map[string]interface{}{
			"taskType": taskType,
			"userID":   userID,
			"data":     string(postData),
		}
	case string:
		// For postID in delete_post
		eventData = map[string]interface{}{
			"taskType": taskType,
			"userID":   userID,
			"data":     v,
		}
	default:
		// For create_feed
		eventData = map[string]interface{}{
			"taskType": taskType,
			"userID":   userID,
		}
	}

	// Add event to Redis stream
	msgID, err := RDB.XAdd(CTX(), &redis.XAddArgs{
		Stream: FeedStreamName,
		Values: eventData,
	}).Result()

	if err != nil {
		log.Printf("Failed to add event to stream: %v", err)
		return err
	}

	fmt.Printf("Event %s added to stream: %s with ID: %s\n", taskType, streamName, msgID)
	return nil
}
