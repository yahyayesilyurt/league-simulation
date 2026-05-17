package config

import (
	"log"

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
		{Name: "Liverpool", Strength: 85},
		{Name: "Arsenal", Strength: 80},
		{Name: "Chelsea", Strength: 75},
	}

	for _, team := range teams {
		var existing model.Team
		result := db.Where("name = ?", team.Name).First(&existing)

		if result.Error != nil {
			if err := db.Create(&team).Error; err != nil {
				log.Printf("Failed to seed team %s: %v", team.Name, err)
			} else {
				log.Printf("Seeded team: %s (strength: %d)", team.Name, team.Strength)
			}
		} else {
			log.Printf("Team already exists, skipping: %s", team.Name)
		}
	}
}

func seedStandings(db *gorm.DB) {
	var teams []model.Team
	if err := db.Find(&teams).Error; err != nil {
		log.Printf("Failed to fetch teams for standings seed: %v", err)
		return
	}

	for _, team := range teams {
		var existing model.Standing
		result := db.Where("team_id = ?", team.ID).First(&existing)

		if result.Error != nil {
			standing := model.Standing{TeamID: team.ID}
			if err := db.Create(&standing).Error; err != nil {
				log.Printf("Failed to seed standing for %s: %v", team.Name, err)
			} else {
				log.Printf("Seeded standing for: %s", team.Name)
			}
		} else {
			log.Printf("Standing already exists, skipping: %s", team.Name)
		}
	}
}