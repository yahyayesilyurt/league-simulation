package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yahyayesilyurt/league-simulation/internal/model"
)

type MockMatchService struct {
	mock.Mock
}

func (m *MockMatchService) UpdateMatchResult(id uint, homeGoals, awayGoals int) (*model.Match, error) {
	args := m.Called(id, homeGoals, awayGoals)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Match), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMatchService) GetMatchByID(id uint) (*model.Match, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Match), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMatchService) PlayMatch(id uint) (*model.Match, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Match), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStandingService) RecalculateAll() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMatchService) PlayWeek(week int) ([]model.Match, error) {
	args := m.Called(week)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Match), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockStandingService struct {
	mock.Mock
}

func (m *MockStandingService) GetStandings() ([]model.Standing, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]model.Standing), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupMatchRouter(matchSvc *MockMatchService, standingSvc *MockStandingService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewMatchHandler(matchSvc, standingSvc)

	r.PUT("/match/:id/result", handler.UpdateResult)
	r.GET("/match/:id", handler.GetMatch)

	return r
}

func parseMatchResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	return response
}

func TestMatchHandler_UpdateResult_Success(t *testing.T) {
	mockMatch := new(MockMatchService)
	mockStanding := new(MockStandingService)

	updatedMatch := &model.Match{ID: 1}
	standings := []model.Standing{{TeamID: 1, Points: 3}}

	mockMatch.On("UpdateMatchResult", uint(1), 2, 1).Return(updatedMatch, nil)
	mockStanding.On("GetStandings").Return(standings, nil)

	r := setupMatchRouter(mockMatch, mockStanding)
	w := httptest.NewRecorder()

	body := []byte(`{"home_goals": 2, "away_goals": 1}`)
	req, _ := http.NewRequest(http.MethodPut, "/match/1/result", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	respBody := parseMatchResponse(t, w)
	assert.Equal(t, "Match result updated successfully", respBody["message"])
	assert.Contains(t, respBody, "match")
	assert.Contains(t, respBody, "standings")

	mockMatch.AssertExpectations(t)
	mockStanding.AssertExpectations(t)
}

func TestMatchHandler_UpdateResult_InvalidID(t *testing.T) {
	r := setupMatchRouter(new(MockMatchService), new(MockStandingService))
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/match/abc/result", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPut, "/match/0/result", nil)
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestMatchHandler_UpdateResult_InvalidJSON(t *testing.T) {
	r := setupMatchRouter(new(MockMatchService), new(MockStandingService))
	w := httptest.NewRecorder()

	body := []byte(`{"home_goals": -1, "away_goals": 2}`)
	req, _ := http.NewRequest(http.MethodPut, "/match/1/result", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	respBody := parseMatchResponse(t, w)
	assert.Contains(t, respBody["error"], "Invalid request body")
}

func TestMatchHandler_UpdateResult_MatchServiceError(t *testing.T) {
	mockMatch := new(MockMatchService)
	mockMatch.On("UpdateMatchResult", uint(1), 2, 1).Return(nil, errors.New("match already played"))

	r := setupMatchRouter(mockMatch, new(MockStandingService))
	w := httptest.NewRecorder()

	body := []byte(`{"home_goals": 2, "away_goals": 1}`)
	req, _ := http.NewRequest(http.MethodPut, "/match/1/result", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	respBody := parseMatchResponse(t, w)
	assert.Equal(t, "match already played", respBody["error"])
}

func TestMatchHandler_UpdateResult_StandingServiceError(t *testing.T) {
	mockMatch := new(MockMatchService)
	mockStanding := new(MockStandingService)

	mockMatch.On("UpdateMatchResult", uint(1), 2, 1).Return(&model.Match{ID: 1}, nil)
	mockStanding.On("GetStandings").Return(nil, errors.New("db error"))

	r := setupMatchRouter(mockMatch, mockStanding)
	w := httptest.NewRecorder()

	body := []byte(`{"home_goals": 2, "away_goals": 1}`)
	req, _ := http.NewRequest(http.MethodPut, "/match/1/result", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMatchHandler_GetMatch_Success(t *testing.T) {
	mockMatch := new(MockMatchService)
	expectedMatch := &model.Match{ID: 5}

	mockMatch.On("GetMatchByID", uint(5)).Return(expectedMatch, nil)

	r := setupMatchRouter(mockMatch, new(MockStandingService))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/match/5", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	respBody := parseMatchResponse(t, w)
	assert.Contains(t, respBody, "match")
}

func TestMatchHandler_GetMatch_InvalidID(t *testing.T) {
	r := setupMatchRouter(new(MockMatchService), new(MockStandingService))
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/match/invalid", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMatchHandler_GetMatch_NotFound(t *testing.T) {
	mockMatch := new(MockMatchService)
	mockMatch.On("GetMatchByID", uint(99)).Return(nil, errors.New("not found"))

	r := setupMatchRouter(mockMatch, new(MockStandingService))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/match/99", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}