package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yahyayesilyurt/league-simulation/internal/service"
)

type LeagueHandler struct {
	leagueSvc     service.LeagueService
	predictionSvc service.PredictionService
}

func NewLeagueHandler(
	leagueSvc service.LeagueService,
	predictionSvc service.PredictionService,
) *LeagueHandler {
	return &LeagueHandler{
		leagueSvc:     leagueSvc,
		predictionSvc: predictionSvc,
	}
}

// GET /league/table
func (h *LeagueHandler) GetStandings(c *gin.Context) {
	standings, err := h.leagueSvc.GetStandings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"standings": standings})
}

// GET /league/fixtures
func (h *LeagueHandler) GetFixtures(c *gin.Context) {
	fixtures, err := h.leagueSvc.GetFixtures()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"fixtures": fixtures})
}

// GET /league/week/:weekNo
func (h *LeagueHandler) GetWeek(c *gin.Context) {
	weekNo, err := strconv.Atoi(c.Param("weekNo"))
	if err != nil || weekNo < 1 || weekNo > 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week number (1-6)"})
		return
	}
	matches, err := h.leagueSvc.GetWeek(weekNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"week": weekNo, "matches": matches})
}

// POST /league/next-week
func (h *LeagueHandler) NextWeek(c *gin.Context) {
	result, err := h.leagueSvc.NextWeek()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GET /league/status
func (h *LeagueHandler) GetStatus(c *gin.Context) {
	status, err := h.leagueSvc.GetStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

// POST /league/play-all
func (h *LeagueHandler) PlayAll(c *gin.Context) {
	result, err := h.leagueSvc.PlayAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// POST /league/reset
func (h *LeagueHandler) Reset(c *gin.Context) {
	status, err := h.leagueSvc.Reset()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "League reset successfully",
		"status":  status,
	})
}

// GET /league/predictions
func (h *LeagueHandler) GetPredictions(c *gin.Context) {
	predictions, err := h.predictionSvc.GetPredictions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"predictions": predictions})
}