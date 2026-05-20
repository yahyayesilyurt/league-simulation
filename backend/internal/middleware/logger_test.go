package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func setupLoggerRouter(buf *bytes.Buffer) *gin.Engine {
	log.Logger = zerolog.New(buf).With().Timestamp().Logger()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestLogger())

	r.GET("/success", func(c *gin.Context) {
		c.Status(http.StatusOK) // 200
	})
	r.GET("/client-error", func(c *gin.Context) {
		c.Status(http.StatusBadRequest) // 400
	})
	r.GET("/server-error", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError) // 500
	})

	return r
}

func TestRequestLogger_InfoLevel(t *testing.T) {
	var buf bytes.Buffer
	r := setupLoggerRouter(&buf)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/success", nil)
	r.ServeHTTP(w, req)

	var logOutput map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logOutput)

	assert.NoError(t, err)
	assert.Equal(t, "info", logOutput["level"])
	assert.Equal(t, float64(200), logOutput["status"])
	assert.Equal(t, "/success", logOutput["path"])
	assert.Equal(t, "GET", logOutput["method"])
	assert.Equal(t, "HTTP request", logOutput["message"])
}

func TestRequestLogger_WarnLevel(t *testing.T) {
	var buf bytes.Buffer
	r := setupLoggerRouter(&buf)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/client-error", nil)
	r.ServeHTTP(w, req)

	var logOutput map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logOutput)

	assert.NoError(t, err)
	assert.Equal(t, "warn", logOutput["level"])
	assert.Equal(t, float64(400), logOutput["status"])
	assert.Equal(t, "/client-error", logOutput["path"])
}

func TestRequestLogger_ErrorLevel(t *testing.T) {
	var buf bytes.Buffer
	r := setupLoggerRouter(&buf)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/server-error", nil)
	r.ServeHTTP(w, req)

	var logOutput map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logOutput)

	assert.NoError(t, err)
	assert.Equal(t, "error", logOutput["level"])
	assert.Equal(t, float64(500), logOutput["status"])
	assert.Equal(t, "/server-error", logOutput["path"])
}