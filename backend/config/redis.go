package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var RedisClient *redis.Client

func ConnectRedis() {
	redisURL := getRedisURL()

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse Redis URL")
	}

	client := redis.NewClient(opts)

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}

	log.Info().Msg("Redis connection established")
	RedisClient = client
}

func getRedisURL() string {
	if url := os.Getenv("UPSTASH_REDIS_URL"); url != "" {
		return url
	}
	return os.Getenv("REDIS_URL")
}