-- +goose Up

ALTER TABLE standings
ADD COLUMN snapshot_date DATE NOT NULL DEFAULT CURRENT_DATE;

ALTER TABLE standings
DROP CONSTRAINT standings_team_id_season_key;

ALTER TABLE standings
ADD CONSTRAINT standings_team_id_season_snapshot_date_key
UNIQUE (team_id, season, snapshot_date);


-- +goose Down

ALTER TABLE standings
DROP CONSTRAINT standings_team_id_season_snapshot_date_key;

ALTER TABLE standings
ADD CONSTRAINT standings_team_id_season_key
UNIQUE (team_id, season);

ALTER TABLE standings
DROP COLUMN snapshot_date;