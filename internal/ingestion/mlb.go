package ingestion

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/JMar2021/sports-data-platform/internal/domain"
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
		domainGame, err := mlb.NormalizeGame(game)
		if err != nil {
			return err
		}
		awayTeamID, err := m.Repository.UpsertTeam(ctx, domainGame.AwayTeam.Sport, domainGame.AwayTeam.ExternalID, domainGame.AwayTeam.Name)
		if err != nil {
			return err
		}
		homeTeamID, err := m.Repository.UpsertTeam(ctx, domainGame.HomeTeam.Sport, domainGame.HomeTeam.ExternalID, domainGame.HomeTeam.Name)
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
			domainGame,
			homeTeamID,
			awayTeamID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MLBIngestor) IngestStandings(ctx context.Context, date string) error {
	season := time.Now().Year()

	standingsResponse, err := m.MLBClient.GetStandings(
		ctx,
		strconv.Itoa(season),
	)
	if err != nil {
		return err
	}

	snapshotDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("invalid snapshot date %q: %w", date, err)
	}

	for _, group := range standingsResponse.Records {
		for _, record := range group.TeamRecords {

			teamID, err := m.Repository.UpsertTeam(
				ctx,
				domain.SportMLB,
				strconv.Itoa(record.Team.ID),
				record.Team.Name,
			)
			if err != nil {
				return fmt.Errorf("upserting team %d: %w", record.Team.ID, err)
			}

			divisionRank, err := strconv.Atoi(record.DivisionRank)
			if err != nil {
				return fmt.Errorf(
					"parsing division rank for team %d: %w",
					record.Team.ID,
					err,
				)
			}

			leagueRank, err := strconv.Atoi(record.LeagueRank)
			if err != nil {
				return fmt.Errorf(
					"parsing league rank for team %d: %w",
					record.Team.ID,
					err,
				)
			}

			var gamesBack float64

			if record.GamesBack == "-" {
				gamesBack = 0
			} else {
				gamesBack, err = strconv.ParseFloat(record.GamesBack, 64)
				if err != nil {
					return fmt.Errorf(
						"parsing games back for team %d: %w",
						record.Team.ID,
						err,
					)
				}
			}

			winningPercentage, err := strconv.ParseFloat(
				record.WinningPercentage,
				64,
			)
			if err != nil {
				return fmt.Errorf(
					"parsing winning percentage for team %d: %w",
					record.Team.ID,
					err,
				)
			}

			lastUpdated, err := time.Parse(
				time.RFC3339,
				record.LastUpdated,
			)
			if err != nil {
				return fmt.Errorf(
					"parsing last updated for team %d: %w",
					record.Team.ID,
					err,
				)
			}

			err = m.Repository.UpsertStandings(
				ctx,
				teamID,
				season,
				snapshotDate,
				divisionRank,
				leagueRank,
				record.GamesPlayed,
				gamesBack,
				record.Wins,
				record.Losses,
				winningPercentage,
				record.RunsScored,
				record.RunsAllowed,
				record.RunDifferential,
				lastUpdated,
			)
			if err != nil {
				return fmt.Errorf(
					"upserting standings for team %d: %w",
					record.Team.ID,
					err,
				)
			}

		}
	}

	return nil

}
