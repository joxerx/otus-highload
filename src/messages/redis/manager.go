package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"otus-highload-messages/models"
)

func DeleteMessageByID(ctx context.Context, messageID int) error {
	keys, err := RDB.Keys(ctx, "chat:*").Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		messages, err := RDB.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			continue
		}

		for _, raw := range messages {
			var msg models.Message
			if err := json.Unmarshal([]byte(raw), &msg); err != nil {
				continue
			}

			if msg.ID == messageID {
				// Use LRem to remove the message
				_, err := RDB.LRem(ctx, key, 1, raw).Result()
				if err != nil {
					return err
				}
				log.Printf("Deleted message %d from chat %s", messageID, key)
				return nil
			}
		}
	}

	return fmt.Errorf("message with ID %d not found", messageID)
}

func MarkMessageAsUnread(ctx context.Context, messageID int) error {
	keys, err := RDB.Keys(ctx, "chat:*").Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		messages, err := RDB.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			continue
		}

		for i, raw := range messages {
			var msg models.Message
			if err := json.Unmarshal([]byte(raw), &msg); err != nil {
				continue
			}

			if msg.ID == messageID {
				msg.IsRead = false
				newRaw, err := json.Marshal(msg)
				if err != nil {
					return err
				}
				if err := RDB.LSet(ctx, key, int64(i), newRaw).Err(); err != nil {
					return err
				}
				log.Printf("Marked message %d as unread in chat %s", messageID, key)
				return nil
			}
		}
	}

	return fmt.Errorf("message with ID %d not found", messageID)
}
