package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	ErrEventsStreamName = os.Getenv("ERR_EVENTS_STREAM")
	EventsStreamName    = os.Getenv("EVENTS_STREAM")
	groupName           = os.Getenv("MAIN_CONSUMER_GROUP")
	consumerName        = os.Getenv("MAIN_CONSUMER")
)

var RDB *redis.Client

func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" {
		redisHost = "redis"
	}
	if redisPort == "" {
		redisPort = "6379"
	}

	RDB = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})
	_, err := RDB.Ping(CTX()).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	log.Println("Redis connected!")
}

func CTX() context.Context {
	return context.Background()
}
