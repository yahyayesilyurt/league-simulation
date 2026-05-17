package service

import (
	"fmt"

	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"github.com/yahyayesilyurt/league-simulation/internal/repository"
)

type LeagueService interface {
	GetStandings() ([]model.Standing, error)
	GetFixtures() ([]model.Match, error)
	GetWeek(week int) ([]model.Match, error)
	GetCurrentWeek() (int, error)
	NextWeek() ([]model.Match, error)
	PlayAll() (map[int][]model.Match, error)
	Reset() error
}

type leagueService struct {
	matchRepo    repository.MatchRepository
	standingRepo repository.StandingRepository
	teamRepo     repository.TeamRepository
	matchSvc     MatchService
}

func NewLeagueService(
	matchRepo repository.MatchRepository,
	standingRepo repository.StandingRepository,
	teamRepo repository.TeamRepository,
) LeagueService {
	return &leagueService{
		matchRepo:    matchRepo,
		standingRepo: standingRepo,
		teamRepo:     teamRepo,
		matchSvc:     NewMatchService(matchRepo, standingRepo, teamRepo),
	}
}

func (s *leagueService) GetStandings() ([]model.Standing, error) {
	return s.standingRepo.GetAll()
}

func (s *leagueService) GetFixtures() ([]model.Match, error) {
	return s.matchRepo.GetAll()
}

func (s *leagueService) GetWeek(week int) ([]model.Match, error) {
	return s.matchRepo.GetByWeek(week)
}

func (s *leagueService) GetCurrentWeek() (int, error) {
	matches, err := s.matchRepo.GetAll()
	if err != nil {
		return 0, err
	}
	currentWeek := 0
	for _, m := range matches {
		if m.Played && m.Week > currentWeek {
			currentWeek = m.Week
		}
	}
	return currentWeek, nil
}

func (s *leagueService) NextWeek() ([]model.Match, error) {
	currentWeek, err := s.GetCurrentWeek()
	if err != nil {
		return nil, err
	}

	nextWeek := currentWeek + 1
	if nextWeek > 6 {
		return nil, fmt.Errorf("league is finished, all 6 weeks have been played")
	}

	return s.matchSvc.PlayWeek(nextWeek)
}

func (s *leagueService) PlayAll() (map[int][]model.Match, error) {
	currentWeek, err := s.GetCurrentWeek()
	if err != nil {
		return nil, err
	}

	if currentWeek >= 6 {
		return nil, fmt.Errorf("league is already finished")
	}

	results := make(map[int][]model.Match)
	for week := currentWeek + 1; week <= 6; week++ {
		matches, err := s.matchSvc.PlayWeek(week)
		if err != nil {
			return nil, fmt.Errorf("error playing week %d: %w", week, err)
		}
		results[week] = matches
	}

	return results, nil
}

func (s *leagueService) Reset() error {
	if err := s.matchRepo.DeleteAll(); err != nil {
		return err
	}
	if err := s.standingRepo.ResetAll(); err != nil {
		return err
	}

	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return err
	}

	fixtureSvc := NewFixtureService(s.matchRepo, s.teamRepo)
	_ = teams
	return fixtureSvc.GenerateFixture()
}