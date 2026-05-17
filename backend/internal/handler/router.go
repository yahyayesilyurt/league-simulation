package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yahyayesilyurt/league-simulation/internal/repository"
	"github.com/yahyayesilyurt/league-simulation/internal/service"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Create repositories
	teamRepo     := repository.NewTeamRepository(db)
	matchRepo    := repository.NewMatchRepository(db)
	standingRepo := repository.NewStandingRepository(db)

	// Create services
	leagueSvc     := service.NewLeagueService(matchRepo, standingRepo, teamRepo)
	predictionSvc := service.NewPredictionService(standingRepo, matchRepo)

	// Create handlers
	leagueHandler := NewLeagueHandler(leagueSvc, predictionSvc)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "League Simulation is running",
		})
	})

	// League routes
	league := r.Group("/league")
	{
		league.GET("/table",       leagueHandler.GetStandings)
		league.GET("/fixtures",    leagueHandler.GetFixtures)
		league.GET("/week/:weekNo", leagueHandler.GetWeek)
		league.GET("/predictions", leagueHandler.GetPredictions)
		league.POST("/next-week",  leagueHandler.NextWeek)
		league.POST("/play-all",   leagueHandler.PlayAll)
		league.POST("/reset",      leagueHandler.Reset)
	}

	return r
}