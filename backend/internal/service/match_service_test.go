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


type MatchSvcMockMatchRepo struct{ mock.Mock }
func (m *MatchSvcMockMatchRepo) GetAll() ([]model.Match, error) { args := m.Called(); return args.Get(0).([]model.Match), args.Error(1) }
func (m *MatchSvcMockMatchRepo) GetByWeek(w int) ([]model.Match, error) { args := m.Called(w); if args.Get(0) != nil { return args.Get(0).([]model.Match), args.Error(1) }; return nil, args.Error(1) }
func (m *MatchSvcMockMatchRepo) GetUnplayed() ([]model.Match, error) { args := m.Called(); return args.Get(0).([]model.Match), args.Error(1) }
func (m *MatchSvcMockMatchRepo) GetByID(id uint) (*model.Match, error) { args := m.Called(id); if args.Get(0) != nil { return args.Get(0).(*model.Match), args.Error(1) }; return nil, args.Error(1) }
func (m *MatchSvcMockMatchRepo) Create(match *model.Match) error { args := m.Called(match); return args.Error(0) }
func (m *MatchSvcMockMatchRepo) Update(match *model.Match) error { args := m.Called(match); return args.Error(0) }
func (m *MatchSvcMockMatchRepo) DeleteAll() error { args := m.Called(); return args.Error(0) }

type MatchSvcMockStandingRepo struct{ mock.Mock }
func (m *MatchSvcMockStandingRepo) GetAll() ([]model.Standing, error) { args := m.Called(); return args.Get(0).([]model.Standing), args.Error(1) }
func (m *MatchSvcMockStandingRepo) GetByTeamID(id uint) (*model.Standing, error) { args := m.Called(id); if args.Get(0) != nil { return args.Get(0).(*model.Standing), args.Error(1) }; return nil, args.Error(1) }
func (m *MatchSvcMockStandingRepo) Update(s *model.Standing) error { args := m.Called(s); return args.Error(0) }
func (m *MatchSvcMockStandingRepo) ResetAll() error { args := m.Called(); return args.Error(0) }

type MatchSvcMockTeamRepo struct{ mock.Mock }
func (m *MatchSvcMockTeamRepo) GetAll() ([]model.Team, error) { args := m.Called(); return args.Get(0).([]model.Team), args.Error(1) }
func (m *MatchSvcMockTeamRepo) GetByID(id uint) (*model.Team, error) { args := m.Called(id); return args.Get(0).(*model.Team), args.Error(1) }
func (m *MatchSvcMockTeamRepo) Create(t *model.Team) error { args := m.Called(t); return args.Error(0) }
func (m *MatchSvcMockTeamRepo) Update(t *model.Team) error { args := m.Called(t); return args.Error(0) }

type MatchSvcMockStandingSvc struct{ mock.Mock }
func (m *MatchSvcMockStandingSvc) RecalculateAll() error { return m.Called().Error(0) }
func (m *MatchSvcMockStandingSvc) GetStandings() ([]model.Standing, error) { return nil, nil }


func setupMatchServiceTest(t *testing.T) (*matchService, *MatchSvcMockMatchRepo, *MatchSvcMockStandingRepo, *MatchSvcMockStandingSvc) {
	mr, _ := miniredis.Run()
	t.Cleanup(mr.Close)
	
	redisClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	appCache := cache.NewCache(redisClient)

	matchRepo := new(MatchSvcMockMatchRepo)
	standingRepo := new(MatchSvcMockStandingRepo)
	teamRepo := new(MatchSvcMockTeamRepo)
	standingSvc := new(MatchSvcMockStandingSvc)

	svc := &matchService{
		matchRepo:    matchRepo,
		standingRepo: standingRepo,
		teamRepo:     teamRepo,
		engine:       NewSimulationEngine(),
		standingSvc:  standingSvc,
		cache:        appCache,
	}

	return svc, matchRepo, standingRepo, standingSvc
}


func TestMatchService_GetMatchByID(t *testing.T) {
	svc, mockRepo, _, _ := setupMatchServiceTest(t)
	mockRepo.On("GetByID", uint(1)).Return(&model.Match{ID: 1}, nil)
	
	match, err := svc.GetMatchByID(1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), match.ID)
}


func TestMatchService_UpdateMatchResult(t *testing.T) {
	svc, mockRepo, _, mockStandingSvc := setupMatchServiceTest(t)

	playedMatch := &model.Match{ID: 1, Played: true}
	mockRepo.On("GetByID", uint(1)).Return(playedMatch, nil).Once()
	mockRepo.On("Update", mock.Anything).Return(nil).Once()
	mockStandingSvc.On("RecalculateAll").Return(nil).Once()

	match, err := svc.UpdateMatchResult(1, 3, 0)
	assert.NoError(t, err)
	assert.Equal(t, 3, *match.HomeGoals)
	assert.Equal(t, 0, *match.AwayGoals)

	mockRepo.On("GetByID", uint(99)).Return(nil, errors.New("not found")).Once()
	_, errNotFound := svc.UpdateMatchResult(99, 1, 1)
	assert.Error(t, errNotFound)

	unplayedMatch := &model.Match{ID: 2, Played: false}
	mockRepo.On("GetByID", uint(2)).Return(unplayedMatch, nil).Once()
	_, errNotPlayed := svc.UpdateMatchResult(2, 1, 1)
	assert.Error(t, errNotPlayed)
	assert.Contains(t, errNotPlayed.Error(), "not been played yet")

	mockRepo.On("GetByID", uint(3)).Return(&model.Match{ID: 3, Played: true}, nil).Once()
	mockRepo.On("Update", mock.Anything).Return(errors.New("db error")).Once()
	_, errUpdate := svc.UpdateMatchResult(3, 1, 1)
	assert.Error(t, errUpdate)

	mockRepo.On("GetByID", uint(4)).Return(&model.Match{ID: 4, Played: true}, nil).Once()
	mockRepo.On("Update", mock.Anything).Return(nil).Once()
	mockStandingSvc.On("RecalculateAll").Return(errors.New("calc error")).Once()
	_, errCalc := svc.UpdateMatchResult(4, 1, 1)
	assert.Error(t, errCalc)
}


func TestMatchService_PlayMatch_Errors(t *testing.T) {
	svc, mockRepo, _, _ := setupMatchServiceTest(t)

	mockRepo.On("GetByID", uint(99)).Return(nil, errors.New("not found")).Once()
	_, err1 := svc.PlayMatch(99)
	assert.Error(t, err1)

	mockRepo.On("GetByID", uint(1)).Return(&model.Match{ID: 1, Played: true}, nil).Once()
	_, err2 := svc.PlayMatch(1)
	assert.Error(t, err2)
	assert.Contains(t, err2.Error(), "already played")
}

func TestMatchService_PlayWeek(t *testing.T) {
	svc, mockRepo, mockStandingRepo, _ := setupMatchServiceTest(t)

	mockRepo.On("GetByWeek", 1).Return(nil, errors.New("db err")).Once()
	_, err1 := svc.PlayWeek(1)
	assert.Error(t, err1)

	mockRepo.On("GetByWeek", 2).Return([]model.Match{}, nil).Once()
	_, err2 := svc.PlayWeek(2)
	assert.Error(t, err2)

	playedMatches := []model.Match{{Played: true}, {Played: true}}
	mockRepo.On("GetByWeek", 3).Return(playedMatches, nil).Once()
	_, err3 := svc.PlayWeek(3)
	assert.Error(t, err3)
	assert.Contains(t, err3.Error(), "already played")

	mixedMatches := []model.Match{
		{ID: 10, Played: true},
		{ID: 11, Played: false, HomeTeamID: 1, AwayTeamID: 2},
	}
	mockRepo.On("GetByWeek", 4).Return(mixedMatches, nil).Once()
	
	mockRepo.On("GetByID", uint(11)).Return(&mixedMatches[1], nil).Once()
	mockStandingRepo.On("GetByTeamID", uint(1)).Return(&model.Standing{}, nil).Once()
	mockStandingRepo.On("GetByTeamID", uint(2)).Return(&model.Standing{}, nil).Once()
	mockRepo.On("Update", mock.Anything).Return(nil).Once()
	mockStandingRepo.On("Update", mock.Anything).Return(nil).Twice()

	results, errSuccess := svc.PlayWeek(4)
	assert.NoError(t, errSuccess)
	assert.Len(t, results, 1) 
}


func TestMatchService_UpdateStandings_Logic(t *testing.T) {
	svc, _, mockStandingRepo, _ := setupMatchServiceTest(t)

	testCases := []struct {
		name       string
		homeGoals  int
		awayGoals  int
		expectHome func(*model.Standing)
		expectAway func(*model.Standing)
	}{
		{
			name:      "Home Team Wins",
			homeGoals: 3, awayGoals: 1,
			expectHome: func(s *model.Standing) { assert.Equal(t, 1, s.Won); assert.Equal(t, 3, s.Points) },
			expectAway: func(s *model.Standing) { assert.Equal(t, 1, s.Lost); assert.Equal(t, 0, s.Points) },
		},
		{
			name:      "Away Team Wins",
			homeGoals: 0, awayGoals: 2,
			expectHome: func(s *model.Standing) { assert.Equal(t, 1, s.Lost); assert.Equal(t, 0, s.Points) },
			expectAway: func(s *model.Standing) { assert.Equal(t, 1, s.Won); assert.Equal(t, 3, s.Points) },
		},
		{
			name:      "Drawn",
			homeGoals: 2, awayGoals: 2,
			expectHome: func(s *model.Standing) { assert.Equal(t, 1, s.Drawn); assert.Equal(t, 1, s.Points) },
			expectAway: func(s *model.Standing) { assert.Equal(t, 1, s.Drawn); assert.Equal(t, 1, s.Points) },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match := &model.Match{
				HomeTeamID: 1,
				AwayTeamID: 2,
				HomeGoals:  &tc.homeGoals,
				AwayGoals:  &tc.awayGoals,
			}

			homeSt := &model.Standing{TeamID: 1}
			awaySt := &model.Standing{TeamID: 2}

			mockStandingRepo.On("GetByTeamID", uint(1)).Return(homeSt, nil).Once()
			mockStandingRepo.On("GetByTeamID", uint(2)).Return(awaySt, nil).Once()
			mockStandingRepo.On("Update", mock.Anything).Return(nil).Twice()

			err := svc.updateStandings(match)

			assert.NoError(t, err)
			tc.expectHome(homeSt)
			tc.expectAway(awaySt)
			
			assert.Equal(t, 1, homeSt.Played)
			assert.Equal(t, tc.homeGoals, homeSt.GoalsFor)
			assert.Equal(t, tc.homeGoals-tc.awayGoals, homeSt.GoalDiff)
		})
	}
}

func TestMatchService_UpdateStandings_Errors(t *testing.T) {
	svc, _, mockStandingRepo, _ := setupMatchServiceTest(t)
	
	hg, ag := 1, 1
	match := &model.Match{HomeTeamID: 1, AwayTeamID: 2, HomeGoals: &hg, AwayGoals: &ag}

	mockStandingRepo.On("GetByTeamID", uint(1)).Return(nil, errors.New("not found")).Once()
	assert.Error(t, svc.updateStandings(match))

	mockStandingRepo.On("GetByTeamID", uint(1)).Return(&model.Standing{}, nil).Once()
	mockStandingRepo.On("GetByTeamID", uint(2)).Return(nil, errors.New("not found")).Once()
	assert.Error(t, svc.updateStandings(match))
}