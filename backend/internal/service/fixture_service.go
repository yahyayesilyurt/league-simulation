package service

import (
	"fmt"

	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"github.com/yahyayesilyurt/league-simulation/internal/repository"
)

type FixtureService interface {
	GenerateFixture() error
	IsFixtureGenerated() (bool, error)
}

type fixtureService struct {
	matchRepo repository.MatchRepository
	teamRepo  repository.TeamRepository
}

func NewFixtureService(
	matchRepo repository.MatchRepository,
	teamRepo repository.TeamRepository,
) FixtureService {
	return &fixtureService{
		matchRepo: matchRepo,
		teamRepo:  teamRepo,
	}
}

func (s *fixtureService) IsFixtureGenerated() (bool, error) {
	matches, err := s.matchRepo.GetAll()
	if err != nil {
		return false, err
	}
	return len(matches) > 0, nil
}

func (s *fixtureService) GenerateFixture() error {
	generated, err := s.IsFixtureGenerated()
	if err != nil {
		return err
	}
	if generated {
		return fmt.Errorf("fixture already generated")
	}

	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to fetch teams: %w", err)
	}
	if len(teams) != 4 {
		return fmt.Errorf("expected 4 teams, got %d", len(teams))
	}

	matches := generateRoundRobin(teams)

	for _, match := range matches {
		m := match 
		if err := s.matchRepo.Create(&m); err != nil {
			return fmt.Errorf("failed to save match: %w", err)
		}
	}

	return nil
}

func generateRoundRobin(teams []model.Team) []model.Match {
	n := len(teams) 
	var matches []model.Match

	firstLeg := [][]int{
		{0, 1, 2, 3}, // Week 1: teams[0]-teams[1], teams[2]-teams[3]
		{0, 2, 3, 1}, // Week 2: teams[0]-teams[2], teams[3]-teams[1]
		{0, 3, 1, 2}, // Week 3: teams[0]-teams[3], teams[1]-teams[2]
	}

	for week, pairs := range firstLeg {
		matches = append(matches, model.Match{
			Week:       week + 1,
			HomeTeamID: teams[pairs[0]].ID,
			AwayTeamID: teams[pairs[1]].ID,
			Played:     false,
		})
		matches = append(matches, model.Match{
			Week:       week + 1,
			HomeTeamID: teams[pairs[2]].ID,
			AwayTeamID: teams[pairs[3]].ID,
			Played:     false,
		})
	}

	for i, match := range matches[:n/2*3] {
		week := (i/2) + 4 // Week 4, 5, 6
		matches = append(matches, model.Match{
			Week:       week,
			HomeTeamID: match.AwayTeamID, 
			AwayTeamID: match.HomeTeamID,
			Played:     false,
		})
	}

	return matches
}