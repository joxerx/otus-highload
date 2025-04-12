package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"otus-highload-counter/models"

	"github.com/redis/go-redis/v9"
)

func StartEventConsumer() {
	ctx := CTX()
	go func() {
		for {
			streams, err := ExternalRDB.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: consumerName,
				Streams:  []string{EventsStreamName, ">"},
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

						_, err := ExternalRDB.XAck(ctx, EventsStreamName, groupName, msg.ID).Result()
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
	err := ExternalRDB.XGroupCreateMkStream(CTX(), EventsStreamName, groupName, "$").Err()
	if err != nil {
		if err.Error() == "BUSYGROUP Consumer Group name already exists" {
			log.Println("Consumer group already exists. Skipping creation.")
			return nil
		}

		// Handle other unexpected errors
		log.Printf("Unexpected error creating group: %v", err)
		return err
	}

	log.Printf("Event group ready %s:%s", EventsStreamName, groupName)
	err = ExternalRDB.XGroupCreateMkStream(CTX(), ErrEventsStreamName, groupName, "$").Err()
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
	eventJSON, ok := values["data"].(string)
	if !ok {
		return fmt.Errorf("event data missing or invalid")
	}

	var event models.CounterEvent
	if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	key := fmt.Sprintf("unreadCount:%d:%d", event.DialogID, event.RecipientID)

	switch event.EventType {
	case "increaseCounter":
		err := RDB.Incr(CTX(), key).Err()
		if err != nil {
			log.Printf("Failed to increase counter: %v", err)
			pushErrorEvent("increaseError", event.MessageID)
		}
	case "decreaseCounter":
		err := RDB.Decr(CTX(), key).Err()
		if err != nil {
			log.Printf("Failed to decrease counter: %v", err)
			pushErrorEvent("decreaseError", event.MessageID)
		}
	default:
		log.Printf("Unknown event type: %s", event.EventType)
	}

	return nil
}

func pushErrorEvent(eventType string, messageID int) {
	data := map[string]interface{}{
		"eventType": eventType,
		"messageId": strconv.Itoa(messageID),
	}

	if err := ExternalRDB.XAdd(CTX(), &redis.XAddArgs{
		Stream: ErrEventsStreamName,
		Values: data,
	}).Err(); err != nil {
		log.Printf("Failed to push %s event: %v", eventType, err)
	}
}
