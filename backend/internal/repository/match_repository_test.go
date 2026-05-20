package repository

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)

    err = db.AutoMigrate(&model.Team{}, &model.Match{})
    require.NoError(t, err)

    return db
}

func TestMatchRepository_CreateAndGetByID(t *testing.T) {
    db := setupTestDB(t)
    repo := NewMatchRepository(db)

    team1 := model.Team{Name: "Home Team", Strength: 80}
    team2 := model.Team{Name: "Away Team", Strength: 80}
    
    db.Create(&team1)
    db.Create(&team2)

    match := &model.Match{
        Week:       1,
        HomeTeamID: team1.ID, 
        AwayTeamID: team2.ID,
        Played:     false,
    }

    err := repo.Create(match)
    assert.NoError(t, err)
    assert.NotZero(t, match.ID) 

    fetchedMatch, err := repo.GetByID(match.ID)
    assert.NoError(t, err)
    assert.Equal(t, match.Week, fetchedMatch.Week)
    
    assert.Equal(t, "Home Team", fetchedMatch.HomeTeam.Name)
    assert.Equal(t, "Away Team", fetchedMatch.AwayTeam.Name)
}

func TestMatchRepository_GetByWeek(t *testing.T) {
    db := setupTestDB(t)
    repo := NewMatchRepository(db)

    repo.Create(&model.Match{Week: 1})
    repo.Create(&model.Match{Week: 1})
    repo.Create(&model.Match{Week: 2})

    matches, err := repo.GetByWeek(1)
    
    assert.NoError(t, err)
    assert.Len(t, matches, 2) 
    assert.Equal(t, 1, matches[0].Week)
}

func TestMatchRepository_GetUnplayed(t *testing.T) {
    db := setupTestDB(t)
    repo := NewMatchRepository(db)

    repo.Create(&model.Match{Played: false})
    repo.Create(&model.Match{Played: false})
    repo.Create(&model.Match{Played: true})

    unplayed, err := repo.GetUnplayed()

    assert.NoError(t, err)
    assert.Len(t, unplayed, 2)
    assert.False(t, unplayed[0].Played)
    assert.False(t, unplayed[1].Played)
}

func TestMatchRepository_GetAll(t *testing.T) {
    db := setupTestDB(t)
    repo := NewMatchRepository(db)

    repo.Create(&model.Match{Week: 2})
    repo.Create(&model.Match{Week: 1})

    matches, err := repo.GetAll()

    assert.NoError(t, err)
    assert.Len(t, matches, 2)
    
    assert.Equal(t, 1, matches[0].Week)
    assert.Equal(t, 2, matches[1].Week)
}

func TestMatchRepository_Update(t *testing.T) {
    db := setupTestDB(t)
    repo := NewMatchRepository(db)

    match := &model.Match{Played: false}
    repo.Create(match)

    match.Played = true
    err := repo.Update(match)
    assert.NoError(t, err)

    fetched, _ := repo.GetByID(match.ID)
    assert.True(t, fetched.Played)
}

func TestMatchRepository_DeleteAll(t *testing.T) {
    db := setupTestDB(t)
    repo := NewMatchRepository(db)

    repo.Create(&model.Match{Week: 1})
    repo.Create(&model.Match{Week: 2})

    err := repo.DeleteAll()
    assert.NoError(t, err)

    matches, err := repo.GetAll()
    assert.NoError(t, err)
    assert.Len(t, matches, 0)
}