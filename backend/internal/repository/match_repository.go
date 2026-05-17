package repository

import (
	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"gorm.io/gorm"
)

type MatchRepository interface {
    GetAll() ([]model.Match, error)
    GetByWeek(week int) ([]model.Match, error)
    GetUnplayed() ([]model.Match, error)
    GetByID(id uint) (*model.Match, error)
    Create(match *model.Match) error
    Update(match *model.Match) error
    DeleteAll() error
}

type matchRepository struct {
    db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
    return &matchRepository{db: db}
}

func (r *matchRepository) GetAll() ([]model.Match, error) {
    var matches []model.Match
    result := r.db.
        Preload("HomeTeam").
        Preload("AwayTeam").
        Order("week ASC, id ASC").
        Find(&matches)
    return matches, result.Error
}

func (r *matchRepository) GetByWeek(week int) ([]model.Match, error) {
    var matches []model.Match
    result := r.db.
        Preload("HomeTeam").
        Preload("AwayTeam").
        Where("week = ?", week).
        Order("id ASC").
        Find(&matches)
    return matches, result.Error
}

func (r *matchRepository) GetUnplayed() ([]model.Match, error) {
    var matches []model.Match
    result := r.db.
        Preload("HomeTeam").
        Preload("AwayTeam").
        Where("played = ?", false).
        Order("week ASC, id ASC").
        Find(&matches)
    return matches, result.Error
}

func (r *matchRepository) GetByID(id uint) (*model.Match, error) {
    var match model.Match
    result := r.db.
        Preload("HomeTeam").
        Preload("AwayTeam").
        First(&match, id)
    return &match, result.Error
}

func (r *matchRepository) Create(match *model.Match) error {
    return r.db.Create(match).Error
}

func (r *matchRepository) Update(match *model.Match) error {
    return r.db.Save(match).Error
}

func (r *matchRepository) DeleteAll() error {
    return r.db.Exec("DELETE FROM matches").Error
}