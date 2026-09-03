package repository

import (
	"context"
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
