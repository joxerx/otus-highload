package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var ExternalRDB *redis.Client

var (
	ErrEventsStreamName = os.Getenv("ERR_EVENTS_STREAM")
	EventsStreamName    = os.Getenv("EVENTS_STREAM")
	groupName           = os.Getenv("MAIN_CONSUMER_GROUP")
	consumerName        = os.Getenv("MAIN_CONSUMER")
)

func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" {
		redisHost = "counter-redis"
	}
	if redisPort == "" {
		redisPort = "6379"
	}

	externalRedisHost := os.Getenv("REDIS_HOST")
	externalRedisPort := os.Getenv("REDIS_PORT")
	if externalRedisHost == "" {
		externalRedisPort = "redis"
	}
	if externalRedisPort == "" {
		externalRedisPort = "6379"
	}

	RDB = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})
	ExternalRDB = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", externalRedisHost, externalRedisPort),
	})
	_, err := RDB.Ping(CTX()).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	_, err = ExternalRDB.Ping(CTX()).Result()
	if err != nil {
		log.Fatalf("failed to connect to ExternalRedis: %v", err)
	}
	log.Println("Redis connected!")
}

func CTX() context.Context {
	return context.Background()
}
