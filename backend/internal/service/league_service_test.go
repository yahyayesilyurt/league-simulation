package service

import (
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yahyayesilyurt/league-simulation/internal/cache"
	"github.com/yahyayesilyurt/league-simulation/internal/model"
)


type MockMatchRepo struct{ mock.Mock }
func (m *MockMatchRepo) GetAll() ([]model.Match, error) { args := m.Called(); return args.Get(0).([]model.Match), args.Error(1) }
func (m *MockMatchRepo) GetByWeek(w int) ([]model.Match, error) { args := m.Called(w); return args.Get(0).([]model.Match), args.Error(1) }
func (m *MockMatchRepo) GetUnplayed() ([]model.Match, error) { args := m.Called(); return args.Get(0).([]model.Match), args.Error(1) }
func (m *MockMatchRepo) GetByID(id uint) (*model.Match, error) { args := m.Called(id); return args.Get(0).(*model.Match), args.Error(1) }
func (m *MockMatchRepo) Create(match *model.Match) error { args := m.Called(match); return args.Error(0) }
func (m *MockMatchRepo) Update(match *model.Match) error { args := m.Called(match); return args.Error(0) }
func (m *MockMatchRepo) DeleteAll() error { args := m.Called(); return args.Error(0) }

type MockStandingRepo struct{ mock.Mock }
func (m *MockStandingRepo) GetAll() ([]model.Standing, error) { args := m.Called(); return args.Get(0).([]model.Standing), args.Error(1) }
func (m *MockStandingRepo) GetByTeamID(id uint) (*model.Standing, error) { args := m.Called(id); return args.Get(0).(*model.Standing), args.Error(1) }
func (m *MockStandingRepo) Update(s *model.Standing) error { args := m.Called(s); return args.Error(0) }
func (m *MockStandingRepo) ResetAll() error { args := m.Called(); return args.Error(0) }

type MockTeamRepo struct{ mock.Mock }
func (m *MockTeamRepo) GetAll() ([]model.Team, error) { args := m.Called(); return args.Get(0).([]model.Team), args.Error(1) }
func (m *MockTeamRepo) GetByID(id uint) (*model.Team, error) { args := m.Called(id); return args.Get(0).(*model.Team), args.Error(1) }
func (m *MockTeamRepo) Create(t *model.Team) error { args := m.Called(t); return args.Error(0) }
func (m *MockTeamRepo) Update(t *model.Team) error { args := m.Called(t); return args.Error(0) }

type MockMatchSvc struct{ mock.Mock }
func (m *MockMatchSvc) PlayWeek(w int) ([]model.Match, error) { args := m.Called(w); return args.Get(0).([]model.Match), args.Error(1) }
func (m *MockMatchSvc) PlayMatch(id uint) (*model.Match, error) { args := m.Called(id); return args.Get(0).(*model.Match), args.Error(1) }
func (m *MockMatchSvc) UpdateMatchResult(id uint, h, a int) (*model.Match, error) { args := m.Called(id, h, a); return args.Get(0).(*model.Match), args.Error(1) }
func (m *MockMatchSvc) GetMatchByID(id uint) (*model.Match, error) { args := m.Called(id); return args.Get(0).(*model.Match), args.Error(1) }

type MockPredictionSvc struct{ mock.Mock }
func (m *MockPredictionSvc) GetPredictions() ([]TeamPrediction, error) { args := m.Called(); return args.Get(0).([]TeamPrediction), args.Error(1) }


func setupLeagueServiceTest(t *testing.T) (*leagueService, *MockMatchRepo, *MockStandingRepo, *MockTeamRepo, *MockMatchSvc, *MockPredictionSvc, *miniredis.Miniredis) {
	mr, _ := miniredis.Run()
	t.Cleanup(mr.Close)
	
	redisClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	appCache := cache.NewCache(redisClient)

	matchRepo := new(MockMatchRepo)
	standingRepo := new(MockStandingRepo)
	teamRepo := new(MockTeamRepo)
	matchSvc := new(MockMatchSvc)
	predSvc := new(MockPredictionSvc)

	svc := &leagueService{
		matchRepo:     matchRepo,
		standingRepo:  standingRepo,
		teamRepo:      teamRepo,
		matchSvc:      matchSvc,
		predictionSvc: predSvc,
		cache:         appCache,
	}

	return svc, matchRepo, standingRepo, teamRepo, matchSvc, predSvc, mr
}

func TestLeagueService_GetStandings(t *testing.T) {
	svc, _, mockStandingRepo, _, _, _, _ := setupLeagueServiceTest(t)

	expectedStandings := []model.Standing{{Points: 10}}
	mockStandingRepo.On("GetAll").Return(expectedStandings, nil).Once()

	standings, err := svc.GetStandings()
	assert.NoError(t, err)
	assert.Equal(t, 10, standings[0].Points)

	standingsCached, err := svc.GetStandings()
	assert.NoError(t, err)
	assert.Equal(t, 10, standingsCached[0].Points)

	mockStandingRepo.AssertExpectations(t)
}

func TestLeagueService_GetStatus(t *testing.T) {
	svc, mockMatchRepo, _, _, _, _, mr := setupLeagueServiceTest(t)

	matches := []model.Match{
		{Played: true, Week: 1}, {Played: true, Week: 1}, {Played: true, Week: 2},
		{Played: false}, {Played: false}, {Played: false}, {Played: false}, {Played: false}, {Played: false}, {Played: false}, {Played: false}, {Played: false},
	}
	mockMatchRepo.On("GetAll").Return(matches, nil)

	status, err := svc.GetStatus()
	assert.NoError(t, err)
	assert.Equal(t, "in_progress", status.Status)
	assert.Equal(t, 2, status.CurrentWeek) 
	assert.Equal(t, 3, status.MatchesPlayed)
	assert.Equal(t, 9, status.MatchesLeft)
	
	mr.FlushAll()

	finishedMatches := []model.Match{{Played: true}, {Played: true}}
	mockMatchRepo.ExpectedCalls = nil 
	mockMatchRepo.On("GetAll").Return(finishedMatches, nil)

	statusFinished, _ := svc.GetStatus()
	assert.Equal(t, "finished", statusFinished.Status)
	assert.True(t, statusFinished.LeagueFinished)
}

func TestLeagueService_NextWeek(t *testing.T) {
	svc, mockMatchRepo, mockStandingRepo, _, mockMatchSvc, mockPredSvc, _ := setupLeagueServiceTest(t)

	mockMatchRepo.On("GetAll").Return([]model.Match{{Played: true, Week: 2}}, nil)
	
	mockMatchSvc.On("PlayWeek", 3).Return([]model.Match{{Week: 3}}, nil)
	mockStandingRepo.On("GetAll").Return([]model.Standing{}, nil)
	mockPredSvc.On("GetPredictions").Return([]TeamPrediction{}, nil)

	result, err := svc.NextWeek()

	assert.NoError(t, err)
	assert.Equal(t, 3, result.Week)
	assert.False(t, result.LeagueFinished) 

	mockMatchRepo.ExpectedCalls = nil
	mockMatchRepo.On("GetAll").Return([]model.Match{{Played: true, Week: 6}}, nil)
	
	_, errFinished := svc.NextWeek()
	assert.Error(t, errFinished)
	assert.Contains(t, errFinished.Error(), "league is finished")
}

func TestLeagueService_PlayAll(t *testing.T) {
	svc, mockMatchRepo, mockStandingRepo, _, mockMatchSvc, mockPredSvc, _ := setupLeagueServiceTest(t)

	mockMatchRepo.On("GetAll").Return([]model.Match{{Played: true, Week: 4}}, nil)

	mockMatchSvc.On("PlayWeek", 5).Return([]model.Match{}, nil)
	mockMatchSvc.On("PlayWeek", 6).Return([]model.Match{}, nil)
	
	standings := []model.Standing{
		{Team: model.Team{Name: "City"}, GoalsFor: 15, GoalsAgainst: 3}, 
		{Team: model.Team{Name: "Chelsea"}, GoalsFor: 5, GoalsAgainst: 10},
	}
	mockStandingRepo.On("GetAll").Return(standings, nil)
	mockPredSvc.On("GetPredictions").Return([]TeamPrediction{}, nil)

	result, err := svc.PlayAll()

	assert.NoError(t, err)
	assert.Equal(t, 2, result.TotalWeeksPlayed) 
	assert.Equal(t, "City", result.Summary.TopScorer)
	assert.Equal(t, "City", result.Summary.BestDefense)
}

func TestLeagueService_Reset(t *testing.T) {
	svc, mockMatchRepo, mockStandingRepo, mockTeamRepo, _, _, _ := setupLeagueServiceTest(t)

	mockMatchRepo.On("DeleteAll").Return(nil)
	mockStandingRepo.On("ResetAll").Return(nil)
	
	mockMatchRepo.On("GetAll").Return([]model.Match{}, nil)
	
	mockTeamRepo.On("GetAll").Return([]model.Team{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}}, nil)
	mockMatchRepo.On("Create", mock.Anything).Return(nil)

	status, err := svc.Reset()

	assert.NoError(t, err)
	assert.Equal(t, "not_started", status.Status)
	assert.Equal(t, 0, status.CurrentWeek)
	assert.Equal(t, 12, status.MatchesLeft)

	mockMatchRepo.ExpectedCalls = nil 
	mockMatchRepo.On("DeleteAll").Return(errors.New("db crash"))
	
	_, errDelete := svc.Reset()
	assert.Error(t, errDelete)
	assert.Contains(t, errDelete.Error(), "failed to delete matches")
}