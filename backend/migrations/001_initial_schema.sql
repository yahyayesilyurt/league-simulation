-- =====================================================
-- League Simulation - Initial Schema
-- =====================================================

-- -----------------------------------------------------
-- TEAMS - Teams and power levels
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS teams (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100)  NOT NULL UNIQUE,
    strength   INT           NOT NULL CHECK (strength BETWEEN 1 AND 100),
    created_at TIMESTAMP     NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP     NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE  teams            IS 'Teams in the league';
COMMENT ON COLUMN teams.strength   IS 'Team power score (1-100). Used in match simulation.';

-- -----------------------------------------------------
-- MATCHES - Fixtures and match results
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS matches (
    id           SERIAL PRIMARY KEY,
    week         INT          NOT NULL CHECK (week BETWEEN 1 AND 6),
    home_team_id INT          NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    away_team_id INT          NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    home_goals   INT          CHECK (home_goals >= 0),
    away_goals   INT          CHECK (away_goals >= 0),
    played       BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP    NOT NULL DEFAULT NOW(),

    -- The same two teams can only play each other once in the same week.
    CONSTRAINT unique_match_per_week UNIQUE (week, home_team_id, away_team_id),

    -- The team can't play with itself.
    CONSTRAINT different_teams CHECK (home_team_id != away_team_id)
);

COMMENT ON TABLE  matches            IS 'League fixtures and match results.';
COMMENT ON COLUMN matches.week       IS 'Week number (1-6). 4 teams, double round-robin = 6 weeks';
COMMENT ON COLUMN matches.played     IS 'Whether the match was played or not. FALSE = not played';
COMMENT ON COLUMN matches.home_goals IS 'Home team goals tally. NULL = match not played';
COMMENT ON COLUMN matches.away_goals IS 'Away goals count. NULL = no match played';

-- -----------------------------------------------------
-- STANDINGS - Scoreboard
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS standings (
    id         SERIAL PRIMARY KEY,
    team_id    INT       NOT NULL UNIQUE REFERENCES teams(id) ON DELETE CASCADE,
    played     INT       NOT NULL DEFAULT 0,
    won        INT       NOT NULL DEFAULT 0,
    drawn      INT       NOT NULL DEFAULT 0,
    lost       INT       NOT NULL DEFAULT 0,
    goals_for      INT   NOT NULL DEFAULT 0,
    goals_against  INT   NOT NULL DEFAULT 0,
    goal_diff      INT   NOT NULL DEFAULT 0,
    points     INT       NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE  standings              IS 'League table — Premier League rules';
COMMENT ON COLUMN standings.goals_for    IS 'Goal scored (GF)';
COMMENT ON COLUMN standings.goals_against IS 'Goal conceded (GA)';
COMMENT ON COLUMN standings.goal_diff    IS 'Goal average = GF - GA (GD)';
COMMENT ON COLUMN standings.points       IS 'Wins=3, Draws=1, Losses=0';

-- -----------------------------------------------------
-- INDEX — Query performance
-- -----------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_matches_week
    ON matches(week);

CREATE INDEX IF NOT EXISTS idx_matches_played
    ON matches(played);

CREATE INDEX IF NOT EXISTS idx_matches_home_team
    ON matches(home_team_id);

CREATE INDEX IF NOT EXISTS idx_matches_away_team
    ON matches(away_team_id);

CREATE INDEX IF NOT EXISTS idx_standings_points
    ON standings(points DESC, goal_diff DESC, goals_for DESC);

-- -----------------------------------------------------
-- updated_at automatic update function
-- -----------------------------------------------------
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER teams_updated_at
    BEFORE UPDATE ON teams
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER matches_updated_at
    BEFORE UPDATE ON matches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER standings_updated_at
    BEFORE UPDATE ON standings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();