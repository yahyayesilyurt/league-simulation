package service

import (
	"math"

	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"github.com/yahyayesilyurt/league-simulation/internal/repository"
)

type PredictionService interface {
	GetPredictions() ([]TeamPrediction, error)
}

type TeamPrediction struct {
	TeamID     uint    `json:"team_id"`
	TeamName   string  `json:"team_name"`
	Percentage float64 `json:"percentage"`
	CurrentPts int     `json:"current_points"`
	MaxPts     int     `json:"max_points"`
}

type predictionService struct {
	standingRepo repository.StandingRepository
	matchRepo    repository.MatchRepository
	teamRepo     repository.TeamRepository
}

func NewPredictionService(
	standingRepo repository.StandingRepository,
	matchRepo repository.MatchRepository,
	teamRepo repository.TeamRepository,
) PredictionService {
	return &predictionService{
		standingRepo: standingRepo,
		matchRepo:    matchRepo,
		teamRepo:     teamRepo,
	}
}

func (s *predictionService) GetPredictions() ([]TeamPrediction, error) {
	standings, err := s.standingRepo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(standings) == 0 {
		return []TeamPrediction{}, nil
	}

	currentWeek, err := s.getCurrentWeek()
	if err != nil {
		return nil, err
	}

	// Before week 4 → equal distribution
	if currentWeek < 4 {
		return s.equalDistribution(standings), nil
	}

	// Is the league over?
	if currentWeek >= 6 {
		return s.finishedLeague(standings), nil
	}

	// Weeks 4-5 → actual calculation
	return s.calculatePredictions(standings, currentWeek)
}

// equalDistribution — 25% equal distribution before week 4
func (s *predictionService) equalDistribution(standings []model.Standing) []TeamPrediction {
	result := make([]TeamPrediction, len(standings))
	for i, st := range standings {
		result[i] = TeamPrediction{
			TeamID:     st.TeamID,
			TeamName:   st.Team.Name,
			Percentage: 25.0,
			CurrentPts: st.Points,
			MaxPts:     st.Points + (6-st.Played)*3,
		}
	}
	return result
}

// finishedLeague — If the league is over, the winner is 100%
func (s *predictionService) finishedLeague(standings []model.Standing) []TeamPrediction {
	result := make([]TeamPrediction, len(standings))

	topPoints := standings[0].Points
	topGD     := standings[0].GoalDiff
	topGF     := standings[0].GoalsFor

	winners := 0
	for _, st := range standings {
		if st.Points == topPoints && st.GoalDiff == topGD && st.GoalsFor == topGF {
			winners++
		}
	}

	for i, st := range standings {
		pct := 0.0
		if st.Points == topPoints && st.GoalDiff == topGD && st.GoalsFor == topGF {
			pct = 100.0 / float64(winners) 
		}
		result[i] = TeamPrediction{
			TeamID:     st.TeamID,
			TeamName:   st.Team.Name,
			Percentage: math.Round(pct*100) / 100,
			CurrentPts: st.Points,
			MaxPts:     st.Points,
		}
	}
	return result
}

// calculatePredictions — Actual calculation in weeks 4-5.
func (s *predictionService) calculatePredictions(standings []model.Standing, currentWeek int) ([]TeamPrediction, error) {
	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}

	strengthMap := make(map[uint]int)
	for _, t := range teams {
		strengthMap[t.ID] = t.Strength
	}

	remainingMatches, err := s.matchRepo.GetUnplayed()
	if err != nil {
		return nil, err
	}

	remainingPtsMap := make(map[uint]float64)
	for _, match := range remainingMatches {
		homeStr := float64(strengthMap[match.HomeTeamID] + HOME_ADVANTAGE)
		awayStr := float64(strengthMap[match.AwayTeamID])
		total   := homeStr + awayStr

		homeWinProb := homeStr / total
		awayWinProb := awayStr / total
		drawProb    := 0.25 

		scale       := 1.0 - drawProb
		homeWinProb  = homeWinProb * scale
		awayWinProb  = awayWinProb * scale

		remainingPtsMap[match.HomeTeamID] += homeWinProb*3 + drawProb
		remainingPtsMap[match.AwayTeamID] += awayWinProb*3 + drawProb
	}

	scores := make(map[uint]float64)
	for _, st := range standings {
		currentPts  := float64(st.Points)
		expectedPts := remainingPtsMap[st.TeamID]
		gdBonus     := float64(st.GoalDiff) * 0.1 
		gfBonus     := float64(st.GoalsFor) * 0.05

		scores[st.TeamID] = currentPts + expectedPts + gdBonus + gfBonus
	}

	minScore := 0.0
	for _, score := range scores {
		if score < minScore {
			minScore = score
		}
	}
	if minScore < 0 {
		for id := range scores {
			scores[id] -= minScore
		}
	}

	totalScore := 0.0
	for _, score := range scores {
		totalScore += score
	}

	result := make([]TeamPrediction, len(standings))
	for i, st := range standings {
		pct := 0.0
		if totalScore > 0 {
			pct = (scores[st.TeamID] / totalScore) * 100
		}
		result[i] = TeamPrediction{
			TeamID:     st.TeamID,
			TeamName:   st.Team.Name,
			Percentage: math.Round(pct*10) / 10, 
			CurrentPts: st.Points,
			MaxPts:     st.Points + countRemainingMatches(remainingMatches, st.TeamID)*3,
		}
	}

	return result, nil
}

func (s *predictionService) getCurrentWeek() (int, error) {
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

func countRemainingMatches(matches []model.Match, teamID uint) int {
	count := 0
	for _, m := range matches {
		if m.HomeTeamID == teamID || m.AwayTeamID == teamID {
			count++
		}
	}
	return count
}