package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	StandingsKey   = "league:standings"
	PredictionsKey = "league:predictions"
	StatusKey      = "league:status"
	TTL            = 5 * time.Minute
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Set(ctx context.Context, key string, value string) error {
	return c.client.Set(ctx, key, value, TTL).Err()
}

func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

func (c *Cache) InvalidateLeague(ctx context.Context) error {
	return c.Delete(ctx, StandingsKey, PredictionsKey, StatusKey)
}

func (c *Cache) Exists(ctx context.Context, key string) bool {
	result, err := c.client.Exists(ctx, key).Result()
	return err == nil && result > 0
}