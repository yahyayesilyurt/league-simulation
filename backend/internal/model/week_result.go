package model

type WeekResult struct {
	Week           int         `json:"week"`
	Matches        []Match     `json:"matches"`
	Standings      []Standing  `json:"standings"`
	Predictions    interface{} `json:"predictions"`
	LeagueFinished bool        `json:"league_finished"`
	Message        string      `json:"message,omitempty"`
}