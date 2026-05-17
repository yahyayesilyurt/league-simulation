package service

import (
	"math"
	"math/rand"
	"time"
)

const (
	HOME_ADVANTAGE = 5   // Home team bonus strength
	BASE_GOALS     = 2.5 // Average expected goals
)

type SimulationEngine struct {
	rng *rand.Rand
}

func NewSimulationEngine() *SimulationEngine {
	return &SimulationEngine{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SimulateMatch - simulates a match between two teams
// homeStrength, awayStrength: Team power scores (1-100)
// Returns: homeGoals, awayGoals
func (e *SimulationEngine) SimulateMatch(homeStrength, awayStrength int) (int, int) {
	effectiveHome := homeStrength + HOME_ADVANTAGE

	totalStrength := float64(effectiveHome + awayStrength)
	homeRatio     := float64(effectiveHome) / totalStrength
	awayRatio     := float64(awayStrength) / totalStrength

	homeLambda := homeRatio * BASE_GOALS * 2
	awayLambda := awayRatio * BASE_GOALS * 2

	homeGoals := e.poissonRandom(homeLambda)
	awayGoals := e.poissonRandom(awayLambda)

	return homeGoals, awayGoals
}

// poissonRandom — It generates random numbers according to the Poisson distribution
func (e *SimulationEngine) poissonRandom(lambda float64) int {
	// Knuth algorithm
	L := math.Exp(-lambda)
	k := 0
	p := 1.0

	for p > L {
		k++
		p *= e.rng.Float64()
	}

	return k - 1
}