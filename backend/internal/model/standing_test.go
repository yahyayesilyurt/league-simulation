package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStanding_TableName(t *testing.T) {
	standing := Standing{}
	
	assert.Equal(t, "standings", standing.TableName())
}

func TestStanding_JSONSerialization(t *testing.T) {
	standing := Standing{
		ID:           1,
		TeamID:       10,
		Played:       5,
		Won:          3,
		Drawn:        1,
		Lost:         1,
		GoalsFor:     10,
		GoalsAgainst: 4,
		GoalDiff:     6,
		Points:       10,
	}

	jsonData, err := json.Marshal(standing)
	assert.NoError(t, err)

	jsonString := string(jsonData)

	assert.Contains(t, jsonString, `"id":1`)
	assert.Contains(t, jsonString, `"team_id":10`)
	assert.Contains(t, jsonString, `"played":5`)
	assert.Contains(t, jsonString, `"won":3`)
	assert.Contains(t, jsonString, `"drawn":1`)
	assert.Contains(t, jsonString, `"lost":1`)
	assert.Contains(t, jsonString, `"goals_for":10`)
	assert.Contains(t, jsonString, `"goals_against":4`)
	assert.Contains(t, jsonString, `"goal_diff":6`)
	assert.Contains(t, jsonString, `"points":10`)
}