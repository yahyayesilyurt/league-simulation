package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yahyayesilyurt/league-simulation/internal/service"
)

type FixtureHandler struct {
	fixtureSvc service.FixtureService
}

func NewFixtureHandler(fixtureSvc service.FixtureService) *FixtureHandler {
	return &FixtureHandler{fixtureSvc: fixtureSvc}
}

// POST /league/generate-fixture
func (h *FixtureHandler) GenerateFixture(c *gin.Context) {
	if err := h.fixtureSvc.GenerateFixture(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Fixture generated successfully (6 weeks, 12 matches)",
	})
}

// GET /league/fixture-status
func (h *FixtureHandler) GetFixtureStatus(c *gin.Context) {
	generated, err := h.fixtureSvc.IsFixtureGenerated()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"fixture_generated": generated})
}