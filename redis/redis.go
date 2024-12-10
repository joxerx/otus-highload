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
	_, err := RDB.Ping(ctx()).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	log.Println("Redis connected!")
}

func IncrementSlaveCounter(slaveID string) error {
	ctx := ctx()
	key := fmt.Sprintf("slave:%s:requests", slaveID)
	return RDB.Incr(ctx, key).Err()
}

func GetLeastLoadedSlave(slaves []string) (string, error) {
	ctx := ctx()
	minRequests := int64(^uint64(0) >> 1)
	var leastLoaded string

	for _, slaveID := range slaves {
		key := fmt.Sprintf("slave:%s:requests", slaveID)
		count, err := RDB.Get(ctx, key).Int64()
		if err != nil && err.Error() != "redis: nil" {
			return "", err
		}

		if count < minRequests {
			minRequests = count
			leastLoaded = slaveID
		}
	}

	return leastLoaded, nil
}

func ctx() context.Context {
	return context.Background()
}
