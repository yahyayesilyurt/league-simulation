package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRedisURL_UpstashPriority(t *testing.T) {
	expectedUpstashURL := "redis://default:secret@upstash-server.com:6379"
	t.Setenv("UPSTASH_REDIS_URL", expectedUpstashURL)
	t.Setenv("REDIS_URL", "redis://localhost:6379")

	url := getRedisURL()

	assert.Equal(t, expectedUpstashURL, url)
}

func TestGetRedisURL_FallbackToRedisURL(t *testing.T) {
	os.Unsetenv("UPSTASH_REDIS_URL")
	expectedLocalURL := "redis://localhost:6379"
	t.Setenv("REDIS_URL", expectedLocalURL)

	url := getRedisURL()

	assert.Equal(t, expectedLocalURL, url)
}

func TestGetRedisURL_Empty(t *testing.T) {
	os.Unsetenv("UPSTASH_REDIS_URL")
	os.Unsetenv("REDIS_URL")

	url := getRedisURL()

	assert.Empty(t, url)
}