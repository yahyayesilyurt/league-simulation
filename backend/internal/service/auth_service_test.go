package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin_Success(t *testing.T) {
	os.Setenv("ADMIN_USERNAME", "admin")
	os.Setenv("ADMIN_PASSWORD", "admin123")
	os.Setenv("JWT_SECRET", "test-secret")
	defer func() {
		os.Unsetenv("ADMIN_USERNAME")
		os.Unsetenv("ADMIN_PASSWORD")
		os.Unsetenv("JWT_SECRET")
	}()

	svc := NewAuthService()
	token, err := svc.Login("admin", "admin123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLogin_WrongPassword(t *testing.T) {
	os.Setenv("ADMIN_USERNAME", "admin")
	os.Setenv("ADMIN_PASSWORD", "admin123")
	defer func() {
		os.Unsetenv("ADMIN_USERNAME")
		os.Unsetenv("ADMIN_PASSWORD")
	}()

	svc := NewAuthService()
	token, err := svc.Login("admin", "wrongpassword")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.EqualError(t, err, "invalid credentials")
}

func TestLogin_WrongUsername(t *testing.T) {
	os.Setenv("ADMIN_USERNAME", "admin")
	os.Setenv("ADMIN_PASSWORD", "admin123")
	defer func() {
		os.Unsetenv("ADMIN_USERNAME")
		os.Unsetenv("ADMIN_PASSWORD")
	}()

	svc := NewAuthService()
	token, err := svc.Login("wronguser", "admin123")

	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestLogin_DefaultCredentials(t *testing.T) {
	os.Unsetenv("ADMIN_USERNAME")
	os.Unsetenv("ADMIN_PASSWORD")
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	svc := NewAuthService()
	token, err := svc.Login("admin", "admin123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLogin_BothWrong(t *testing.T) {
	svc := NewAuthService()
	token, err := svc.Login("wrong", "wrong")

	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestLogin_GenerateTokenError(t *testing.T) {
	t.Setenv("ADMIN_USERNAME", "admin")
	t.Setenv("ADMIN_PASSWORD", "admin123")
	
	os.Unsetenv("JWT_SECRET")

	svc := NewAuthService()
	token, err := svc.Login("admin", "admin123")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.EqualError(t, err, "JWT_SECRET not set")
}