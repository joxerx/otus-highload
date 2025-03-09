package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"otus-highload/models"
	"sync"
	"time"
)

const maxConcurrentTasks = 100

func StartTaskConsumer(stopChan <-chan struct{}) {
	key := "task_queue"

	semaphore := make(chan struct{}, maxConcurrentTasks)
	var wg sync.WaitGroup

	for {
		select {
		case <-stopChan:
			log.Println("Stopping task consumer...")
			wg.Wait()
			close(semaphore)
			log.Println("Task consumer stopped.")
			return
		default:
			result, err := RDB.BRPop(CTX(), 0, key).Result()
			if err != nil {
				log.Printf("Error consuming task: %v", err)
				continue
			}

			if len(result) > 1 {
				taskJSON := result[1]
				log.Printf("Processing task: %s", taskJSON)

				semaphore <- struct{}{}
				wg.Add(1)

				go func(taskJSON string) {
					defer func() {
						<-semaphore
						wg.Done()
					}()

					var task map[string]interface{}
					if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
						log.Printf("Failed to parse task: %v", err)
						return
					}

					userID, ok := task["userID"].(string)
					if !ok {
						log.Printf("Invalid userID in task: %s", taskJSON)
						return
					}

					taskType, ok := task["task"].(string)
					if !ok {
						log.Printf("Invalid task type in task: %s", taskJSON)
						return
					}

					data := task["data"]

					log.Printf("Executing task type '%s' for user '%s'", taskType, userID)
					err := HandleTask(taskType, userID, data)
					if taskType == "create_feed" || taskType == "delete_post" {
						setKey := "task_queue_set:" + taskType
						taskID := userID + ":" + taskType
						RDB.SRem(CTX(), setKey, taskID)
					}

					if err != nil {
						log.Printf("Error executing task for user '%s': %v", userID, err)
					}
				}(taskJSON)
			}
		}
	}
}

func EnqueueTask(userID, taskType string, data interface{}) error {
	task := map[string]interface{}{
		"userID": userID,
		"task":   taskType,
		"data":   data,
	}
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	key := "task_queue"

	if taskType == "create_feed" || taskType == "delete_post" {
		setKey := "task_queue_set:" + taskType
		taskID := userID + ":" + taskType

		exists, err := RDB.SIsMember(CTX(), setKey, taskID).Result()
		if err != nil {
			return err
		}

		if exists {
			log.Printf("Task %s for user %s already exists, skipping enqueue", taskType, userID)
			return nil
		}

		if err := RDB.SAdd(CTX(), setKey, taskID).Err(); err != nil {
			return err
		}

		if err := RDB.LPush(CTX(), key, taskJSON).Err(); err != nil {
			return err
		}
		RDB.Expire(CTX(), setKey, 45*time.Hour)

	} else {
		if err := RDB.LPush(CTX(), key, taskJSON).Err(); err != nil {
			return err
		}
	}

	return nil
}

func HandleTask(taskType string, userID string, data interface{}) error {
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
