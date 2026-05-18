package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yahyayesilyurt/league-simulation/internal/service"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// Login godoc
// @Summary      Admin login
// @Description  You can take JWT tokens with your username and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      object{username=string,password=string}  true  "Username and password"
// @Success      200          {object}  object{token=string,message=string}
// @Failure      400          {object}  object{error=string}
// @Failure      401          {object}  object{error=string}
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username and password are required",
		})
		return
	}

	token, err := h.authSvc.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "Login successful",
	})
}