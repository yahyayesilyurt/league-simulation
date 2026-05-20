package config

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestInitLogger_DevelopmentEnv(t *testing.T) {
	t.Setenv("APP_ENV", "development")

	InitLogger()

	assert.Equal(t, zerolog.DebugLevel, zerolog.GlobalLevel())
}

func TestInitLogger_ProductionEnv(t *testing.T) {
	t.Setenv("APP_ENV", "production")

	InitLogger()

	assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())
}

func TestInitLogger_EmptyEnv(t *testing.T) {
	os.Unsetenv("APP_ENV")

	InitLogger()

	assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())
}