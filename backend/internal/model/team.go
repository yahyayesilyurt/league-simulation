package model

import "time"

type Team struct {
    ID        uint      `json:"id"         gorm:"primaryKey;autoIncrement"`
    Name      string    `json:"name"       gorm:"size:100;not null;uniqueIndex"`
    Strength  int       `json:"strength"   gorm:"not null;check:strength BETWEEN 1 AND 100"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (Team) TableName() string {
    return "teams"
}