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

// GET /league/fixtures

func TestGetFixtures_Success(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	fixtures := []model.Match{
		{ID: 1, Week: 1, HomeTeamID: 1, AwayTeamID: 2},
		{ID: 2, Week: 1, HomeTeamID: 3, AwayTeamID: 4},
	}

	mockLeague.On("GetFixtures").Return(fixtures, nil)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodGet, "/league/fixtures", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseResponse(t, w)
	assert.Contains(t, body, "fixtures")
	mockLeague.AssertExpectations(t)
}

// POST /league/play-all

func TestPlayAll_Success(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	result := &model.PlayAllResult{
		TotalWeeksPlayed: 6,
		Weeks:            []model.WeekResult{},
	}

	mockLeague.On("PlayAll").Return(result, nil)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodPost, "/league/play-all", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	mockLeague.AssertExpectations(t)
}

func TestPlayAll_AlreadyFinished(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	mockLeague.On("PlayAll").Return(
		(*model.PlayAllResult)(nil),
		assert.AnError,
	)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodPost, "/league/play-all", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockLeague.AssertExpectations(t)
}

// GET /league/predictions

func TestGetPredictions_Empty(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	mockPred.On("GetPredictions").Return([]service.TeamPrediction{}, nil)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodGet, "/league/predictions", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPred.AssertExpectations(t)
}

// Week number edge cases

func TestGetWeek_Week6_Valid(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	mockLeague.On("GetWeek", 6).Return([]model.Match{}, nil)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodGet, "/league/week/6", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	mockLeague.AssertExpectations(t)
}

func TestGetWeek_Week1_Valid(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	mockLeague.On("GetWeek", 1).Return([]model.Match{}, nil)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodGet, "/league/week/1", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	mockLeague.AssertExpectations(t)
}

// GET /league/status — in_progress

func TestGetStatus_InProgress(t *testing.T) {
	mockLeague := new(MockLeagueService)
	mockPred   := new(MockPredictionService)

	status := &model.LeagueStatus{
		CurrentWeek:    3,
		LeagueFinished: false,
		MatchesPlayed:  6,
		MatchesLeft:    6,
		Status:         "in_progress",
	}

	mockLeague.On("GetStatus").Return(status, nil)

	r := setupTestRouter(mockLeague, mockPred)
	w := makeRequest(t, r, http.MethodGet, "/league/status", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseResponse(t, w)
	assert.Equal(t, "in_progress", body["status"])
	assert.Equal(t, float64(3), body["current_week"])
	mockLeague.AssertExpectations(t)
}

// GET /league/table - Service Error

func TestGetStandings_ServiceError(t *testing.T) {
    mockLeague := new(MockLeagueService)
    mockPred   := new(MockPredictionService)

    mockLeague.On("GetStandings").Return(([]model.Standing)(nil), assert.AnError)

    r := setupTestRouter(mockLeague, mockPred)
    w := makeRequest(t, r, http.MethodGet, "/league/table", nil)

    assert.Equal(t, http.StatusInternalServerError, w.Code)

    body := parseResponse(t, w)
    assert.Contains(t, body, "error")
    mockLeague.AssertExpectations(t)
}

// GET /league/status - Service Error

func TestGetStatus_ServiceError(t *testing.T) {
    mockLeague := new(MockLeagueService)
    mockPred   := new(MockPredictionService)

    mockLeague.On("GetStatus").Return((*model.LeagueStatus)(nil), assert.AnError)

    r := setupTestRouter(mockLeague, mockPred)
    w := makeRequest(t, r, http.MethodGet, "/league/status", nil)

    assert.Equal(t, http.StatusInternalServerError, w.Code)
    
    body := parseResponse(t, w)
    assert.Contains(t, body, "error")
    mockLeague.AssertExpectations(t)
}

// GET /league/predictions - Service Error

func TestGetPredictions_ServiceError(t *testing.T) {
    mockLeague := new(MockLeagueService)
    mockPred   := new(MockPredictionService)

    mockPred.On("GetPredictions").Return(([]service.TeamPrediction)(nil), assert.AnError)

    r := setupTestRouter(mockLeague, mockPred)
    w := makeRequest(t, r, http.MethodGet, "/league/predictions", nil)

    assert.Equal(t, http.StatusInternalServerError, w.Code)
    mockPred.AssertExpectations(t)
}

// POST /league/reset - Service Error

func TestReset_ServiceError(t *testing.T) {
    mockLeague := new(MockLeagueService)
    mockPred   := new(MockPredictionService)

    mockLeague.On("Reset").Return((*model.LeagueStatus)(nil), assert.AnError)

    r := setupTestRouter(mockLeague, mockPred)
    w := makeRequest(t, r, http.MethodPost, "/league/reset", nil)

    assert.Equal(t, http.StatusInternalServerError, w.Code)
    
    body := parseResponse(t, w)
    assert.Contains(t, body, "error")
    mockLeague.AssertExpectations(t)
}