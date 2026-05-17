package model

type WeekResult struct {
	Week           int         `json:"week"`
	Matches        []Match     `json:"matches"`
	Standings      []Standing  `json:"standings"`
	Predictions    interface{} `json:"predictions"`
	LeagueFinished bool        `json:"league_finished"`
	Message        string      `json:"message,omitempty"`
}

type LeagueStatus struct {
	CurrentWeek    int    `json:"current_week"`
	TotalWeeks     int    `json:"total_weeks"`
	LeagueFinished bool   `json:"league_finished"`
	MatchesPlayed  int    `json:"matches_played"`
	MatchesLeft    int    `json:"matches_left"`
	Status         string `json:"status"` // "not_started", "in_progress", "finished"
}

type LeagueSummary struct {
	Champion       *Standing  `json:"champion"`
	FinalStandings []Standing `json:"final_standings"`
	TopScorer      string     `json:"top_scorer_team"`
	BestDefense    string     `json:"best_defense_team"`
	TotalGoals     int        `json:"total_goals"`
	TotalMatches   int        `json:"total_matches"`
}

type PlayAllResult struct {
	TotalWeeksPlayed int            `json:"total_weeks_played"`
	Weeks            []WeekResult   `json:"weeks"`
	Summary          *LeagueSummary `json:"summary"`
}