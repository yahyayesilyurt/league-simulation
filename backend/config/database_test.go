package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDSN_WithSupabaseURL(t *testing.T) {
	expectedURL := "postgresql://postgres:password@db.supabase.co:5432/postgres"
	t.Setenv("SUPABASE_URL", expectedURL)

	dsn := buildDSN()

	assert.Equal(t, expectedURL, dsn)
}

func TestBuildDSN_WithIndividualDBVars(t *testing.T) {
	os.Unsetenv("SUPABASE_URL") 
	
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "admin")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "league_db")

	dsn := buildDSN()

	expectedDSN := "host=localhost port=5432 user=admin password=secret dbname=league_db sslmode=disable TimeZone=UTC"
	assert.Equal(t, expectedDSN, dsn)
}

func TestBuildDSN_EmptyVars(t *testing.T) {
	os.Unsetenv("SUPABASE_URL")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")

	dsn := buildDSN()

	expectedDSN := "host= port= user= password= dbname= sslmode=disable TimeZone=UTC"
	assert.Equal(t, expectedDSN, dsn)
}