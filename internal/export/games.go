package export

import (
	"encoding/json"
	"os"

	"github.com/JMar2021/sports-data-platform/internal/domain"
)

type GameData struct {
	Sport       domain.Sport `json:"sport"`
	ExternalID  string       `json:"external_id"`
	ScheduledAt string       `json:"scheduled_at"`
	AwayTeam    TeamData     `json:"away_team"`
	HomeTeam    TeamData     `json:"home_team"`
	AwayScore   int          `json:"away_score"`
	HomeScore   int          `json:"home_score"`
	Status      string       `json:"status"`
	Venue       string       `json:"venue"`
}

type TeamData struct {
	Name string `json:"name"`
}

type GamesData struct {
	Games []GameData `json:"games"`
}

func ExportGames(path string, games []domain.Game) error {
	data := GamesData{
		Games: make([]GameData, 0, len(games)),
	}

	for _, game := range games {
		data.Games = append(data.Games, GameData{
			Sport:       game.Sport,
			ExternalID:  game.ExternalID,
			ScheduledAt: game.ScheduledAt.Format("2006-01-02T15:04:05Z07:00"),
			AwayTeam:    TeamData{Name: game.AwayTeam.Name},
			HomeTeam:    TeamData{Name: game.HomeTeam.Name},
			AwayScore:   game.AwayScore,
			HomeScore:   game.HomeScore,
			Status:      game.Status,
			Venue:       game.Venue,
		})
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	return encoder.Encode(data)
}
