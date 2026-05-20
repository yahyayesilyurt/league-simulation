package config

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	token, err := GenerateToken("admin", "admin")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateToken_NoSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	token, err := GenerateToken("admin", "admin")

	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestValidateToken_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	token, _ := GenerateToken("admin", "admin")
	claims, err := ValidateToken(token)

	assert.NoError(t, err)
	assert.Equal(t, "admin", claims.Username)
	assert.Equal(t, "admin", claims.Role)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	claims, err := ValidateToken("invalid.token.here")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret-one")
	token, _ := GenerateToken("admin", "admin")

	os.Setenv("JWT_SECRET", "secret-two")
	defer os.Unsetenv("JWT_SECRET")

	claims, err := ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_EmptyToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	claims, err := ValidateToken("")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_UnexpectedSigningMethod(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	claims := &Claims{Username: "hacker", Role: "admin"}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	parsedClaims, err := ValidateToken(tokenString)

	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
	assert.Contains(t, err.Error(), "unexpected signing method")
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	secret := "test-secret"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	claims := &Claims{
		Username: "admin",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), 
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "league-simulation",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))

	parsedClaims, err := ValidateToken(tokenString)

	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
	assert.Contains(t, err.Error(), "token has invalid claims: token is expired")
}