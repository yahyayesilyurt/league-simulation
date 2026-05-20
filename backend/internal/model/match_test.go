package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch_TableName(t *testing.T) {
	match := Match{}
	
	assert.Equal(t, "matches", match.TableName())
}

func TestMatch_JSONSerialization(t *testing.T) {
	unplayedMatch := Match{
		ID:         1,
		HomeTeamID: 10,
		AwayTeamID: 20,
		Played:     false,
	}

	unplayedJSON, err := json.Marshal(unplayedMatch)
	assert.NoError(t, err)

	assert.Contains(t, string(unplayedJSON), `"home_goals":null`)
	assert.Contains(t, string(unplayedJSON), `"away_goals":null`)
	assert.Contains(t, string(unplayedJSON), `"played":false`)

	homeGoals := 2
	awayGoals := 1
	playedMatch := Match{
		ID:         2,
		HomeTeamID: 10,
		AwayTeamID: 20,
		HomeGoals:  &homeGoals,
		AwayGoals:  &awayGoals,
		Played:     true,
	}

	playedJSON, err := json.Marshal(playedMatch)
	assert.NoError(t, err)

	assert.Contains(t, string(playedJSON), `"home_goals":2`)
	assert.Contains(t, string(playedJSON), `"away_goals":1`)
	assert.Contains(t, string(playedJSON), `"played":true`)
}