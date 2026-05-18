package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yahyayesilyurt/league-simulation/internal/model"
)

func TestGenerateRoundRobin_CorrectMatchCount(t *testing.T) {
	teams := []model.Team{
		{ID: 1, Name: "Man City",  Strength: 90},
		{ID: 2, Name: "Liverpool", Strength: 85},
		{ID: 3, Name: "Arsenal",   Strength: 80},
		{ID: 4, Name: "Chelsea",   Strength: 75},
	}

	matches := generateRoundRobin(teams)

	assert.Len(t, matches, 12, "Should generate 12 matches")
}

func TestGenerateRoundRobin_CorrectWeekCount(t *testing.T) {
	teams := []model.Team{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4},
	}

	matches := generateRoundRobin(teams)

	weeks := make(map[int]bool)
	for _, m := range matches {
		weeks[m.Week] = true
	}

	assert.Len(t, weeks, 6, "Should have 6 weeks")
}

func TestGenerateRoundRobin_TwoMatchesPerWeek(t *testing.T) {
	teams := []model.Team{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4},
	}

	matches := generateRoundRobin(teams)

	weekMatches := make(map[int]int)
	for _, m := range matches {
		weekMatches[m.Week]++
	}

	for week, count := range weekMatches {
		assert.Equal(t, 2, count, "Week %d should have 2 matches", week)
	}
}

func TestGenerateRoundRobin_NoTeamPlaysItself(t *testing.T) {
	teams := []model.Team{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4},
	}

	matches := generateRoundRobin(teams)

	for _, m := range matches {
		assert.NotEqual(t, m.HomeTeamID, m.AwayTeamID,
			"A team should not play against itself")
	}
}

func TestGenerateRoundRobin_EachTeamPlaysSixMatches(t *testing.T) {
	teams := []model.Team{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4},
	}

	matches := generateRoundRobin(teams)

	matchCount := make(map[uint]int)
	for _, m := range matches {
		matchCount[m.HomeTeamID]++
		matchCount[m.AwayTeamID]++
	}

	for teamID, count := range matchCount {
		assert.Equal(t, 6, count,
			"Team %d should play 6 matches", teamID)
	}
}

func TestGenerateRoundRobin_SecondLegReversesHomeAway(t *testing.T) {
	teams := []model.Team{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4},
	}

	matches := generateRoundRobin(teams)

	firstLeg := make(map[[2]uint]bool)
	for _, m := range matches {
		if m.Week <= 3 {
			firstLeg[[2]uint{m.HomeTeamID, m.AwayTeamID}] = true
		}
	}

	for _, m := range matches {
		if m.Week > 3 {
			reversed := [2]uint{m.AwayTeamID, m.HomeTeamID}
			assert.True(t, firstLeg[reversed],
				"Second leg should reverse home/away teams")
		}
	}
}