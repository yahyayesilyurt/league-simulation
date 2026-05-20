package repository

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"gorm.io/gorm"
)

func setupStandingTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Team{}, &model.Standing{})
	require.NoError(t, err)

	return db
}

func TestStandingRepository_GetAll_Ordering(t *testing.T) {
	db := setupStandingTestDB(t)
	repo := NewStandingRepository(db)

	team1 := model.Team{Name: "Team A", Strength: 80}
	team2 := model.Team{Name: "Team B", Strength: 80}
	team3 := model.Team{Name: "Team C", Strength: 80}
	db.Create(&team1)
	db.Create(&team2)
	db.Create(&team3)

	repo.Update(&model.Standing{TeamID: team3.ID, Points: 6, GoalDiff: 1, GoalsFor: 1})
	repo.Update(&model.Standing{TeamID: team2.ID, Points: 9, GoalDiff: 2, GoalsFor: 3})
	repo.Update(&model.Standing{TeamID: team1.ID, Points: 9, GoalDiff: 5, GoalsFor: 6})

	standings, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, standings, 3)

	assert.Equal(t, "Team A", standings[0].Team.Name)
	assert.Equal(t, 9, standings[0].Points)

	assert.Equal(t, "Team B", standings[1].Team.Name)

	assert.Equal(t, "Team C", standings[2].Team.Name)
	assert.Equal(t, 6, standings[2].Points)
}

func TestStandingRepository_GetByTeamID(t *testing.T) {
	db := setupStandingTestDB(t)
	repo := NewStandingRepository(db)

	team := model.Team{Name: "Team Name", Strength: 75}
	db.Create(&team)

	standing := &model.Standing{
		TeamID: team.ID,
		Points: 12,
		Won:    4,
	}
	repo.Update(standing)

	fetched, err := repo.GetByTeamID(team.ID)

	assert.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, 12, fetched.Points)
	
	assert.Equal(t, "Team Name", fetched.Team.Name)
}

func TestStandingRepository_Update(t *testing.T) {
	db := setupStandingTestDB(t)
	repo := NewStandingRepository(db)

	team := model.Team{Name: "Team Test", Strength: 85}
	db.Create(&team)

	standing := model.Standing{TeamID: team.ID, Points: 0, Played: 0}
	repo.Update(&standing) 

	standing.Points = 3
	standing.Played = 1
	standing.Won = 1
	err := repo.Update(&standing)

	assert.NoError(t, err)

	var updatedStanding model.Standing
	db.First(&updatedStanding, "team_id = ?", team.ID)

	assert.Equal(t, 3, updatedStanding.Points)
	assert.Equal(t, 1, updatedStanding.Played)
	assert.Equal(t, 1, updatedStanding.Won)
}

func TestStandingRepository_ResetAll(t *testing.T) {
	db := setupStandingTestDB(t)
	repo := NewStandingRepository(db)

	team1 := model.Team{Name: "Team 1", Strength: 80}
	team2 := model.Team{Name: "Team 2", Strength: 80}
	db.Create(&team1)
	db.Create(&team2)

	db.Create(&model.Standing{TeamID: team1.ID, Played: 5, Won: 3, Points: 9, GoalsFor: 10})
	db.Create(&model.Standing{TeamID: team2.ID, Played: 5, Won: 1, Points: 3, GoalsFor: 4})

	err := repo.ResetAll()
	assert.NoError(t, err)

	standings, _ := repo.GetAll()
	for _, st := range standings {
		assert.Equal(t, 0, st.Played)
		assert.Equal(t, 0, st.Won)
		assert.Equal(t, 0, st.Drawn)
		assert.Equal(t, 0, st.Lost)
		assert.Equal(t, 0, st.GoalsFor)
		assert.Equal(t, 0, st.GoalsAgainst)
		assert.Equal(t, 0, st.GoalDiff)
		assert.Equal(t, 0, st.Points)
	}
}