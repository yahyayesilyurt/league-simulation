package repository

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"gorm.io/gorm"
)

func setupTeamTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Team{})
	require.NoError(t, err)

	return db
}

func TestTeamRepository_CreateAndGetByID(t *testing.T) {
	db := setupTeamTestDB(t)
	repo := NewTeamRepository(db)

	team := &model.Team{
		Name:     "Team A",
		Strength: 85, 
	}

	err := repo.Create(team)
	assert.NoError(t, err)
	assert.NotZero(t, team.ID)

	fetchedTeam, err := repo.GetByID(team.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Team A", fetchedTeam.Name)
	assert.Equal(t, 85, fetchedTeam.Strength)
}

func TestTeamRepository_GetByID_NotFound(t *testing.T) {
	db := setupTeamTestDB(t)
	repo := NewTeamRepository(db)

	fetchedTeam, err := repo.GetByID(999)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Equal(t, uint(0), fetchedTeam.ID)
}

func TestTeamRepository_GetAll(t *testing.T) {
	db := setupTeamTestDB(t)
	repo := NewTeamRepository(db)

	team1 := &model.Team{Name: "Team A", Strength: 70}
	team2 := &model.Team{Name: "Team B", Strength: 75}
	team3 := &model.Team{Name: "Team C", Strength: 80}

	repo.Create(team1)
	repo.Create(team2)
	repo.Create(team3)

	teams, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, teams, 3)

	assert.Equal(t, "Team A", teams[0].Name)
	assert.Equal(t, "Team B", teams[1].Name)
	assert.Equal(t, "Team C", teams[2].Name)
}

func TestTeamRepository_Update(t *testing.T) {
	db := setupTeamTestDB(t)
	repo := NewTeamRepository(db)

	team := &model.Team{
		Name:     "Old Name",
		Strength: 60,
	}
	repo.Create(team)

	team.Name = "New Name"
	team.Strength = 90
	err := repo.Update(team)

	assert.NoError(t, err)

	fetchedTeam, _ := repo.GetByID(team.ID)
	assert.Equal(t, "New Name", fetchedTeam.Name)
	assert.Equal(t, 90, fetchedTeam.Strength)
}