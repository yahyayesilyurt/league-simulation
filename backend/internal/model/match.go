package model

import "time"

type Match struct {
    ID         uint      `json:"id"           gorm:"primaryKey;autoIncrement"`
    Week       int       `json:"week"         gorm:"not null"`
    HomeTeamID uint      `json:"home_team_id" gorm:"not null"`
    AwayTeamID uint      `json:"away_team_id" gorm:"not null"`
    HomeGoals  *int      `json:"home_goals"`   
    AwayGoals  *int      `json:"away_goals"`  
    Played     bool      `json:"played"       gorm:"not null;default:false"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`

    // Relationships
    HomeTeam Team `json:"home_team" gorm:"foreignKey:HomeTeamID"`
    AwayTeam Team `json:"away_team" gorm:"foreignKey:AwayTeamID"`
}

func (Match) TableName() string {
    return "matches"
}