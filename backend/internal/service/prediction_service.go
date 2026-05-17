package service

import (
	"github.com/yahyayesilyurt/league-simulation/internal/repository"
)

type PredictionService interface {
	GetPredictions() (map[string]float64, error)
}

type predictionService struct {
	standingRepo repository.StandingRepository
	matchRepo    repository.MatchRepository
}

func NewPredictionService(
	standingRepo repository.StandingRepository,
	matchRepo repository.MatchRepository,
) PredictionService {
	return &predictionService{
		standingRepo: standingRepo,
		matchRepo:    matchRepo,
	}
}

func (s *predictionService) GetPredictions() (map[string]float64, error) {
	standings, err := s.standingRepo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string]float64)
	if len(standings) == 0 {
		return result, nil
	}

	equal := 100.0 / float64(len(standings))
	for _, s := range standings {
		result[s.Team.Name] = equal
	}
	return result, nil
}