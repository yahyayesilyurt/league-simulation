package repository

import (
	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"gorm.io/gorm"
)

type StandingRepository interface {
    GetAll() ([]model.Standing, error)
    GetByTeamID(teamID uint) (*model.Standing, error)
    Update(standing *model.Standing) error
    ResetAll() error
}

type standingRepository struct {
    db *gorm.DB
}

func NewStandingRepository(db *gorm.DB) StandingRepository {
    return &standingRepository{db: db}
}

func (r *standingRepository) GetAll() ([]model.Standing, error) {
    var standings []model.Standing
    result := r.db.
        Preload("Team").
        Order("points DESC, goal_diff DESC, goals_for DESC").
        Find(&standings)
    return standings, result.Error
}

func (r *standingRepository) GetByTeamID(teamID uint) (*model.Standing, error) {
    var standing model.Standing
    result := r.db.
        Preload("Team").
        Where("team_id = ?", teamID).
        First(&standing)
    return &standing, result.Error
}

func (r *standingRepository) Update(standing *model.Standing) error {
    return r.db.Save(standing).Error
}

func (r *standingRepository) ResetAll() error {
    return r.db.Exec(`
        UPDATE standings SET
            played=0, won=0, drawn=0, lost=0,
            goals_for=0, goals_against=0,
            goal_diff=0, points=0
    `).Error
}