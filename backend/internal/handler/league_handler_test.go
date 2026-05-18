package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"github.com/yahyayesilyurt/league-simulation/internal/service"
)

// GET /league/table

func TestGetStandings_Success(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	standings := []model.Standing{
		{TeamID: 1, Points: 9, Team: model.Team{Name: "Man City"}},
		{TeamID: 2, Points: 6, Team: model.Team{Name: "Liverpool"}},
	}

	mockLeague.On("GetStandings").Return(standings, nil)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodGet, "/league/table", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseResponse(t, w)
	assert.Contains(t, body, "standings")
	mockLeague.AssertExpectations(t)
}

func TestGetStandings_EmptyTable(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	mockLeague.On("GetStandings").Return([]model.Standing{}, nil)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodGet, "/league/table", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	mockLeague.AssertExpectations(t)
}

// GET /league/week/:weekNo

func TestGetWeek_ValidWeek(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	homeGoals := 2
	awayGoals := 1
	matches := []model.Match{
		{
			ID:        1,
			Week:      1,
			HomeGoals: &homeGoals,
			AwayGoals: &awayGoals,
			Played:    true,
		},
	}

	mockLeague.On("GetWeek", 1).Return(matches, nil)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodGet, "/league/week/1", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseResponse(t, w)
	assert.Equal(t, float64(1), body["week"])
	assert.Contains(t, body, "matches")
	mockLeague.AssertExpectations(t)
}

func TestGetWeek_InvalidWeek_Zero(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodGet, "/league/week/0", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	body := parseResponse(t, w)
	assert.Contains(t, body, "error")
}

func TestGetWeek_InvalidWeek_Seven(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodGet, "/league/week/7", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetWeek_InvalidWeek_String(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodGet, "/league/week/abc", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// GET /league/status

func TestGetStatus_NotStarted(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	status := &model.LeagueStatus{
		CurrentWeek:    0,
		TotalWeeks:     6,
		LeagueFinished: false,
		MatchesPlayed:  0,
		MatchesLeft:    12,
		Status:         "not_started",
	}

	mockLeague.On("GetStatus").Return(status, nil)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodGet, "/league/status", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseResponse(t, w)
	assert.Equal(t, "not_started", body["status"])
	assert.Equal(t, float64(0), body["current_week"])
	assert.Equal(t, float64(12), body["matches_left"])
	mockLeague.AssertExpectations(t)
}

func TestGetStatus_Finished(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	status := &model.LeagueStatus{
		CurrentWeek:    6,
		LeagueFinished: true,
		MatchesPlayed:  12,
		MatchesLeft:    0,
		Status:         "finished",
	}

	mockLeague.On("GetStatus").Return(status, nil)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodGet, "/league/status", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseResponse(t, w)
	assert.Equal(t, "finished", body["status"])
	assert.Equal(t, true, body["league_finished"])
}

// POST /league/next-week

func TestNextWeek_Success(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	homeGoals := 2
	awayGoals := 0
	result := &model.WeekResult{
		Week: 1,
		Matches: []model.Match{
			{HomeGoals: &homeGoals, AwayGoals: &awayGoals, Played: true},
		},
		Standings:      []model.Standing{},
		LeagueFinished: false,
	}

	mockLeague.On("NextWeek").Return(result, nil)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodPost, "/league/next-week", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseResponse(t, w)
	assert.Equal(t, float64(1), body["week"])
	assert.Equal(t, false, body["league_finished"])
	mockLeague.AssertExpectations(t)
}

func TestNextWeek_LeagueFinished(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	mockLeague.On("NextWeek").Return(
		(*model.WeekResult)(nil),
		assert.AnError,
	)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodPost, "/league/next-week", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	body := parseResponse(t, w)
	assert.Contains(t, body, "error")
}

// GET /league/predictions

func TestGetPredictions_Success(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	predictions := []service.TeamPrediction{
		{TeamName: "Man City",  Percentage: 45.0},
		{TeamName: "Liverpool", Percentage: 30.0},
		{TeamName: "Arsenal",   Percentage: 15.0},
		{TeamName: "Chelsea",   Percentage: 10.0},
	}

	mockPred.On("GetPredictions").Return(predictions, nil)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodGet, "/league/predictions", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseResponse(t, w)
	assert.Contains(t, body, "predictions")
	mockPred.AssertExpectations(t)
}

// POST /league/reset

func TestReset_Success(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	status := &model.LeagueStatus{
		CurrentWeek: 0,
		Status:      "not_started",
	}

	mockLeague.On("Reset").Return(status, nil)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodPost, "/league/reset", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseResponse(t, w)
	assert.Contains(t, body, "message")
	assert.Contains(t, body, "status")
	mockLeague.AssertExpectations(t)
}