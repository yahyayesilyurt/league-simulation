package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yahyayesilyurt/league-simulation/internal/service"
)

type MatchHandler struct {
	matchSvc    service.MatchService
	standingSvc service.StandingService
}

func NewMatchHandler(
	matchSvc service.MatchService,
	standingSvc service.StandingService,
) *MatchHandler {
	return &MatchHandler{
		matchSvc:    matchSvc,
		standingSvc: standingSvc,
	}
}

// PUT /match/:id/result
func (h *MatchHandler) UpdateResult(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	var req struct {
		HomeGoals int `json:"home_goals" binding:"min=0"`
		AwayGoals int `json:"away_goals" binding:"min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body. home_goals and away_goals must be >= 0",
		})
		return
	}

	match, err := h.matchSvc.UpdateMatchResult(uint(id), req.HomeGoals, req.AwayGoals)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	standings, err := h.standingSvc.GetStandings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Match result updated successfully",
		"match":     match,
		"standings": standings,
	})
}

// GET /match/:id
func (h *MatchHandler) GetMatch(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	match, err := h.matchSvc.GetMatchByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"match": match})
}