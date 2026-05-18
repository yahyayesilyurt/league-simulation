package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yahyayesilyurt/league-simulation/internal/cache"
	"github.com/yahyayesilyurt/league-simulation/internal/middleware"
	"github.com/yahyayesilyurt/league-simulation/internal/repository"
	"github.com/yahyayesilyurt/league-simulation/internal/service"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	appCache := cache.NewCache(redisClient)

	// Repositories
	teamRepo     := repository.NewTeamRepository(db)
	matchRepo    := repository.NewMatchRepository(db)
	standingRepo := repository.NewStandingRepository(db)

	// Services
	matchSvc      := service.NewMatchService(matchRepo, standingRepo, teamRepo, appCache)
	leagueSvc     := service.NewLeagueService(matchRepo, standingRepo, teamRepo, appCache)
	predictionSvc := service.NewPredictionService(standingRepo, matchRepo, teamRepo)
	fixtureSvc    := service.NewFixtureService(matchRepo, teamRepo)
	standingSvc   := service.NewStandingService(standingRepo, matchRepo, teamRepo)
	authSvc       := service.NewAuthService()

	// Handlers
	leagueHandler  := NewLeagueHandler(leagueSvc, predictionSvc)
	fixtureHandler := NewFixtureHandler(fixtureSvc)
	matchHandler   := NewMatchHandler(matchSvc, standingSvc)
	authHandler    := NewAuthHandler(authSvc)

	// Public routes 
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "League Simulation is running",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth
	r.POST("/auth/login", authHandler.Login)

	// League
	league := r.Group("/league")
	{
		league.GET("/table",          leagueHandler.GetStandings)
		league.GET("/fixtures",       leagueHandler.GetFixtures)
		league.GET("/week/:weekNo",   leagueHandler.GetWeek)
		league.GET("/predictions",    leagueHandler.GetPredictions)
		league.GET("/status",         leagueHandler.GetStatus)
		league.GET("/fixture-status", fixtureHandler.GetFixtureStatus)
		league.POST("/next-week",     leagueHandler.NextWeek)
		league.POST("/play-all",      leagueHandler.PlayAll)
	}

	// Match
	match := r.Group("/match")
	{
		match.GET("/:id", matchHandler.GetMatch)
	}

	// Protected routes
	admin := r.Group("/")
	admin.Use(middleware.AuthRequired(), middleware.AdminOnly())
	{
		admin.POST("/league/reset",            leagueHandler.Reset)
		admin.POST("/league/generate-fixture", fixtureHandler.GenerateFixture)
		admin.PUT("/match/:id/result", matchHandler.UpdateResult)
	}

	return r
}