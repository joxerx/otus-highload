package redis

import (
	"context"
	"encoding/json"
	"log"
	"otus-highload-messages/models"
	"time"

	"github.com/redis/go-redis/v9"
)

func StartEventConsumer() {
	ctx := CTX()
	go func() {
		for {
			streams, err := RDB.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: consumerName,
				Streams:  []string{ErrEventsStreamName, ">"},
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
						if err := handleEvent(msg.Values); err != nil {
							log.Printf("Failed to process notification: %v", err)
							return
						}

						_, err := RDB.XAck(ctx, ErrEventsStreamName, groupName, msg.ID).Result()
						if err != nil {
							log.Printf("Failed to acknowledge message %s: %v", msg.ID, err)
						} else {
							log.Printf("Event %s acknowledged", msg.ID)
						}
					}(message)
				}
			}
		}
	}()
}

func CreateEventGroup() error {
	err := RDB.XGroupCreateMkStream(CTX(), ErrEventsStreamName, groupName, "$").Err()
	if err != nil {
		if err.Error() == "BUSYGROUP Consumer Group name already exists" {
			log.Println("Consumer group already exists. Skipping creation.")
			return nil
		}

		// Handle other unexpected errors
		log.Printf("Unexpected error creating group: %v", err)
		return err
	}

	log.Printf("Event group ready %s:%s", ErrEventsStreamName, groupName)
	return nil

}
func handleEvent(values map[string]interface{}) error {
	ctx := context.Background()

	raw, ok := values["data"].(string)
	if !ok {
		return nil
	}

	var event models.ErrorEvent
	if err := json.Unmarshal([]byte(raw), &event); err != nil {
		log.Printf("Failed to parse error event: %v", err)
		return err
	}

	switch event.EventType {
	case "increaseError":
		log.Printf("increaseError received for message ID: %d. Deleting message.", event.MessageID)
		if err := DeleteMessageByID(ctx, event.MessageID); err != nil {
			log.Printf("Failed to delete message %d: %v", event.MessageID, err)
			return err
		}
	case "decreaseError":
		log.Printf("decreaseError received for message ID: %d. Marking message as unread.", event.MessageID)
		if err := MarkMessageAsUnread(ctx, event.MessageID); err != nil {
			log.Printf("Failed to mark message %d as unread: %v", event.MessageID, err)
			return err
		}
	default:
		log.Printf("Unknown event type: %s", event.EventType)
	}

	return nil
}

func PushEvent(eventType string, chatID string, userID string, messageID int) {
	data := map[string]interface{}{
		"eventType":   eventType,
		"dialogID":    chatID,
		"recipientID": userID,
		"messageID":   messageID,
	}
	if err := RDB.XAdd(CTX(), &redis.XAddArgs{
		Stream: EventsStreamName,
		Values: data,
	}).Err(); err != nil {
		log.Printf("Failed to push %s event: %v", eventType, err)
	}
}
