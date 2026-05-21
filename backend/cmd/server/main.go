// @title           League Simulation API
// @version         1.0
// @description     A four-team football league simulation API. Match simulation, league table, and championship prediction based on Premier League rules.
// @termsOfService  http://swagger.io/terms/

// @license.name  MIT

// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token. Format: "Bearer <token>"

package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/yahyayesilyurt/league-simulation/config"
	_ "github.com/yahyayesilyurt/league-simulation/docs"
	"github.com/yahyayesilyurt/league-simulation/internal/handler"
	"github.com/yahyayesilyurt/league-simulation/internal/repository"
	"github.com/yahyayesilyurt/league-simulation/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using environment variables")
	}

	config.InitLogger()
	config.ConnectDatabase()
	config.ConnectRedis()
	config.SeedDatabase(config.DB)

	matchRepo  := repository.NewMatchRepository(config.DB)
	teamRepo   := repository.NewTeamRepository(config.DB)
	fixtureSvc := service.NewFixtureService(matchRepo, teamRepo)

	if err := fixtureSvc.GenerateFixture(); err != nil {
		log.Info().Str("info", err.Error()).Msg("Fixture status")
	} else {
		log.Info().Msg("Fixture generated successfully")
	}

	r := handler.SetupRouter(config.DB, config.RedisClient)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Str("port", port).Msg("Server starting")
	r.Run(":" + port)
}