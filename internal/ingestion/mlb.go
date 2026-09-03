package ingestion

import (
	"context"
	"fmt"
	"time"

	"github.com/JMar2021/sports-data-platform/internal/mlb"
	"github.com/JMar2021/sports-data-platform/internal/repository"
)

type MLBIngestor struct {
	MLBClient  *mlb.Client
	Repository *repository.Repository
}

func (m *MLBIngestor) IngestSchedule(ctx context.Context, date string) error {
	scheduleResponse, err := m.MLBClient.GetSchedule(ctx, date)
	if err != nil {
		return err
	}

	if len(scheduleResponse.Dates) == 0 {
		fmt.Println("No games scheduled for this date.")
		return nil
	}
	for _, game := range scheduleResponse.Dates[0].Games {
		fmt.Printf("%s %d @ %s %d\n", game.Teams.Away.Team.Name, game.Teams.Away.Score, game.Teams.Home.Team.Name, game.Teams.Home.Score)
		awayTeam := game.Teams.Away.Team
		homeTeam := game.Teams.Home.Team
		awayTeamID, err := m.Repository.UpsertTeam(ctx, awayTeam.ID, awayTeam.Name)
		if err != nil {
			return err
		}
		homeTeamID, err := m.Repository.UpsertTeam(ctx, homeTeam.ID, homeTeam.Name)
		if err != nil {
			return err
		}
		gameDate, err := time.Parse(time.RFC3339, game.GameDate)
		if err != nil {
			return err
		}

		fmt.Printf(
			"%s → DB team ID %d\n",
			awayTeam.Name,
			awayTeamID,
		)

		fmt.Printf(
			"%s → DB team ID %d\n",
			homeTeam.Name,
			homeTeamID,
		)
		err = m.Repository.UpsertGame(
			ctx,
			game.GamePk,
			gameDate,
			awayTeamID,
			homeTeamID,
			game.Teams.Away.Score,
			game.Teams.Home.Score,
			game.Status.AbstractGameState,
			game.Venue.Name,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
