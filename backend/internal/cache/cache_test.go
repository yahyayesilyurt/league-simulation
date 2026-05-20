package cache

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockData struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

func setupCacheTest(t *testing.T) (*miniredis.Miniredis, *Cache) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	t.Cleanup(mr.Close)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return mr, NewCache(client)
}

func TestCache_SetAndGetJSON_Success(t *testing.T) {
	_, cache := setupCacheTest(t)
	ctx := context.Background()

	data := mockData{Name: "Test", Score: 100}

	err := cache.SetJSON(ctx, "test_key", data)
	assert.NoError(t, err)

	var fetchedData mockData
	err = cache.GetJSON(ctx, "test_key", &fetchedData)

	assert.NoError(t, err)
	assert.Equal(t, data.Name, fetchedData.Name)
	assert.Equal(t, data.Score, fetchedData.Score)
}

func TestCache_GetJSON_CacheMiss(t *testing.T) {
	_, cache := setupCacheTest(t)
	ctx := context.Background()

	var fetchedData mockData
	err := cache.GetJSON(ctx, "non_existent_key", &fetchedData)

	assert.Error(t, err)
	assert.Equal(t, ErrCacheMiss, err)
}

func TestCache_GetJSON_UnmarshalError(t *testing.T) {
	mr, cache := setupCacheTest(t)
	ctx := context.Background()

	mr.Set("bad_json_key", "{bad_json: true")

	var fetchedData mockData
	err := cache.GetJSON(ctx, "bad_json_key", &fetchedData)

	assert.Error(t, err)
	assert.NotEqual(t, ErrCacheMiss, err) 
}

func TestCache_SetJSON_MarshalError(t *testing.T) {
	_, cache := setupCacheTest(t)
	ctx := context.Background()

	unsupportedData := make(chan int)

	err := cache.SetJSON(ctx, "error_key", unsupportedData)

	assert.Error(t, err)
}

func TestCache_Delete(t *testing.T) {
	_, cache := setupCacheTest(t)
	ctx := context.Background()

	cache.SetJSON(ctx, "key_to_delete", mockData{Name: "To be deleted"})

	err := cache.Delete(ctx, "key_to_delete")
	assert.NoError(t, err)

	var data mockData
	err = cache.GetJSON(ctx, "key_to_delete", &data)
	assert.Equal(t, ErrCacheMiss, err)
}

func TestCache_InvalidateLeague(t *testing.T) {
	_, cache := setupCacheTest(t)
	ctx := context.Background()

	cache.SetJSON(ctx, StandingsKey, "data")
	cache.SetJSON(ctx, PredictionsKey, "data")
	cache.SetJSON(ctx, StatusKey, "data")
	
	cache.SetJSON(ctx, "other_key", "data")

	err := cache.InvalidateLeague(ctx)
	assert.NoError(t, err)

	var dummy string
	assert.Equal(t, ErrCacheMiss, cache.GetJSON(ctx, StandingsKey, &dummy))
	assert.Equal(t, ErrCacheMiss, cache.GetJSON(ctx, PredictionsKey, &dummy))
	assert.Equal(t, ErrCacheMiss, cache.GetJSON(ctx, StatusKey, &dummy))

	err = cache.GetJSON(ctx, "other_key", &dummy)
	assert.NoError(t, err)
}