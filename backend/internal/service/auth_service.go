package service

import (
	"errors"
	"os"

	"github.com/yahyayesilyurt/league-simulation/config"
)

type AuthService interface {
	Login(username, password string) (string, error)
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Login(username, password string) (string, error) {
	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPass := os.Getenv("ADMIN_PASSWORD")

	if adminUser == "" {
		adminUser = "admin"
	}
	if adminPass == "" {
		adminPass = "admin123"
	}

	if username != adminUser || password != adminPass {
		return "", errors.New("invalid credentials")
	}

	token, err := config.GenerateToken(username, "admin")
	if err != nil {
		return "", err
	}

	return token, nil
}