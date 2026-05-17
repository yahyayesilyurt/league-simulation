package service

import (
	"fmt"

	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"github.com/yahyayesilyurt/league-simulation/internal/repository"
)

type MatchService interface {
	PlayMatch(matchID uint) (*model.Match, error)
	PlayWeek(week int) ([]model.Match, error)
	UpdateMatchResult(matchID uint, homeGoals, awayGoals int) (*model.Match, error)
}

type matchService struct {
	matchRepo    repository.MatchRepository
	standingRepo repository.StandingRepository
	engine       *SimulationEngine
}

func NewMatchService(
	matchRepo repository.MatchRepository,
	standingRepo repository.StandingRepository,
) MatchService {
	return &matchService{
		matchRepo:    matchRepo,
		standingRepo: standingRepo,
		engine:       NewSimulationEngine(),
	}
}

// PlayMatch — It simulates a single match and records the result
func (s *matchService) PlayMatch(matchID uint) (*model.Match, error) {
	match, err := s.matchRepo.GetByID(matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found: %w", err)
	}
	if match.Played {
		return nil, fmt.Errorf("match already played")
	}

	homeGoals, awayGoals := s.engine.SimulateMatch(
		match.HomeTeam.Strength,
		match.AwayTeam.Strength,
	)

	match.HomeGoals = &homeGoals
	match.AwayGoals = &awayGoals
	match.Played = true

	if err := s.matchRepo.Update(match); err != nil {
		return nil, fmt.Errorf("failed to update match: %w", err)
	}

	if err := s.updateStandings(match); err != nil {
		return nil, fmt.Errorf("failed to update standings: %w", err)
	}

	return match, nil
}

// PlayWeek — It plays all the matches of a week
func (s *matchService) PlayWeek(week int) ([]model.Match, error) {
	matches, err := s.matchRepo.GetByWeek(week)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("no matches found for week %d", week)
	}

	var results []model.Match
	for _, match := range matches {
		if match.Played {
			continue
		}
		result, err := s.PlayMatch(match.ID)
		if err != nil {
			return nil, err
		}
		results = append(results, *result)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("week %d already played", week)
	}

	return results, nil
}

// UpdateMatchResult — Manually updates the match result
func (s *matchService) UpdateMatchResult(matchID uint, homeGoals, awayGoals int) (*model.Match, error) {
	match, err := s.matchRepo.GetByID(matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found: %w", err)
	}
	if !match.Played {
		return nil, fmt.Errorf("match has not been played yet")
	}

	if err := s.revertStandings(match); err != nil {
		return nil, fmt.Errorf("failed to revert standings: %w", err)
	}

	match.HomeGoals = &homeGoals
	match.AwayGoals = &awayGoals
	if err := s.matchRepo.Update(match); err != nil {
		return nil, fmt.Errorf("failed to update match: %w", err)
	}

	if err := s.updateStandings(match); err != nil {
		return nil, fmt.Errorf("failed to update standings: %w", err)
	}

	return match, nil
}

// updateStandings — Updates the points table based on match results
func (s *matchService) updateStandings(match *model.Match) error {
	homeGoals := *match.HomeGoals
	awayGoals := *match.AwayGoals

	homeStanding, err := s.standingRepo.GetByTeamID(match.HomeTeamID)
	if err != nil {
		return err
	}
	awayStanding, err := s.standingRepo.GetByTeamID(match.AwayTeamID)
	if err != nil {
		return err
	}

	homeStanding.Played++
	awayStanding.Played++

	homeStanding.GoalsFor     += homeGoals
	homeStanding.GoalsAgainst += awayGoals
	awayStanding.GoalsFor     += awayGoals
	awayStanding.GoalsAgainst += homeGoals

	homeStanding.GoalDiff = homeStanding.GoalsFor - homeStanding.GoalsAgainst
	awayStanding.GoalDiff = awayStanding.GoalsFor - awayStanding.GoalsAgainst

	switch {
	case homeGoals > awayGoals: // The home team won
		homeStanding.Won++
		homeStanding.Points += 3
		awayStanding.Lost++
	case homeGoals < awayGoals: // They won away
		awayStanding.Won++
		awayStanding.Points += 3
		homeStanding.Lost++
	default: // drew
		homeStanding.Drawn++
		homeStanding.Points++
		awayStanding.Drawn++
		awayStanding.Points++
	}

	if err := s.standingRepo.Update(homeStanding); err != nil {
		return err
	}
	return s.standingRepo.Update(awayStanding)
}

// revertStandings — retrieves the old match result from the standings
func (s *matchService) revertStandings(match *model.Match) error {
	homeGoals := *match.HomeGoals
	awayGoals := *match.AwayGoals

	homeStanding, err := s.standingRepo.GetByTeamID(match.HomeTeamID)
	if err != nil {
		return err
	}
	awayStanding, err := s.standingRepo.GetByTeamID(match.AwayTeamID)
	if err != nil {
		return err
	}

	homeStanding.Played--
	awayStanding.Played--
	homeStanding.GoalsFor     -= homeGoals
	homeStanding.GoalsAgainst -= awayGoals
	awayStanding.GoalsFor     -= awayGoals
	awayStanding.GoalsAgainst -= homeGoals
	homeStanding.GoalDiff      = homeStanding.GoalsFor - homeStanding.GoalsAgainst
	awayStanding.GoalDiff      = awayStanding.GoalsFor - awayStanding.GoalsAgainst

	switch {
	case homeGoals > awayGoals:
		homeStanding.Won--
		homeStanding.Points -= 3
		awayStanding.Lost--
	case homeGoals < awayGoals:
		awayStanding.Won--
		awayStanding.Points -= 3
		homeStanding.Lost--
	default:
		homeStanding.Drawn--
		homeStanding.Points--
		awayStanding.Drawn--
		awayStanding.Points--
	}

	if err := s.standingRepo.Update(homeStanding); err != nil {
		return err
	}
	return s.standingRepo.Update(awayStanding)
}