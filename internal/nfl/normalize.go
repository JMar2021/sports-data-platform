package nfl

import (
	"fmt"
	"strconv"
	"time"

	"github.com/JMar2021/sports-data-platform/internal/domain"
)

func NormalizeGame(competition Competition) (domain.Game, error) {
	if len(competition.Competitors) != 2 {
		return domain.Game{}, fmt.Errorf(
			"expected 2 competitors, got %d",
			len(competition.Competitors),
		)
	}

	scheduledAt, err := time.Parse("2006-01-02T15:04Z", competition.Date)
	if err != nil {
		return domain.Game{}, fmt.Errorf(
			"parse game date %q: %w",
			competition.Date,
			err,
		)
	}

	var homeTeam domain.Team
	var awayTeam domain.Team
	var homeScore int
	var awayScore int

	for _, competitor := range competition.Competitors {
		score, err := strconv.Atoi(competitor.Score)
		if err != nil {
			return domain.Game{}, fmt.Errorf(
				"parse score %q for team %q: %w",
				competitor.Score,
				competitor.Team.Name,
				err,
			)
		}

		switch competitor.HomeAway {
		case "home":
			homeTeam = NormalizeTeam(competitor.Team)
			homeScore = score

		case "away":
			awayTeam = NormalizeTeam(competitor.Team)
			awayScore = score

		default:
			return domain.Game{}, fmt.Errorf(
				"unknown homeAway value %q",
				competitor.HomeAway,
			)
		}
	}

	status := "Not started"
	if competition.Status.Type.Completed {
		status = "Done"
	}

	return domain.Game{
		Sport:       domain.SportNFL,
		ExternalID:  competition.ID,
		ScheduledAt: scheduledAt,
		HomeTeam:    homeTeam,
		AwayTeam:    awayTeam,
		HomeScore:   homeScore,
		AwayScore:   awayScore,
		Status:      status,
		Venue:       competition.Venue.Name,
	}, nil
}

func NormalizeTeam(team Team) domain.Team {
	return domain.Team{
		Sport:      domain.SportNFL,
		ExternalID: team.ID,
		Name:       team.Name,
	}
}
