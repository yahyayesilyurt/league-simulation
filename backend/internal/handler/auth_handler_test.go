package handler

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yahyayesilyurt/league-simulation/internal/service"
)

// Mock Auth Service
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func setupAuthRouter(authSvc service.AuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAuthHandler(authSvc)
	r.POST("/auth/login", h.Login)
	return r
}

func TestLogin_Success(t *testing.T) {
	mockAuth := new(MockAuthService)
	mockAuth.On("Login", "admin", "admin123").
		Return("mock.jwt.token", nil)

	r := setupAuthRouter(mockAuth)
	w := makeRequest(t, r, http.MethodPost, "/auth/login", map[string]string{
		"username": "admin",
		"password": "admin123",
	})

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseResponse(t, w)
	assert.Equal(t, "mock.jwt.token", body["token"])
	assert.Equal(t, "Login successful", body["message"])
	mockAuth.AssertExpectations(t)
}

func TestLogin_WrongCredentials(t *testing.T) {
	mockAuth := new(MockAuthService)
	mockAuth.On("Login", "admin", "yanlis").
		Return("", assert.AnError)

	r := setupAuthRouter(mockAuth)
	w := makeRequest(t, r, http.MethodPost, "/auth/login", map[string]string{
		"username": "admin",
		"password": "yanlis",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	body := parseResponse(t, w)
	assert.Contains(t, body, "error")
}

func TestLogin_MissingUsername(t *testing.T) {
	mockAuth := new(MockAuthService)

	r := setupAuthRouter(mockAuth)
	w := makeRequest(t, r, http.MethodPost, "/auth/login", map[string]string{
		"password": "admin123",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_MissingPassword(t *testing.T) {
	mockAuth := new(MockAuthService)

	r := setupAuthRouter(mockAuth)
	w := makeRequest(t, r, http.MethodPost, "/auth/login", map[string]string{
		"username": "admin",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_EmptyBody(t *testing.T) {
	mockAuth := new(MockAuthService)

	r := setupAuthRouter(mockAuth)
	w := makeRequest(t, r, http.MethodPost, "/auth/login", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}