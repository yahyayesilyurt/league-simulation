package model

import "time"

type Standing struct {
    ID           uint      `json:"id"            gorm:"primaryKey;autoIncrement"`
    TeamID       uint      `json:"team_id"       gorm:"not null;uniqueIndex"`
    Played       int       `json:"played"        gorm:"not null;default:0"`
    Won          int       `json:"won"           gorm:"not null;default:0"`
    Drawn        int       `json:"drawn"         gorm:"not null;default:0"`
    Lost         int       `json:"lost"          gorm:"not null;default:0"`
    GoalsFor     int       `json:"goals_for"     gorm:"not null;default:0"`
    GoalsAgainst int       `json:"goals_against" gorm:"not null;default:0"`
    GoalDiff     int       `json:"goal_diff"     gorm:"not null;default:0"`
    Points       int       `json:"points"        gorm:"not null;default:0"`
    UpdatedAt    time.Time `json:"updated_at"`

    // Relationships
    Team Team `json:"team" gorm:"foreignKey:TeamID"`
}

func (Standing) TableName() string {
    return "standings"
}