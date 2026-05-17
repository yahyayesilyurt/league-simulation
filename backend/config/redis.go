package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() {
	redisURL := getRedisURL()

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opts)

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Printf("Redis connection established")
	RedisClient = client
}

func getRedisURL() string {
	if url := os.Getenv("UPSTASH_REDIS_URL"); url != "" {
		return url
	}
	return os.Getenv("REDIS_URL")
}