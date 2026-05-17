package service

import (
	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"github.com/yahyayesilyurt/league-simulation/internal/repository"
)

type StandingService interface {
	GetStandings() ([]model.Standing, error)
	RecalculateAll() error
}

type standingService struct {
	standingRepo repository.StandingRepository
	matchRepo    repository.MatchRepository
	teamRepo     repository.TeamRepository
}

func NewStandingService(
	standingRepo repository.StandingRepository,
	matchRepo repository.MatchRepository,
	teamRepo repository.TeamRepository,
) StandingService {
	return &standingService{
		standingRepo: standingRepo,
		matchRepo:    matchRepo,
		teamRepo:     teamRepo,
	}
}

// GetStandings — It returns the order of the League standings
// Ranking: Points → Goal difference → Goals scored
func (s *standingService) GetStandings() ([]model.Standing, error) {
	return s.standingRepo.GetAll()
}

// RecalculateAll — Calculations from scratch based on all matches played
// Called when match result is finalized
func (s *standingService) RecalculateAll() error {
	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return err
	}

	if err := s.standingRepo.ResetAll(); err != nil {
		return err
	}

	matches, err := s.matchRepo.GetAll()
	if err != nil {
		return err
	}

	standingMap := make(map[uint]*model.Standing)
	for _, team := range teams {
		t := team
		standing, err := s.standingRepo.GetByTeamID(t.ID)
		if err != nil {
			return err
		}
		standingMap[t.ID] = standing
	}

	for _, match := range matches {
		if !match.Played || match.HomeGoals == nil || match.AwayGoals == nil {
			continue
		}

		homeGoals := *match.HomeGoals
		awayGoals := *match.AwayGoals

		home := standingMap[match.HomeTeamID]
		away := standingMap[match.AwayTeamID]

		home.Played++
		away.Played++

		home.GoalsFor     += homeGoals
		home.GoalsAgainst += awayGoals
		away.GoalsFor     += awayGoals
		away.GoalsAgainst += homeGoals

		home.GoalDiff = home.GoalsFor - home.GoalsAgainst
		away.GoalDiff = away.GoalsFor - away.GoalsAgainst

		switch {
		case homeGoals > awayGoals:
			home.Won++
			home.Points += 3
			away.Lost++
		case homeGoals < awayGoals:
			away.Won++
			away.Points += 3
			home.Lost++
		default:
			home.Drawn++
			home.Points++
			away.Drawn++
			away.Points++
		}
	}

	for _, standing := range standingMap {
		if err := s.standingRepo.Update(standing); err != nil {
			return err
		}
	}

	return nil
}