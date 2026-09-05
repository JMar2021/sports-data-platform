-- +goose Up

CREATE TABLE standings (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id),
    season INTEGER NOT NULL,
    division_rank INTEGER NOT NULL,
    league_rank INTEGER NOT NULL,
    games_played INTEGER NOT NULL,
    games_back NUMERIC NOT NULL,
    wins INTEGER NOT NULL,
    losses INTEGER NOT NULL,
    winning_percentage NUMERIC NOT NULL,
    runs_scored INTEGER NOT NULL,
    runs_allowed INTEGER NOT NULL,
    run_differential INTEGER NOT NULL,
    last_updated TIMESTAMP NOT NULL,

    UNIQUE (team_id, season)
);

-- +goose Down

DROP TABLE standings;