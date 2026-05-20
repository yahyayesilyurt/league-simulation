package config

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"gorm.io/gorm"
)

func setupSeedDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Team{}, &model.Standing{})
	require.NoError(t, err)

	return db
}

func TestSeedDatabase_InitialRun(t *testing.T) {
	db := setupSeedDB(t)

	SeedDatabase(db)

	var teams []model.Team
	db.Find(&teams)
	assert.Len(t, teams, 4)

	var city model.Team
	db.Where("name = ?", "Manchester City").First(&city)
	assert.Equal(t, 90, city.Strength)

	var standings []model.Standing
	db.Find(&standings)
	assert.Len(t, standings, 4)
	
	assert.Equal(t, city.ID, standings[0].TeamID)
}

func TestSeedDatabase_Idempotency(t *testing.T) {
	db := setupSeedDB(t)

	SeedDatabase(db)
	SeedDatabase(db) 

	var teamCount int64
	db.Model(&model.Team{}).Count(&teamCount)
	assert.Equal(t, int64(4), teamCount)

	var standingCount int64
	db.Model(&model.Standing{}).Count(&standingCount)
	assert.Equal(t, int64(4), standingCount)
}