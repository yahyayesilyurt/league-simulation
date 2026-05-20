package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yahyayesilyurt/league-simulation/internal/model"
)

type FixtureSvcMockMatchRepo struct { mock.Mock }
func (m *FixtureSvcMockMatchRepo) GetAll() ([]model.Match, error) { args := m.Called(); if args.Get(0) != nil { return args.Get(0).([]model.Match), args.Error(1) }; return nil, args.Error(1) }
func (m *FixtureSvcMockMatchRepo) GetByWeek(w int) ([]model.Match, error) { return nil, nil }
func (m *FixtureSvcMockMatchRepo) GetUnplayed() ([]model.Match, error) { return nil, nil }
func (m *FixtureSvcMockMatchRepo) GetByID(id uint) (*model.Match, error) { return nil, nil }
func (m *FixtureSvcMockMatchRepo) Create(match *model.Match) error { return m.Called(match).Error(0) }
func (m *FixtureSvcMockMatchRepo) Update(match *model.Match) error { return nil }
func (m *FixtureSvcMockMatchRepo) DeleteAll() error { return nil }

type FixtureSvcMockTeamRepo struct { mock.Mock }
func (m *FixtureSvcMockTeamRepo) GetAll() ([]model.Team, error) { args := m.Called(); if args.Get(0) != nil { return args.Get(0).([]model.Team), args.Error(1) }; return nil, args.Error(1) }
func (m *FixtureSvcMockTeamRepo) GetByID(id uint) (*model.Team, error) { return nil, nil }
func (m *FixtureSvcMockTeamRepo) Create(t *model.Team) error { return nil }
func (m *FixtureSvcMockTeamRepo) Update(t *model.Team) error { return nil }

func TestGenerateRoundRobin_CorrectMatchCount(t *testing.T) {
	teams := []model.Team{
		{ID: 1, Name: "Man City",  Strength: 90},
		{ID: 2, Name: "Liverpool", Strength: 85},
		{ID: 3, Name: "Arsenal",   Strength: 80},
		{ID: 4, Name: "Chelsea",   Strength: 75},
	}

	matches := generateRoundRobin(teams)

	assert.Len(t, matches, 12, "Should generate 12 matches")
}

func TestGenerateRoundRobin_CorrectWeekCount(t *testing.T) {
	teams := []model.Team{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4},
	}

	matches := generateRoundRobin(teams)

	weeks := make(map[int]bool)
	for _, m := range matches {
		weeks[m.Week] = true
	}

	assert.Len(t, weeks, 6, "Should have 6 weeks")
}

func TestGenerateRoundRobin_TwoMatchesPerWeek(t *testing.T) {
	teams := []model.Team{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4},
	}

	matches := generateRoundRobin(teams)

	weekMatches := make(map[int]int)
	for _, m := range matches {
		weekMatches[m.Week]++
	}

	for week, count := range weekMatches {
		assert.Equal(t, 2, count, "Week %d should have 2 matches", week)
	}
}

func TestGenerateRoundRobin_NoTeamPlaysItself(t *testing.T) {
	teams := []model.Team{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4},
	}

	matches := generateRoundRobin(teams)

	for _, m := range matches {
		assert.NotEqual(t, m.HomeTeamID, m.AwayTeamID,
			"A team should not play against itself")
	}
}

func TestGenerateRoundRobin_EachTeamPlaysSixMatches(t *testing.T) {
	teams := []model.Team{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4},
	}

	matches := generateRoundRobin(teams)

	matchCount := make(map[uint]int)
	for _, m := range matches {
		matchCount[m.HomeTeamID]++
		matchCount[m.AwayTeamID]++
	}

	for teamID, count := range matchCount {
		assert.Equal(t, 6, count,
			"Team %d should play 6 matches", teamID)
	}
}

func TestGenerateRoundRobin_SecondLegReversesHomeAway(t *testing.T) {
	teams := []model.Team{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4},
	}

	matches := generateRoundRobin(teams)

	firstLeg := make(map[[2]uint]bool)
	for _, m := range matches {
		if m.Week <= 3 {
			firstLeg[[2]uint{m.HomeTeamID, m.AwayTeamID}] = true
		}
	}

	for _, m := range matches {
		if m.Week > 3 {
			reversed := [2]uint{m.AwayTeamID, m.HomeTeamID}
			assert.True(t, firstLeg[reversed],
				"Second leg should reverse home/away teams")
		}
	}
}

func TestFixtureService_IsFixtureGenerated(t *testing.T) {
	mockMatch := new(FixtureSvcMockMatchRepo)
	svc := NewFixtureService(mockMatch, nil)

	mockMatch.On("GetAll").Return([]model.Match{{ID: 1}}, nil).Once()
	generated, err := svc.IsFixtureGenerated()
	assert.NoError(t, err)
	assert.True(t, generated)

	mockMatch.On("GetAll").Return([]model.Match{}, nil).Once()
	generated2, err2 := svc.IsFixtureGenerated()
	assert.NoError(t, err2)
	assert.False(t, generated2)

	mockMatch.On("GetAll").Return(nil, errors.New("db error")).Once()
	_, err3 := svc.IsFixtureGenerated()
	assert.Error(t, err3)
}

func TestFixtureService_GenerateFixture_Success(t *testing.T) {
	mockMatch := new(FixtureSvcMockMatchRepo)
	mockTeam := new(FixtureSvcMockTeamRepo)
	svc := NewFixtureService(mockMatch, mockTeam)

	mockMatch.On("GetAll").Return([]model.Match{}, nil).Once()
	
	teams := []model.Team{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}}
	mockTeam.On("GetAll").Return(teams, nil).Once()
	
	mockMatch.On("Create", mock.Anything).Return(nil).Times(12)

	err := svc.GenerateFixture()
	assert.NoError(t, err)
	mockMatch.AssertExpectations(t)
}

func TestFixtureService_GenerateFixture_Errors(t *testing.T) {
	mockMatch := new(FixtureSvcMockMatchRepo)
	mockTeam := new(FixtureSvcMockTeamRepo)
	svc := NewFixtureService(mockMatch, mockTeam)

	mockMatch.On("GetAll").Return([]model.Match{{ID: 1}}, nil).Once()
	err1 := svc.GenerateFixture()
	assert.Error(t, err1)
	assert.Contains(t, err1.Error(), "already generated")

	mockMatch.On("GetAll").Return([]model.Match{}, nil).Once()
	mockTeam.On("GetAll").Return(nil, errors.New("team db err")).Once()
	err2 := svc.GenerateFixture()
	assert.Error(t, err2)
	assert.Contains(t, err2.Error(), "failed to fetch teams")

	mockMatch.On("GetAll").Return([]model.Match{}, nil).Once()
	mockTeam.On("GetAll").Return([]model.Team{{ID: 1}, {ID: 2}}, nil).Once() 
	err3 := svc.GenerateFixture()
	assert.Error(t, err3)
	assert.Contains(t, err3.Error(), "expected 4 teams")

	mockMatch.On("GetAll").Return([]model.Match{}, nil).Once()
	mockTeam.On("GetAll").Return([]model.Team{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}}, nil).Once()
	mockMatch.On("Create", mock.Anything).Return(errors.New("insert failed")).Once()
	
	err4 := svc.GenerateFixture()
	assert.Error(t, err4)
	assert.Contains(t, err4.Error(), "failed to save match")
}