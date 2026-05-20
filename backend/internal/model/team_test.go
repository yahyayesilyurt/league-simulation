package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeam_TableName(t *testing.T) {
	team := Team{}

	assert.Equal(t, "teams", team.TableName())
}

func TestTeam_JSONSerialization(t *testing.T) {
	team := Team{
		ID:       1,
		Name:     "Arsenal",
		Strength: 85,
	}

	jsonData, err := json.Marshal(team)
	assert.NoError(t, err)

	jsonString := string(jsonData)

	assert.Contains(t, jsonString, `"id":1`)
	assert.Contains(t, jsonString, `"name":"Arsenal"`)
	assert.Contains(t, jsonString, `"strength":85`)
}