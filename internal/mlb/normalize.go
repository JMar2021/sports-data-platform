package mlb

import (
	"fmt"
	"strconv"
	"time"

	"github.com/JMar2021/sports-data-platform/internal/domain"
)

func NormalizeGame(game Game) (domain.Game, error) {
	scheduledAt, err := time.Parse(time.RFC3339, game.GameDate)
	if err != nil {
		return domain.Game{}, fmt.Errorf(
			"parse game date %q: %w",
			game.GameDate,
			err,
		)
	}

	return domain.Game{
		Sport:       domain.SportMLB,
		ExternalID:  strconv.Itoa(game.GamePk),
		ScheduledAt: scheduledAt,
		HomeTeam:    NormalizeTeam(game.Teams.Home.Team),
		AwayTeam:    NormalizeTeam(game.Teams.Away.Team),
		HomeScore:   game.Teams.Home.Score,
		AwayScore:   game.Teams.Away.Score,
		Status:      game.Status.AbstractGameState,
		Venue:       game.Venue.Name,
	}, nil
}

func NormalizeTeam(team Team) domain.Team {
	return domain.Team{Sport: domain.SportMLB, ExternalID: strconv.Itoa(team.ID), Name: team.Name}
}
