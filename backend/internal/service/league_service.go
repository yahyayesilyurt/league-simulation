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
	NextWeek() (*model.WeekResult, error)
	PlayAll() ([]model.WeekResult, error)
	Reset() error
}

type leagueService struct {
	matchRepo    repository.MatchRepository
	standingRepo repository.StandingRepository
	teamRepo     repository.TeamRepository
	matchSvc     MatchService
	predictionSvc PredictionService
}

func NewLeagueService(
	matchRepo repository.MatchRepository,
	standingRepo repository.StandingRepository,
	teamRepo repository.TeamRepository,
) LeagueService {
	return &leagueService{
		matchRepo:     matchRepo,
		standingRepo:  standingRepo,
		teamRepo:      teamRepo,
		matchSvc:      NewMatchService(matchRepo, standingRepo, teamRepo),
		predictionSvc: NewPredictionService(standingRepo, matchRepo, teamRepo),
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

func (s *leagueService) NextWeek() (*model.WeekResult, error) {
	currentWeek, err := s.GetCurrentWeek()
	if err != nil {
		return nil, err
	}

	nextWeek := currentWeek + 1
	if nextWeek > 6 {
		return nil, fmt.Errorf("league is finished, all 6 weeks have been played")
	}

	matches, err := s.matchSvc.PlayWeek(nextWeek)
	if err != nil {
		return nil, err
	}

	standings, err := s.standingRepo.GetAll()
	if err != nil {
		return nil, err
	}

	predictions, err := s.predictionSvc.GetPredictions()
	if err != nil {
		return nil, err
	}

	return &model.WeekResult{
		Week:           nextWeek,
		Matches:        matches,
		Standings:      standings,
		Predictions:    predictions,
		LeagueFinished: nextWeek == 6,
	}, nil
}

func (s *leagueService) PlayAll() ([]model.WeekResult, error) {
	currentWeek, err := s.GetCurrentWeek()
	if err != nil {
		return nil, err
	}
	if currentWeek >= 6 {
		return nil, fmt.Errorf("league is already finished")
	}

	var results []model.WeekResult
	for week := currentWeek + 1; week <= 6; week++ {
		matches, err := s.matchSvc.PlayWeek(week)
		if err != nil {
			return nil, fmt.Errorf("error playing week %d: %w", week, err)
		}

		standings, err := s.standingRepo.GetAll()
		if err != nil {
			return nil, err
		}

		predictions, err := s.predictionSvc.GetPredictions()
		if err != nil {
			return nil, err
		}

		results = append(results, model.WeekResult{
			Week:           week,
			Matches:        matches,
			Standings:      standings,
			Predictions:    predictions,
			LeagueFinished: week == 6,
		})
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
	fixtureSvc := NewFixtureService(s.matchRepo, s.teamRepo)
	return fixtureSvc.GenerateFixture()
}