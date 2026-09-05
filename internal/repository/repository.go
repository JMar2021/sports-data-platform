package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

type GameResult struct {
	GameID    int
	AwayTeam  string
	AwayScore int
	HomeTeam  string
	HomeScore int
	Status    string
	Venue     string
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) GetDatabaseVersion(ctx context.Context) (string, error) {
	var version string
	err := r.DB.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		return "", err
	}
	return version, nil
}

func (r *Repository) UpsertTeam(ctx context.Context, mlbID int, name string) (int, error) {
	var teamId int
	err := r.DB.QueryRow(ctx,
		`INSERT INTO teams (mlb_id, name)
         VALUES ($1, $2)
         ON CONFLICT (mlb_id)
         DO UPDATE SET name = EXCLUDED.name
         RETURNING id`,
		mlbID, name).Scan(&teamId)
	if err != nil {
		return 0, err
	}
	return teamId, nil
}

func (r *Repository) UpsertGame(ctx context.Context, mlbID int, gameDate time.Time, awayTeamID int, homeTeamID int, awayScore int, homeScore int, status string, venue string) error {
	_, err := r.DB.Exec(ctx,
		`INSERT INTO games (mlb_game_id, game_date, away_team_id, home_team_id, away_score, home_score, status, venue)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (mlb_game_id)
		 DO UPDATE SET game_date = EXCLUDED.game_date,
					   home_score = EXCLUDED.home_score,
					   away_score = EXCLUDED.away_score,
					   status = EXCLUDED.status,
					   venue = EXCLUDED.venue`,
		mlbID, gameDate, awayTeamID, homeTeamID, awayScore, homeScore, status, venue)
	return err
}

func (r *Repository) GetGames(ctx context.Context) ([]GameResult, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT
			g.mlb_game_id,
			away.name AS away_team,
			g.away_score,
			home.name AS home_team,
			g.home_score,
			g.status,
			g.venue
		FROM games g
		JOIN teams away ON g.away_team_id = away.id
		JOIN teams home ON g.home_team_id = home.id
		ORDER BY g.game_date`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []GameResult

	for rows.Next() {
		var game GameResult

		err := rows.Scan(
			&game.GameID,
			&game.AwayTeam,
			&game.AwayScore,
			&game.HomeTeam,
			&game.HomeScore,
			&game.Status,
			&game.Venue,
		)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}

func (r *Repository) UpsertStandings(
	ctx context.Context,
	teamID int,
	season int,
	snapshotDate time.Time,
	divisionRank int,
	leagueRank int,
	gamesPlayed int,
	gamesBack float64,
	wins int,
	losses int,
	winningPercentage float64,
	runsScored int,
	runsAllowed int,
	runDifferential int,
	lastUpdated time.Time,
) error {
	_, err := r.DB.Exec(
		ctx,
		`INSERT INTO standings (
			team_id,
			season,
			snapshot_date,
			division_rank,
			league_rank,
			games_played,
			games_back,
			wins,
			losses,
			winning_percentage,
			runs_scored,
			runs_allowed,
			run_differential,
			last_updated
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13, $14
		)
		ON CONFLICT (team_id, season, snapshot_date)
		DO UPDATE SET
			division_rank = EXCLUDED.division_rank,
			league_rank = EXCLUDED.league_rank,
			games_played = EXCLUDED.games_played,
			games_back = EXCLUDED.games_back,
			wins = EXCLUDED.wins,
			losses = EXCLUDED.losses,
			winning_percentage = EXCLUDED.winning_percentage,
			runs_scored = EXCLUDED.runs_scored,
			runs_allowed = EXCLUDED.runs_allowed,
			run_differential = EXCLUDED.run_differential,
			last_updated = EXCLUDED.last_updated`,
		teamID,
		season,
		snapshotDate,
		divisionRank,
		leagueRank,
		gamesPlayed,
		gamesBack,
		wins,
		losses,
		winningPercentage,
		runsScored,
		runsAllowed,
		runDifferential,
		lastUpdated,
	)

	return err
}

func (r *Repository) CreateJob(
	ctx context.Context,
	id string,
	key string,
	sport string,
	operation string,
	jobDate time.Time,
	status string,
	attempts int,
	createdAt time.Time,
) (bool, error) {
	tag, err := r.DB.Exec(ctx, `
		INSERT INTO jobs (
			id,
			key,
			sport,
			operation,
			job_date,
			status,
			attempts,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT DO NOTHING
	`,
		id,
		key,
		sport,
		operation,
		jobDate,
		status,
		attempts,
		createdAt,
	)

	if err != nil {
		return false, err
	}

	return tag.RowsAffected() == 1, nil
}

func (r *Repository) StartJob(
	ctx context.Context,
	id string,
	startedAt time.Time,
) error {
	tag, err := r.DB.Exec(ctx, `
		UPDATE jobs
		SET status = $1,
			started_at = $2,
			error = NULL
		WHERE id = $3
		  AND status IN ('pending', 'retrying')
	`,
		"running",
		startedAt,
		id,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("job %s cannot transition to running", id)
	}

	return nil
}

func (r *Repository) RetryJob(
	ctx context.Context,
	id string,
	attempts int,
	errMessage string,
) error {
	tag, err := r.DB.Exec(ctx, `
		UPDATE jobs
		SET status = $1,
			attempts = $2,
			error = $3
		WHERE id = $4
		  AND status = 'running'
	`,
		"retrying",
		attempts,
		errMessage,
		id,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("job %s cannot transition to retrying", id)
	}

	return nil
}

func (r *Repository) CompleteJob(
	ctx context.Context,
	id string,
	completedAt time.Time,
) error {
	tag, err := r.DB.Exec(ctx, `
		UPDATE jobs
		SET status = $1,
			completed_at = $2,
			error = NULL
		WHERE id = $3
		  AND status = 'running'
	`,
		"completed",
		completedAt,
		id,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("job %s cannot transition to completed", id)
	}

	return nil
}

func (r *Repository) FailJob(
	ctx context.Context,
	id string,
	attempts int,
	errMessage string,
	completedAt time.Time,
) error {
	tag, err := r.DB.Exec(ctx, `
		UPDATE jobs
		SET status = $1,
			attempts = $2,
			error = $3,
			completed_at = $4
		WHERE id = $5
		  AND status IN ('running', 'pending', 'retrying')
	`,
		"failed",
		attempts,
		errMessage,
		completedAt,
		id,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("job %s cannot transition to failed", id)
	}

	return nil
}

type StaleJob struct {
	ID        string
	Key       string
	Sport     string
	Operation string
	Date      time.Time
	Attempts  int
}

func (r *Repository) RecoverStaleJobs(
	ctx context.Context,
	staleBefore time.Time,
) ([]StaleJob, error) {
	rows, err := r.DB.Query(ctx, `
		UPDATE jobs
		SET status = 'retrying',
			attempts = attempts + 1,
			error = 'recovered after worker crash'
		WHERE status = 'running'
		  AND started_at < $1
		RETURNING
			id,
			key,
			sport,
			operation,
			job_date,
			attempts
	`,
		staleBefore,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []StaleJob

	for rows.Next() {
		var job StaleJob

		if err := rows.Scan(
			&job.ID,
			&job.Key,
			&job.Sport,
			&job.Operation,
			&job.Date,
			&job.Attempts,
		); err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}
