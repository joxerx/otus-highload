package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
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
