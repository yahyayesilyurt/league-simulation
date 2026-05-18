package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimulateMatch_ReturnsNonNegativeGoals(t *testing.T) {
	engine := NewSimulationEngine()

	homeGoals, awayGoals := engine.SimulateMatch(80, 75)

	assert.GreaterOrEqual(t, homeGoals, 0, "Home goals must be >= 0")
	assert.GreaterOrEqual(t, awayGoals, 0, "Away goals must be >= 0")
}

func TestSimulateMatch_StrongerTeamWinsMoreOften(t *testing.T) {
	engine := NewSimulationEngine()

	strongWins := 0
	total      := 1000

	for i := 0; i < total; i++ {
		homeGoals, awayGoals := engine.SimulateMatch(90, 50)
		if homeGoals > awayGoals {
			strongWins++
		}
	}

	winRate := float64(strongWins) / float64(total)
	assert.Greater(t, winRate, 0.45, "Stronger team should win more often")
}

func TestSimulateMatch_HomeAdvantage(t *testing.T) {
	engine := NewSimulationEngine()

	homeWins := 0
	awayWins := 0
	total    := 1000

	for i := 0; i < total; i++ {
		homeGoals, awayGoals := engine.SimulateMatch(75, 75)
		if homeGoals > awayGoals {
			homeWins++
		} else if awayGoals > homeGoals {
			awayWins++
		}
	}

	assert.Greater(t, homeWins, awayWins, "Home team should win more due to home advantage")
}

func TestPoissonRandom_ReturnsNonNegative(t *testing.T) {
	engine := NewSimulationEngine()

	for i := 0; i < 100; i++ {
		result := engine.poissonRandom(2.5)
		assert.GreaterOrEqual(t, result, 0, "Poisson result must be >= 0")
	}
}

func TestPoissonRandom_AverageCloseToLambda(t *testing.T) {
	engine := NewSimulationEngine()
	lambda := 2.5
	total  := 10000
	sum    := 0

	for i := 0; i < total; i++ {
		sum += engine.poissonRandom(lambda)
	}

	avg := float64(sum) / float64(total)

	assert.InDelta(t, lambda, avg, 0.3, "Average should be close to lambda")
}