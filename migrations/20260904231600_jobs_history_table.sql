-- +goose Up

CREATE TABLE jobs (
    id TEXT PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    sport TEXT NOT NULL,
    operation TEXT NOT NULL,
    job_date DATE NOT NULL,
    status TEXT NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    error TEXT
);

-- +goose Down

DROP TABLE jobs;