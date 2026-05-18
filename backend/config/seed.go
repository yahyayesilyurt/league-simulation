package config

import (
	"github.com/rs/zerolog/log"
	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) {
	seedTeams(db)
	seedStandings(db)
}

func seedTeams(db *gorm.DB) {
	teams := []model.Team{
		{Name: "Manchester City", Strength: 90},
		{Name: "Liverpool",       Strength: 85},
		{Name: "Arsenal",         Strength: 80},
		{Name: "Chelsea",         Strength: 75},
	}

	for _, team := range teams {
		var existing model.Team
		result := db.Where("name = ?", team.Name).First(&existing)
		if result.Error != nil {
			if err := db.Create(&team).Error; err != nil {
				log.Error().Err(err).Str("team", team.Name).Msg("Failed to seed team")
			} else {
				log.Info().Str("team", team.Name).Int("strength", team.Strength).Msg("Seeded team")
			}
		} else {
			log.Debug().Str("team", team.Name).Msg("Team already exists, skipping")
		}
	}
}

func seedStandings(db *gorm.DB) {
	var teams []model.Team
	if err := db.Find(&teams).Error; err != nil {
		log.Error().Err(err).Msg("Failed to fetch teams for standings seed")
		return
	}

	for _, team := range teams {
		var existing model.Standing
		result := db.Where("team_id = ?", team.ID).First(&existing)
		if result.Error != nil {
			standing := model.Standing{TeamID: team.ID}
			if err := db.Create(&standing).Error; err != nil {
				log.Error().Err(err).Str("team", team.Name).Msg("Failed to seed standing")
			} else {
				log.Debug().Str("team", team.Name).Msg("Seeded standing")
			}
		}
	}
}