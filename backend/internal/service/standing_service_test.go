package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yahyayesilyurt/league-simulation/internal/model"
)

// Mock Repositories

type MockStandingRepository struct {
	mock.Mock
}

func (m *MockStandingRepository) GetAll() ([]model.Standing, error) {
	args := m.Called()
	return args.Get(0).([]model.Standing), args.Error(1)
}

func (m *MockStandingRepository) GetByTeamID(teamID uint) (*model.Standing, error) {
	args := m.Called(teamID)
	return args.Get(0).(*model.Standing), args.Error(1)
}

func (m *MockStandingRepository) Update(s *model.Standing) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *MockStandingRepository) ResetAll() error {
	args := m.Called()
	return args.Error(0)
}

type MockMatchRepository struct {
	mock.Mock
}

func (m *MockMatchRepository) GetAll() ([]model.Match, error) {
	args := m.Called()
	return args.Get(0).([]model.Match), args.Error(1)
}

func (m *MockMatchRepository) GetByWeek(week int) ([]model.Match, error) {
	args := m.Called(week)
	return args.Get(0).([]model.Match), args.Error(1)
}

func (m *MockMatchRepository) GetUnplayed() ([]model.Match, error) {
	args := m.Called()
	return args.Get(0).([]model.Match), args.Error(1)
}

func (m *MockMatchRepository) GetByID(id uint) (*model.Match, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Match), args.Error(1)
}

func (m *MockMatchRepository) Create(match *model.Match) error {
	args := m.Called(match)
	return args.Error(0)
}

func (m *MockMatchRepository) Update(match *model.Match) error {
	args := m.Called(match)
	return args.Error(0)
}

func (m *MockMatchRepository) DeleteAll() error {
	args := m.Called()
	return args.Error(0)
}

type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) GetAll() ([]model.Team, error) {
	args := m.Called()
	return args.Get(0).([]model.Team), args.Error(1)
}

func (m *MockTeamRepository) GetByID(id uint) (*model.Team, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Team), args.Error(1)
}

func (m *MockTeamRepository) Create(team *model.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

func (m *MockTeamRepository) Update(team *model.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

// Tests

func TestGetStandings_ReturnsSortedByPoints(t *testing.T) {
	mockStanding := new(MockStandingRepository)
	mockMatch    := new(MockMatchRepository)
	mockTeam     := new(MockTeamRepository)

	standings := []model.Standing{
		{TeamID: 1, Points: 6,  GoalDiff: 3, Team: model.Team{Name: "Arsenal"}},
		{TeamID: 2, Points: 9,  GoalDiff: 5, Team: model.Team{Name: "Man City"}},
		{TeamID: 3, Points: 3,  GoalDiff: -2, Team: model.Team{Name: "Chelsea"}},
		{TeamID: 4, Points: 7,  GoalDiff: 1, Team: model.Team{Name: "Liverpool"}},
	}

	mockStanding.On("GetAll").Return(standings, nil)

	svc := NewStandingService(mockStanding, mockMatch, mockTeam)
	result, err := svc.GetStandings()

	assert.NoError(t, err)
	assert.Len(t, result, 4)
	mockStanding.AssertExpectations(t)
}

func TestRecalculateAll_HomeWin(t *testing.T) {
	mockStanding := new(MockStandingRepository)
	mockMatch    := new(MockMatchRepository)
	mockTeam     := new(MockTeamRepository)

	homeGoals := 3
	awayGoals := 1

	matches := []model.Match{
		{
			ID:         1,
			Week:       1,
			HomeTeamID: 1,
			AwayTeamID: 2,
			HomeGoals:  &homeGoals,
			AwayGoals:  &awayGoals,
			Played:     true,
		},
	}

	teams := []model.Team{
		{ID: 1, Name: "Man City"},
		{ID: 2, Name: "Liverpool"},
	}

	homeStanding := &model.Standing{TeamID: 1}
	awayStanding := &model.Standing{TeamID: 2}

	mockTeam.On("GetAll").Return(teams, nil)
	mockMatch.On("GetAll").Return(matches, nil)
	mockStanding.On("ResetAll").Return(nil)
	mockStanding.On("GetByTeamID", uint(1)).Return(homeStanding, nil)
	mockStanding.On("GetByTeamID", uint(2)).Return(awayStanding, nil)
	mockStanding.On("Update", mock.AnythingOfType("*model.Standing")).Return(nil)

	svc := NewStandingService(mockStanding, mockMatch, mockTeam)
	err := svc.RecalculateAll()

	assert.NoError(t, err)

	assert.Equal(t, 3, homeStanding.Points)
	assert.Equal(t, 1, homeStanding.Won)
	assert.Equal(t, 0, homeStanding.Lost)
	assert.Equal(t, 3, homeStanding.GoalsFor)
	assert.Equal(t, 1, homeStanding.GoalsAgainst)
	assert.Equal(t, 2, homeStanding.GoalDiff)

	assert.Equal(t, 0, awayStanding.Points)
	assert.Equal(t, 0, awayStanding.Won)
	assert.Equal(t, 1, awayStanding.Lost)
	assert.Equal(t, 1, awayStanding.GoalsFor)
	assert.Equal(t, 3, awayStanding.GoalsAgainst)
	assert.Equal(t, -2, awayStanding.GoalDiff)

	mockStanding.AssertExpectations(t)
}

func TestRecalculateAll_Draw(t *testing.T) {
	mockStanding := new(MockStandingRepository)
	mockMatch    := new(MockMatchRepository)
	mockTeam     := new(MockTeamRepository)

	homeGoals := 2
	awayGoals := 2

	matches := []model.Match{
		{
			ID:         1,
			HomeTeamID: 1,
			AwayTeamID: 2,
			HomeGoals:  &homeGoals,
			AwayGoals:  &awayGoals,
			Played:     true,
		},
	}

	teams        := []model.Team{{ID: 1}, {ID: 2}}
	homeStanding := &model.Standing{TeamID: 1}
	awayStanding := &model.Standing{TeamID: 2}

	mockTeam.On("GetAll").Return(teams, nil)
	mockMatch.On("GetAll").Return(matches, nil)
	mockStanding.On("ResetAll").Return(nil)
	mockStanding.On("GetByTeamID", uint(1)).Return(homeStanding, nil)
	mockStanding.On("GetByTeamID", uint(2)).Return(awayStanding, nil)
	mockStanding.On("Update", mock.AnythingOfType("*model.Standing")).Return(nil)

	svc := NewStandingService(mockStanding, mockMatch, mockTeam)
	err := svc.RecalculateAll()

	assert.NoError(t, err)

	assert.Equal(t, 1, homeStanding.Points)
	assert.Equal(t, 1, homeStanding.Drawn)
	assert.Equal(t, 1, awayStanding.Points)
	assert.Equal(t, 1, awayStanding.Drawn)
}

func TestRecalculateAll_AwayWin(t *testing.T) {
	mockStanding := new(MockStandingRepository)
	mockMatch    := new(MockMatchRepository)
	mockTeam     := new(MockTeamRepository)

	homeGoals := 0
	awayGoals := 2

	matches := []model.Match{
		{
			ID:         1,
			HomeTeamID: 1,
			AwayTeamID: 2,
			HomeGoals:  &homeGoals,
			AwayGoals:  &awayGoals,
			Played:     true,
		},
	}

	teams        := []model.Team{{ID: 1}, {ID: 2}}
	homeStanding := &model.Standing{TeamID: 1}
	awayStanding := &model.Standing{TeamID: 2}

	mockTeam.On("GetAll").Return(teams, nil)
	mockMatch.On("GetAll").Return(matches, nil)
	mockStanding.On("ResetAll").Return(nil)
	mockStanding.On("GetByTeamID", uint(1)).Return(homeStanding, nil)
	mockStanding.On("GetByTeamID", uint(2)).Return(awayStanding, nil)
	mockStanding.On("Update", mock.AnythingOfType("*model.Standing")).Return(nil)

	svc := NewStandingService(mockStanding, mockMatch, mockTeam)
	err := svc.RecalculateAll()

	assert.NoError(t, err)

	assert.Equal(t, 3, awayStanding.Points)
	assert.Equal(t, 1, awayStanding.Won)
	assert.Equal(t, 0, homeStanding.Points)
	assert.Equal(t, 1, homeStanding.Lost)
}

func TestRecalculateAll_UnplayedMatchesIgnored(t *testing.T) {
	mockStanding := new(MockStandingRepository)
	mockMatch    := new(MockMatchRepository)
	mockTeam     := new(MockTeamRepository)

	matches := []model.Match{
		{
			ID:         1,
			HomeTeamID: 1,
			AwayTeamID: 2,
			HomeGoals:  nil,
			AwayGoals:  nil,
			Played:     false,
		},
	}

	teams        := []model.Team{{ID: 1}, {ID: 2}}
	homeStanding := &model.Standing{TeamID: 1}
	awayStanding := &model.Standing{TeamID: 2}

	mockTeam.On("GetAll").Return(teams, nil)
	mockMatch.On("GetAll").Return(matches, nil)
	mockStanding.On("ResetAll").Return(nil)
	mockStanding.On("GetByTeamID", uint(1)).Return(homeStanding, nil)
	mockStanding.On("GetByTeamID", uint(2)).Return(awayStanding, nil)
	mockStanding.On("Update", mock.AnythingOfType("*model.Standing")).Return(nil)

	svc := NewStandingService(mockStanding, mockMatch, mockTeam)
	err := svc.RecalculateAll()

	assert.NoError(t, err)

	assert.Equal(t, 0, homeStanding.Points)
	assert.Equal(t, 0, awayStanding.Points)
	assert.Equal(t, 0, homeStanding.Played)
	assert.Equal(t, 0, awayStanding.Played)
}

func TestGetStandings_Error(t *testing.T) {
	mockStanding := new(MockStandingRepository)
	mockStanding.On("GetAll").Return([]model.Standing{}, errors.New("db error"))

	svc := NewStandingService(mockStanding, nil, nil)
	_, err := svc.GetStandings()

	assert.Error(t, err)
}

func TestRecalculateAll_Errors(t *testing.T) {
	mockTeam1 := new(MockTeamRepository)
	mockTeam1.On("GetAll").Return([]model.Team{}, errors.New("team fetch error"))
	svc1 := NewStandingService(nil, nil, mockTeam1)
	assert.Error(t, svc1.RecalculateAll())

	mockTeam2 := new(MockTeamRepository)
	mockTeam2.On("GetAll").Return([]model.Team{{ID: 1}}, nil)
	mockStanding2 := new(MockStandingRepository)
	mockStanding2.On("ResetAll").Return(errors.New("reset error"))
	svc2 := NewStandingService(mockStanding2, nil, mockTeam2)
	assert.Error(t, svc2.RecalculateAll())

	mockTeam3 := new(MockTeamRepository)
	mockTeam3.On("GetAll").Return([]model.Team{{ID: 1}}, nil)
	mockStanding3 := new(MockStandingRepository)
	mockStanding3.On("ResetAll").Return(nil)
	mockMatch3 := new(MockMatchRepository)
	mockMatch3.On("GetAll").Return([]model.Match{}, errors.New("match fetch error"))
	svc3 := NewStandingService(mockStanding3, mockMatch3, mockTeam3)
	assert.Error(t, svc3.RecalculateAll())

	mockTeam4 := new(MockTeamRepository)
	mockTeam4.On("GetAll").Return([]model.Team{{ID: 1}}, nil)
	mockStanding4 := new(MockStandingRepository)
	mockStanding4.On("ResetAll").Return(nil)
	mockMatch4 := new(MockMatchRepository)
	mockMatch4.On("GetAll").Return([]model.Match{}, nil)
	mockStanding4.On("GetByTeamID", uint(1)).Return((*model.Standing)(nil), errors.New("standing fetch error"))
	svc4 := NewStandingService(mockStanding4, mockMatch4, mockTeam4)
	assert.Error(t, svc4.RecalculateAll())

	mockTeam5 := new(MockTeamRepository)
	mockTeam5.On("GetAll").Return([]model.Team{{ID: 1}}, nil)
	mockStanding5 := new(MockStandingRepository)
	mockStanding5.On("ResetAll").Return(nil)
	mockMatch5 := new(MockMatchRepository)
	mockMatch5.On("GetAll").Return([]model.Match{}, nil) 
	
	mockStanding5.On("GetByTeamID", uint(1)).Return(&model.Standing{TeamID: 1}, nil)
	mockStanding5.On("Update", mock.AnythingOfType("*model.Standing")).Return(errors.New("update error"))
	
	svc5 := NewStandingService(mockStanding5, mockMatch5, mockTeam5)
	assert.Error(t, svc5.RecalculateAll())
}