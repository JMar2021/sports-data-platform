-- +goose Up

ALTER TABLE jobs
DROP CONSTRAINT jobs_key_key;

CREATE UNIQUE INDEX jobs_active_key_idx
ON jobs (key)
WHERE status IN ('pending', 'running', 'retrying');

-- +goose Down

DROP INDEX jobs_active_key_idx;

ALTER TABLE jobs
ADD CONSTRAINT jobs_key_key UNIQUE (key);