package nfl_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/JMar2021/sports-data-platform/internal/nfl"
)

func TestGetSchedule(t *testing.T) {
	client := nfl.NewClient(http.DefaultClient)

	ctx := context.Background()

	schedule, err := client.GetSchedule(ctx, "20260909")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Events returned: %d", len(schedule.Events))

	for _, event := range schedule.Events {
		t.Logf("Event: %+v", event)

		for _, competition := range event.Competitions {
			game, err := nfl.NormalizeGame(competition)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("Normalized game: %+v", game)

			for _, competitor := range competition.Competitors {
				t.Logf(
					"%s: %s (%s)",
					competitor.HomeAway,
					competitor.Team.Name,
					competitor.Score,
				)
			}
		}
	}
}
