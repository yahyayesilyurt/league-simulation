package repository

import (
	"github.com/yahyayesilyurt/league-simulation/internal/model"
	"gorm.io/gorm"
)

type TeamRepository interface {
    GetAll() ([]model.Team, error)
    GetByID(id uint) (*model.Team, error)
    Create(team *model.Team) error
    Update(team *model.Team) error
}

type teamRepository struct {
    db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
    return &teamRepository{db: db}
}

func (r *teamRepository) GetAll() ([]model.Team, error) {
    var teams []model.Team
    result := r.db.Order("id ASC").Find(&teams)
    return teams, result.Error
}

func (r *teamRepository) GetByID(id uint) (*model.Team, error) {
    var team model.Team
    result := r.db.First(&team, id)
    return &team, result.Error
}

func (r *teamRepository) Create(team *model.Team) error {
    return r.db.Create(team).Error
}

func (r *teamRepository) Update(team *model.Team) error {
    return r.db.Save(team).Error
}