package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yahyayesilyurt/league-simulation/internal/model"
)

func TestGetPredictions_EqualDistribution_WeekZero(t *testing.T) {
	mockStanding := new(MockStandingRepository)
	mockMatch    := new(MockMatchRepository)
	mockTeam     := new(MockTeamRepository)

	standings := []model.Standing{
		{TeamID: 1, Points: 0, Team: model.Team{Name: "Man City"}},
		{TeamID: 2, Points: 0, Team: model.Team{Name: "Liverpool"}},
		{TeamID: 3, Points: 0, Team: model.Team{Name: "Arsenal"}},
		{TeamID: 4, Points: 0, Team: model.Team{Name: "Chelsea"}},
	}

	mockStanding.On("GetAll").Return(standings, nil)
	mockMatch.On("GetAll").Return([]model.Match{}, nil)

	svc := NewPredictionService(mockStanding, mockMatch, mockTeam)
	result, err := svc.GetPredictions()

	assert.NoError(t, err)
	assert.Len(t, result, 4)

	for _, p := range result {
		assert.Equal(t, 25.0, p.Percentage)
	}
}

func TestGetPredictions_FinishedLeague(t *testing.T) {
	mockStanding := new(MockStandingRepository)
	mockMatch    := new(MockMatchRepository)
	mockTeam     := new(MockTeamRepository)

	standings := []model.Standing{
		{TeamID: 1, Points: 15, GoalDiff: 10, GoalsFor: 20, Team: model.Team{Name: "Man City"}},
		{TeamID: 2, Points: 10, GoalDiff: 3,  GoalsFor: 12, Team: model.Team{Name: "Liverpool"}},
		{TeamID: 3, Points: 6,  GoalDiff: -3, GoalsFor: 8,  Team: model.Team{Name: "Arsenal"}},
		{TeamID: 4, Points: 2,  GoalDiff: -10,GoalsFor: 4,  Team: model.Team{Name: "Chelsea"}},
	}

	playedMatches := make([]model.Match, 12)
	for i := range playedMatches {
		week := (i / 2) + 1
		playedMatches[i] = model.Match{Played: true, Week: week}
	}

	mockStanding.On("GetAll").Return(standings, nil)
	mockMatch.On("GetAll").Return(playedMatches, nil)

	svc := NewPredictionService(mockStanding, mockMatch, mockTeam)
	result, err := svc.GetPredictions()

	assert.NoError(t, err)
	assert.Len(t, result, 4)

	leader := result[0]
	assert.Equal(t, "Man City", leader.TeamName)
	assert.Equal(t, 100.0, leader.Percentage)

	for _, p := range result[1:] {
		assert.Equal(t, 0.0, p.Percentage)
	}
}

func TestGetPredictions_EmptyStandings(t *testing.T) {
	mockStanding := new(MockStandingRepository)
	mockMatch    := new(MockMatchRepository)
	mockTeam     := new(MockTeamRepository)

	mockStanding.On("GetAll").Return([]model.Standing{}, nil)

	svc := NewPredictionService(mockStanding, mockMatch, mockTeam)
	result, err := svc.GetPredictions()

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetPredictions_Week4_RealCalculation(t *testing.T) {
	mockStanding := new(MockStandingRepository)
	mockMatch    := new(MockMatchRepository)
	mockTeam     := new(MockTeamRepository)

	standings := []model.Standing{
		{TeamID: 1, Points: 12, GoalDiff: 8, GoalsFor: 15, Team: model.Team{Name: "Man City"}},
		{TeamID: 2, Points: 7,  GoalDiff: 2, GoalsFor: 10, Team: model.Team{Name: "Liverpool"}},
		{TeamID: 3, Points: 4,  GoalDiff: -3,GoalsFor: 6,  Team: model.Team{Name: "Arsenal"}},
		{TeamID: 4, Points: 1,  GoalDiff: -7,GoalsFor: 3,  Team: model.Team{Name: "Chelsea"}},
	}

	teams := []model.Team{
		{ID: 1, Name: "Man City",  Strength: 90},
		{ID: 2, Name: "Liverpool", Strength: 85},
		{ID: 3, Name: "Arsenal",   Strength: 80},
		{ID: 4, Name: "Chelsea",   Strength: 75},
	}

	playedMatches := []model.Match{}
	for i := 0; i < 8; i++ {
		week := (i / 2) + 1
		playedMatches = append(playedMatches, model.Match{
			Played: true,
			Week:   week,
		})
	}

	unplayedMatches := []model.Match{
		{ID: 9,  Week: 5, HomeTeamID: 1, AwayTeamID: 2, Played: false},
		{ID: 10, Week: 5, HomeTeamID: 3, AwayTeamID: 4, Played: false},
		{ID: 11, Week: 6, HomeTeamID: 1, AwayTeamID: 3, Played: false},
		{ID: 12, Week: 6, HomeTeamID: 2, AwayTeamID: 4, Played: false},
	}

	allMatches := append(playedMatches, unplayedMatches...)

	mockStanding.On("GetAll").Return(standings, nil)
	mockMatch.On("GetAll").Return(allMatches, nil)
	mockMatch.On("GetUnplayed").Return(unplayedMatches, nil)
	mockTeam.On("GetAll").Return(teams, nil)

	svc := NewPredictionService(mockStanding, mockMatch, mockTeam)
	result, err := svc.GetPredictions()

	assert.NoError(t, err)
	assert.Len(t, result, 4)

	total := 0.0
	for _, p := range result {
		total += p.Percentage
	}
	assert.InDelta(t, 100.0, total, 1.0)

	assert.Greater(t, result[0].Percentage, result[len(result)-1].Percentage)
}

func TestEqualDistribution_SumIs100(t *testing.T) {
	mockStanding := new(MockStandingRepository)
	mockMatch    := new(MockMatchRepository)
	mockTeam     := new(MockTeamRepository)

	standings := []model.Standing{
		{TeamID: 1, Team: model.Team{Name: "A"}},
		{TeamID: 2, Team: model.Team{Name: "B"}},
		{TeamID: 3, Team: model.Team{Name: "C"}},
		{TeamID: 4, Team: model.Team{Name: "D"}},
	}

	mockStanding.On("GetAll").Return(standings, nil)
	mockMatch.On("GetAll").Return([]model.Match{}, nil)

	svc := NewPredictionService(mockStanding, mockMatch, mockTeam)
	result, err := svc.GetPredictions()

	assert.NoError(t, err)

	total := 0.0
	for _, p := range result {
		total += p.Percentage
	}
	assert.InDelta(t, 100.0, total, 0.1)
}

func TestGetPredictions_Errors(t *testing.T) {
	mockStanding := new(MockStandingRepository)
	mockStanding.On("GetAll").Return([]model.Standing{}, errors.New("db error")).Once()
	svc1 := NewPredictionService(mockStanding, nil, nil)
	_, err1 := svc1.GetPredictions()
	assert.Error(t, err1)

	mockStanding2 := new(MockStandingRepository)
	mockMatch2 := new(MockMatchRepository)
	mockStanding2.On("GetAll").Return([]model.Standing{{TeamID: 1}}, nil).Once()
	mockMatch2.On("GetAll").Return([]model.Match{}, errors.New("match db error")).Once()
	svc2 := NewPredictionService(mockStanding2, mockMatch2, nil)
	_, err2 := svc2.GetPredictions()
	assert.Error(t, err2)

	mockStanding3 := new(MockStandingRepository)
	mockMatch3 := new(MockMatchRepository)
	mockTeam3 := new(MockTeamRepository)
	
	mockStanding3.On("GetAll").Return([]model.Standing{{TeamID: 1}}, nil).Once()
	mockMatch3.On("GetAll").Return([]model.Match{{Played: true, Week: 4}}, nil).Once()
	mockTeam3.On("GetAll").Return([]model.Team{}, errors.New("team db error")).Once()
	
	svc3 := NewPredictionService(mockStanding3, mockMatch3, mockTeam3)
	_, err3 := svc3.GetPredictions()
	assert.Error(t, err3)

	mockStanding4 := new(MockStandingRepository)
	mockMatch4 := new(MockMatchRepository)
	mockTeam4 := new(MockTeamRepository)

	mockStanding4.On("GetAll").Return([]model.Standing{{TeamID: 1}}, nil).Once()
	mockMatch4.On("GetAll").Return([]model.Match{{Played: true, Week: 4}}, nil).Once()
	mockTeam4.On("GetAll").Return([]model.Team{{ID: 1}}, nil).Once()
	mockMatch4.On("GetUnplayed").Return([]model.Match{}, errors.New("unplayed db error")).Once()

	svc4 := NewPredictionService(mockStanding4, mockMatch4, mockTeam4)
	_, err4 := svc4.GetPredictions()
	assert.Error(t, err4)
}

func TestGetPredictions_FinishedLeague_Tie(t *testing.T) {
	mockStanding := new(MockStandingRepository)
	mockMatch := new(MockMatchRepository)

	standings := []model.Standing{
		{TeamID: 1, Points: 15, GoalDiff: 10, GoalsFor: 20, Team: model.Team{Name: "City"}},
		{TeamID: 2, Points: 15, GoalDiff: 10, GoalsFor: 20, Team: model.Team{Name: "Arsenal"}},
		{TeamID: 3, Points: 6, GoalDiff: -3, GoalsFor: 8, Team: model.Team{Name: "Liverpool"}},
	}
	playedMatches := []model.Match{{Played: true, Week: 6}} 

	mockStanding.On("GetAll").Return(standings, nil)
	mockMatch.On("GetAll").Return(playedMatches, nil)

	svc := NewPredictionService(mockStanding, mockMatch, nil)
	result, err := svc.GetPredictions()

	assert.NoError(t, err)
	assert.Equal(t, 50.0, result[0].Percentage) // City %50
	assert.Equal(t, 50.0, result[1].Percentage) // Arsenal %50
	assert.Equal(t, 0.0, result[2].Percentage)  // Liverpool %0
}

func TestGetPredictions_Calculate_MinScoreNegative(t *testing.T) {
	mockStanding := new(MockStandingRepository)
	mockMatch := new(MockMatchRepository)
	mockTeam := new(MockTeamRepository)

	standings := []model.Standing{
		{TeamID: 1, Points: 10, GoalDiff: 10, Team: model.Team{Name: "Good Team"}},
		{TeamID: 2, Points: 0, GoalDiff: -100, Team: model.Team{Name: "Bad Team"}}, 
	}
	
	teams := []model.Team{{ID: 1}, {ID: 2}}
	allMatches := []model.Match{{Played: true, Week: 4}}
	unplayed := []model.Match{{ID: 1, HomeTeamID: 1, AwayTeamID: 2}}

	mockStanding.On("GetAll").Return(standings, nil)
	mockMatch.On("GetAll").Return(allMatches, nil)
	mockTeam.On("GetAll").Return(teams, nil)
	mockMatch.On("GetUnplayed").Return(unplayed, nil)

	svc := NewPredictionService(mockStanding, mockMatch, mockTeam)
	result, err := svc.GetPredictions()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	
	total := result[0].Percentage + result[1].Percentage
	assert.InDelta(t, 100.0, total, 1.0)
}