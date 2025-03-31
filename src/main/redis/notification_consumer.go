package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"otus-highload/models"
	"otus-highload/websocket"

	"github.com/redis/go-redis/v9"
)

func StartNotificationConsumer() {
	ctx := CTX()
	go func() {
		for {
			streams, err := RDB.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: consumerName,
				Streams:  []string{NotificationsStreamName, ">"},
				Count:    10,
				Block:    5 * time.Second,
			}).Result()

			if err != nil && err != redis.Nil {
				log.Printf("Failed to read from notification stream: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			for _, stream := range streams {
				for _, message := range stream.Messages {
					go func(msg redis.XMessage) {
						if err := handleNotification(msg.Values); err != nil {
							log.Printf("Failed to process notification: %v", err)
							return
						}

						_, err := RDB.XAck(ctx, NotificationsStreamName, groupName, msg.ID).Result()
						if err != nil {
							log.Printf("Failed to acknowledge message %s: %v", msg.ID, err)
						} else {
							log.Printf("Notification %s acknowledged", msg.ID)
						}
					}(message)
				}
			}
		}
	}()
}

func handleNotification(values map[string]interface{}) error {
	taskType, ok := values["taskType"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid taskType")
	}

	userID, ok := values["userID"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid userID")
	}

	if rawData, exists := values["data"]; exists {
		dataStr, ok := rawData.(string)
		if !ok {
			return fmt.Errorf("invalid data format")
		}

		if taskType == "notify_friend" {
			var post models.Post
			if err := json.Unmarshal([]byte(dataStr), &post); err != nil {
				return fmt.Errorf("failed to deserialize post: %v", err)
			}
			notification := map[string]string{
				"postId":         post.ID,
				"postText":       post.Text,
				"author_user_id": post.UserID,
			}

			websocket.SendToUser(userID, notification)
			log.Printf("Sent notification to %s for post %s", userID, post.ID)

		}
	}

	return nil
}

func CreateNotificationGroup() error {
	err := RDB.XGroupCreateMkStream(CTX(), NotificationsStreamName, groupName, "$").Err()
	if err != nil {
		if err.Error() == "BUSYGROUP Consumer Group name already exists" {
			log.Println("Consumer group already exists. Skipping creation.")
			return nil
		}

		// Handle other unexpected errors
		log.Printf("Unexpected error creating group: %v", err)
		return err
	}

	log.Printf("Notification group ready %s:%s", NotificationsStreamName, groupName)
	return nil
}
