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

// GetStandings godoc
// @Summary      Scoreboard
// @Description  Returns a league table sorted according to Premier League rules.
// @Tags         league
// @Produce      json
// @Success      200  {object}  object{standings=[]model.Standing}
// @Failure      500  {object}  object{error=string}
// @Router       /league/table [get]
func (h *LeagueHandler) GetStandings(c *gin.Context) {
	standings, err := h.leagueSvc.GetStandings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"standings": standings})
}

// GetFixtures godoc
// @Summary      All fixtures
// @Description  It returns the entire match schedule for 6 weeks.
// @Tags         league
// @Produce      json
// @Success      200  {object}  object{fixtures=[]model.Match}
// @Failure      500  {object}  object{error=string}
// @Router       /league/fixtures [get]
func (h *LeagueHandler) GetFixtures(c *gin.Context) {
	fixtures, err := h.leagueSvc.GetFixtures()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"fixtures": fixtures})
}

// GetWeek godoc
// @Summary      Week details
// @Description  It returns the matches of the specified week.
// @Tags         league
// @Produce      json
// @Param        weekNo  path      int  true  "Week number (1-6)"
// @Success      200     {object}  object{week=int,matches=[]model.Match}
// @Failure      400     {object}  object{error=string}
// @Router       /league/week/{weekNo} [get]
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

// NextWeek godoc
// @Summary      Play next week
// @Description  It simulates next week's matches and updates the league table.
// @Tags         league
// @Produce      json
// @Success      200  {object}  model.WeekResult
// @Failure      400  {object}  object{error=string}
// @Router       /league/next-week [post]
func (h *LeagueHandler) NextWeek(c *gin.Context) {
	result, err := h.leagueSvc.NextWeek()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetStatus godoc
// @Summary      League status
// @Description  Returns the current state of the league (not_started, in_progress, finished)
// @Tags         league
// @Produce      json
// @Success      200  {object}  model.LeagueStatus
// @Failure      500  {object}  object{error=string}
// @Router       /league/status [get]
func (h *LeagueHandler) GetStatus(c *gin.Context) {
	status, err := h.leagueSvc.GetStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

// PlayAll godoc
// @Summary      Play all remaining weeks
// @Description  Plays through all remaining weeks and returns results for each week.
// @Tags         league
// @Produce      json
// @Success      200  {object}  model.PlayAllResult
// @Failure      400  {object}  object{error=string}
// @Router       /league/play-all [post]
func (h *LeagueHandler) PlayAll(c *gin.Context) {
	result, err := h.leagueSvc.PlayAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Reset godoc
// @Summary      Reset the league (Admin)
// @Description  It deletes all matches, resets the standings, and recreates the fixtures.
// @Tags         league
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  object{message=string,status=model.LeagueStatus}
// @Failure      401  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /league/reset [post]
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

// GetPredictions godoc
// @Summary      Championship predictions
// @Description  Returns each team's championship percentage starting from week 4.
// @Tags         league
// @Produce      json
// @Success      200  {object}  object{predictions=[]service.TeamPrediction}
// @Failure      500  {object}  object{error=string}
// @Router       /league/predictions [get]
func (h *LeagueHandler) GetPredictions(c *gin.Context) {
	predictions, err := h.predictionSvc.GetPredictions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"predictions": predictions})
}