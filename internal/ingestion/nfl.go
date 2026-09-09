package ingestion

import (
	"context"
	"fmt"

	"github.com/JMar2021/sports-data-platform/internal/nfl"
	"github.com/JMar2021/sports-data-platform/internal/repository"
)

type NFLIngestor struct {
	NFLClient  *nfl.Client
	Repository *repository.Repository
}

func (m *NFLIngestor) IngestSchedule(ctx context.Context, date string) error {
	scheduleResponse, err := m.NFLClient.GetSchedule(ctx, date)
	if err != nil {
		return err
	}

	if len(scheduleResponse.Events) == 0 {
		fmt.Println("No games scheduled for this date.")
		return nil
	}

	for _, event := range scheduleResponse.Events {
		for _, competition := range event.Competitions {
			domainGame, err := nfl.NormalizeGame(competition)
			if err != nil {
				return err
			}

			fmt.Printf(
				"%s %d @ %s %d\n",
				domainGame.AwayTeam.Name,
				domainGame.AwayScore,
				domainGame.HomeTeam.Name,
				domainGame.HomeScore,
			)

			awayTeamID, err := m.Repository.UpsertTeam(
				ctx,
				domainGame.AwayTeam.Sport,
				domainGame.AwayTeam.ExternalID,
				domainGame.AwayTeam.Name,
			)
			if err != nil {
				return err
			}

			homeTeamID, err := m.Repository.UpsertTeam(
				ctx,
				domainGame.HomeTeam.Sport,
				domainGame.HomeTeam.ExternalID,
				domainGame.HomeTeam.Name,
			)
			if err != nil {
				return err
			}

			fmt.Printf(
				"%s → DB team ID %d\n",
				domainGame.AwayTeam.Name,
				awayTeamID,
			)

			fmt.Printf(
				"%s → DB team ID %d\n",
				domainGame.HomeTeam.Name,
				homeTeamID,
			)

			err = m.Repository.UpsertGame(
				ctx,
				domainGame,
				homeTeamID,
				awayTeamID,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
