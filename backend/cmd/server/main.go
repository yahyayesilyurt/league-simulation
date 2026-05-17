package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yahyayesilyurt/league-simulation/config"
	"github.com/yahyayesilyurt/league-simulation/internal/handler"
	"github.com/yahyayesilyurt/league-simulation/internal/repository"
	"github.com/yahyayesilyurt/league-simulation/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config.ConnectDatabase()
	config.ConnectRedis()
	config.SeedDatabase(config.DB)

	matchRepo  := repository.NewMatchRepository(config.DB)
	teamRepo   := repository.NewTeamRepository(config.DB)
	fixtureSvc := service.NewFixtureService(matchRepo, teamRepo)

	if err := fixtureSvc.GenerateFixture(); err != nil {
		log.Printf("Fixture info: %s", err.Error())
	} else {
		log.Println("Fixture generated successfully")
	}

	r := handler.SetupRouter(config.DB, config.RedisClient)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}