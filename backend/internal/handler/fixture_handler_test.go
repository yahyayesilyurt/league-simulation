package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFixtureService struct {
	mock.Mock
}

func (m *MockFixtureService) GenerateFixture() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockFixtureService) IsFixtureGenerated() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func setupFixtureRouter(mockSvc *MockFixtureService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewFixtureHandler(mockSvc)
	
	r.POST("/league/generate-fixture", handler.GenerateFixture)
	r.GET("/league/fixture-status", handler.GetFixtureStatus)
	
	return r
}

func parseFixtureResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	return response
}

func TestFixtureHandler_GenerateFixture_Success(t *testing.T) {
	mockSvc := new(MockFixtureService)
	mockSvc.On("GenerateFixture").Return(nil)

	r := setupFixtureRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/league/generate-fixture", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	body := parseFixtureResponse(t, w)
	assert.Contains(t, body["message"], "successfully")
	mockSvc.AssertExpectations(t)
}

func TestFixtureHandler_GenerateFixture_Error(t *testing.T) {
	mockSvc := new(MockFixtureService)
	mockSvc.On("GenerateFixture").Return(errors.New("fixture already generated"))

	r := setupFixtureRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/league/generate-fixture", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	body := parseFixtureResponse(t, w)
	assert.Equal(t, "fixture already generated", body["error"])
	mockSvc.AssertExpectations(t)
}

func TestFixtureHandler_GetFixtureStatus_Generated(t *testing.T) {
	mockSvc := new(MockFixtureService)
	mockSvc.On("IsFixtureGenerated").Return(true, nil)

	r := setupFixtureRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/league/fixture-status", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseFixtureResponse(t, w)
	assert.Equal(t, true, body["fixture_generated"])
	mockSvc.AssertExpectations(t)
}

func TestFixtureHandler_GetFixtureStatus_NotGenerated(t *testing.T) {
	mockSvc := new(MockFixtureService)
	mockSvc.On("IsFixtureGenerated").Return(false, nil)

	r := setupFixtureRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/league/fixture-status", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseFixtureResponse(t, w)
	assert.Equal(t, false, body["fixture_generated"])
	mockSvc.AssertExpectations(t)
}

func TestFixtureHandler_GetFixtureStatus_Error(t *testing.T) {
	mockSvc := new(MockFixtureService)
	mockSvc.On("IsFixtureGenerated").Return(false, errors.New("db connection failed"))

	r := setupFixtureRouter(mockSvc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/league/fixture-status", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	body := parseFixtureResponse(t, w)
	assert.Equal(t, "db connection failed", body["error"])
	mockSvc.AssertExpectations(t)
}