-- +goose Up

ALTER TABLE teams
ADD COLUMN description TEXT;

-- +goose Down

ALTER TABLE teams
DROP COLUMN description;