package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"github.com/yahyayesilyurt/league-simulation/internal/service"
)

// Mock Services

type MockLeagueService struct {
	mock.Mock
}

func (m *MockLeagueService) GetStandings() ([]model.Standing, error) {
	args := m.Called()
	return args.Get(0).([]model.Standing), args.Error(1)
}

func (m *MockLeagueService) GetFixtures() ([]model.Match, error) {
	args := m.Called()
	return args.Get(0).([]model.Match), args.Error(1)
}

func (m *MockLeagueService) GetWeek(week int) ([]model.Match, error) {
	args := m.Called(week)
	return args.Get(0).([]model.Match), args.Error(1)
}

func (m *MockLeagueService) GetCurrentWeek() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockLeagueService) GetStatus() (*model.LeagueStatus, error) {
	args := m.Called()
	return args.Get(0).(*model.LeagueStatus), args.Error(1)
}

func (m *MockLeagueService) NextWeek() (*model.WeekResult, error) {
	args := m.Called()
	return args.Get(0).(*model.WeekResult), args.Error(1)
}

func (m *MockLeagueService) PlayAll() (*model.PlayAllResult, error) {
	args := m.Called()
	return args.Get(0).(*model.PlayAllResult), args.Error(1)
}

func (m *MockLeagueService) Reset() (*model.LeagueStatus, error) {
	args := m.Called()
	return args.Get(0).(*model.LeagueStatus), args.Error(1)
}

type MockPredictionService struct {
	mock.Mock
}

func (m *MockPredictionService) GetPredictions() ([]service.TeamPrediction, error) {
	args := m.Called()
	return args.Get(0).([]service.TeamPrediction), args.Error(1)
}

// Test Helper Functions

func setupTestRouter(leagueSvc *MockLeagueService, predSvc *MockPredictionService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	h := NewLeagueHandler(leagueSvc, predSvc)

	league := r.Group("/league")
	{
		league.GET("/table",        h.GetStandings)
		league.GET("/fixtures",     h.GetFixtures)
		league.GET("/week/:weekNo", h.GetWeek)
		league.GET("/predictions",  h.GetPredictions)
		league.GET("/status",       h.GetStatus)
		league.POST("/next-week",   h.NextWeek)
		league.POST("/play-all",    h.PlayAll)
		league.POST("/reset",       h.Reset)
	}

	return r
}

func makeRequest(t *testing.T, r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	return result
}