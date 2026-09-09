-- +goose Up

-- Teams: migrate from MLB-specific identifiers to
-- sport + external_id.

ALTER TABLE teams
    ADD COLUMN sport TEXT,
    ADD COLUMN external_id TEXT;

UPDATE teams
SET
    sport = 'mlb',
    external_id = mlb_id::TEXT;

ALTER TABLE teams
    ALTER COLUMN sport SET NOT NULL,
    ALTER COLUMN external_id SET NOT NULL;

ALTER TABLE teams
    DROP CONSTRAINT teams_mlb_id_key;

ALTER TABLE teams
    DROP COLUMN mlb_id;

ALTER TABLE teams
    ADD CONSTRAINT teams_sport_external_id_key
    UNIQUE (sport, external_id);


-- Games: migrate from MLB-specific identifiers to
-- sport + external_id.

ALTER TABLE games
    ADD COLUMN sport TEXT,
    ADD COLUMN external_id TEXT;

UPDATE games
SET
    sport = 'mlb',
    external_id = mlb_game_id::TEXT;

ALTER TABLE games
    ALTER COLUMN sport SET NOT NULL,
    ALTER COLUMN external_id SET NOT NULL;

ALTER TABLE games
    DROP CONSTRAINT games_mlb_game_id_key;

ALTER TABLE games
    DROP COLUMN mlb_game_id;

ALTER TABLE games
    ADD CONSTRAINT games_sport_external_id_key
    UNIQUE (sport, external_id);


-- +goose Down

ALTER TABLE games
    ADD COLUMN mlb_game_id INTEGER;

UPDATE games
SET mlb_game_id = external_id::INTEGER
WHERE sport = 'mlb';

ALTER TABLE games
    ALTER COLUMN mlb_game_id SET NOT NULL;

ALTER TABLE games
    DROP CONSTRAINT games_sport_external_id_key;

ALTER TABLE games
    DROP COLUMN sport,
    DROP COLUMN external_id;

ALTER TABLE games
    ADD CONSTRAINT games_mlb_game_id_key
    UNIQUE (mlb_game_id);


ALTER TABLE teams
    ADD COLUMN mlb_id INTEGER;

UPDATE teams
SET mlb_id = external_id::INTEGER
WHERE sport = 'mlb';

ALTER TABLE teams
    ALTER COLUMN mlb_id SET NOT NULL;

ALTER TABLE teams
    DROP CONSTRAINT teams_sport_external_id_key;

ALTER TABLE teams
    DROP COLUMN sport,
    DROP COLUMN external_id;

ALTER TABLE teams
    ADD CONSTRAINT teams_mlb_id_key
    UNIQUE (mlb_id);
