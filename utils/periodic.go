package utils

import (
	"log"
	"otus-highload/db"
	"otus-highload/redis"
	"time"
)

const batchSize = 1000
const periodicHours = 1

func enqueueTasksInBatches(userIDs []string) {
	for i := 0; i < len(userIDs); i += batchSize {
		end := i + batchSize
		if end > len(userIDs) {
			end = len(userIDs)
		}

		batch := userIDs[i:end]

		for _, userID := range batch {
			if err := redis.EnqueueTask(userID, "create_feed", nil); err != nil {
				log.Printf("Failed to enqueue task for user %s: %v", userID, err)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func enqueueTasksForAllUsers() {
	rows, err := db.MasterDB.Query("SELECT id FROM users")
	if err != nil {
		log.Printf("Error fetching user IDs: %v", err)
		return
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			log.Printf("Error scanning user ID: %v", err)
			continue
		}
		userIDs = append(userIDs, userID)
	}

	enqueueTasksInBatches(userIDs)
}

func StartPeriodicTask() {
	enqueueTasksForAllUsers()

	ticker := time.NewTicker(periodicHours * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		enqueueTasksForAllUsers()
	}
}
