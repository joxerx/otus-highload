package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"otus-highload/models"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	FeedStreamName          = os.Getenv("FEED_STREAM")
	NotificationsStreamName = os.Getenv("NOTIFICATIONS_STREAM")
	groupName               = os.Getenv("MAIN_CONSUMER_GROUP")
	consumerName            = os.Getenv("MAIN_CONSUMER")
)

func StartFeedTasksConsumer() {
	for {
		ctx := CTX()
		streams, err := RDB.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerName,
			Streams:  []string{FeedStreamName, ">"},
			Count:    10,
			Block:    time.Second,
		}).Result()

		if err != nil && err != redis.Nil {
			log.Printf("Failed to read from stream: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				go func(msg redis.XMessage) {
					if err := handleEntry(msg.Values); err != nil {
						log.Printf("Failed to process message: %v", err)
						return
					}

					_, err := RDB.XAck(ctx, FeedStreamName, groupName, msg.ID).Result()
					if err != nil {
						log.Printf("Failed to acknowledge message %s: %v", msg.ID, err)
					} else {
						log.Printf("Message %s acknowledged", msg.ID)
					}
				}(message)
			}
		}
	}
}

func handleEntry(values map[string]interface{}) error {
	taskType, ok := values["taskType"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid taskType")
	}

	userID, ok := values["userID"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid userID")
	}

	// Handle data serialization
	var data interface{}

	if rawData, exists := values["data"]; exists {
		dataStr, ok := rawData.(string)
		if !ok {
			return fmt.Errorf("invalid data format")
		}

		if taskType == "append_post" || taskType == "update_post" {
			var post models.Post
			if err := json.Unmarshal([]byte(dataStr), &post); err != nil {
				return fmt.Errorf("failed to deserialize post: %v", err)
			}
			data = post
		} else if taskType == "delete_post" {
			data = dataStr
		}
	}

	return handleTask(taskType, userID, data)
}

func handleTask(taskType string, userID string, data interface{}) error {
	switch taskType {
	case "create_feed":
		return CreateWholeFeed(userID)

	case "append_post":
		post, ok := data.(models.Post)
		if !ok {
			return fmt.Errorf("invalid data for append_post")
		}
		return AppendPostToFeed(userID, post)

	case "delete_post":
		postID, ok := data.(string)
		if !ok {
			return fmt.Errorf("invalid data for delete_post")
		}
		return DeletePostFromFeed(userID, postID)

	case "update_post":
		updatedPost, ok := data.(models.Post)
		if !ok {
			return fmt.Errorf("invalid data for update_post")
		}
		return UpdatePostInFeed(userID, updatedPost)

	default:
		return fmt.Errorf("unknown task type: %s", taskType)
	}
}

func CreateFeederGroup() error {
	err := RDB.XGroupCreateMkStream(CTX(), FeedStreamName, groupName, "$").Err()
	if err != nil {
		if err.Error() == "BUSYGROUP Consumer Group name already exists" {
			log.Println("Consumer group already exists. Skipping creation.")
			return nil
		}

		// Handle other unexpected errors
		log.Printf("Unexpected error creating group: %v", err)
		return err
	}

	log.Printf("Feed group ready %s:%s", FeedStreamName, groupName)
	return nil
}
