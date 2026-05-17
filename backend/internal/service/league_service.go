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
	GetStatus() (*model.LeagueStatus, error)
	NextWeek() (*model.WeekResult, error)
	PlayAll() (*model.PlayAllResult, error)
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

func (s *leagueService) PlayAll() (*model.PlayAllResult, error) {
	currentWeek, err := s.GetCurrentWeek()
	if err != nil {
		return nil, err
	}
	if currentWeek >= 6 {
		return nil, fmt.Errorf("league is already finished")
	}

	var weekResults []model.WeekResult

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

		weekResults = append(weekResults, model.WeekResult{
			Week:           week,
			Matches:        matches,
			Standings:      standings,
			Predictions:    predictions,
			LeagueFinished: week == 6,
		})
	}

	finalStandings, err := s.standingRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return &model.PlayAllResult{
		TotalWeeksPlayed: len(weekResults),
		Weeks:            weekResults,
		Summary:          s.buildSummary(finalStandings),
	}, nil
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

func (s *leagueService) GetStatus() (*model.LeagueStatus, error) {
	matches, err := s.matchRepo.GetAll()
	if err != nil {
		return nil, err
	}

	played := 0
	for _, m := range matches {
		if m.Played {
			played++
		}
	}

	total       := len(matches) 
	left        := total - played
	currentWeek, _ := s.GetCurrentWeek()

	status := "not_started"
	if played > 0 && left > 0 {
		status = "in_progress"
	} else if left == 0 && total > 0 {
		status = "finished"
	}

	return &model.LeagueStatus{
		CurrentWeek:    currentWeek,
		TotalWeeks:     6,
		LeagueFinished: left == 0 && total > 0,
		MatchesPlayed:  played,
		MatchesLeft:    left,
		Status:         status,
	}, nil
}

func (s *leagueService) buildSummary(standings []model.Standing) *model.LeagueSummary {
	if len(standings) == 0 {
		return nil
	}

	champion := standings[0]

	topScorer    := standings[0]
	bestDefense  := standings[0]
	totalGoals   := 0

	for _, st := range standings {
		totalGoals += st.GoalsFor
		if st.GoalsFor > topScorer.GoalsFor {
			topScorer = st
		}
		if st.GoalsAgainst < bestDefense.GoalsAgainst {
			bestDefense = st
		}
	}

	return &model.LeagueSummary{
		Champion:       &champion,
		FinalStandings: standings,
		TopScorer:      topScorer.Team.Name,
		BestDefense:    bestDefense.Team.Name,
		TotalGoals:     totalGoals,
		TotalMatches:   12,
	}
}