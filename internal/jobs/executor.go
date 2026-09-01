package jobs

import (
	"fmt"

	"github.com/JMar2021/sports-data-platform/internal/mlb"
)

type Executor struct {
	Client *mlb.Client
}

func (e *Executor) Execute(job Job) error {
	switch job.Sport {
	case SportMLB:
		switch job.Operation {
		case OperationGetSchedule:
			scheduleResponse, err := e.Client.GetSchedule(job.Date)
			if err != nil {
				return err
			}
			// Process the scheduleResponse as needed
			if len(scheduleResponse.Dates) == 0 {
				fmt.Println("No games scheduled for this date.")
				return nil
			}
			for _, game := range scheduleResponse.Dates[0].Games {
				fmt.Printf("%s %d @ %s %d\n", game.Teams.Away.Team.Name, game.Teams.Away.Score, game.Teams.Home.Team.Name, game.Teams.Home.Score)
			}
			return nil
		default:
			return fmt.Errorf("unsupported operation: %s", job.Operation)
		}
	default:
		return fmt.Errorf("unsupported sport: %s", job.Sport)
	}
}
