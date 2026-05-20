package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupRouterTestDeps(t *testing.T) (*gorm.DB, *redis.Client) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return db, redisClient
}

func TestSetupRouter_HealthCheck(t *testing.T) {
	db, redisClient := setupRouterTestDeps(t)
	r := SetupRouter(db, redisClient)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "League Simulation is running", response["message"])
}

func TestSetupRouter_ProtectedRoutesMiddleware(t *testing.T) {
	db, redisClient := setupRouterTestDeps(t)
	r := SetupRouter(db, redisClient)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/league/reset", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSetupRouter_CORS(t *testing.T) {
	db, redisClient := setupRouterTestDeps(t)
	r := SetupRouter(db, redisClient)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodOptions, "/league/table", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "GET")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
}