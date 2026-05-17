package service

import (
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
	matchRepo   repository.MatchRepository
	standingRepo repository.StandingRepository
	teamRepo    repository.TeamRepository
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
	return nil, nil
}

func (s *leagueService) PlayAll() (map[int][]model.Match, error) {
	return nil, nil
}

func (s *leagueService) Reset() error {
	return nil
}