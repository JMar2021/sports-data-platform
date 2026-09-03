-- +goose Up
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    mlb_id INTEGER NOT NULL UNIQUE,
    name TEXT NOT NULL
);

CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    mlb_game_id INTEGER NOT NULL UNIQUE,
    game_date TIMESTAMP NOT NULL,
    away_team_id INTEGER NOT NULL REFERENCES teams(id),
    home_team_id INTEGER NOT NULL REFERENCES teams(id),
    away_score INTEGER NOT NULL,
    home_score INTEGER NOT NULL,
    status TEXT NOT NULL,
    venue TEXT
);

-- +goose Down

DROP TABLE games;
DROP TABLE teams;