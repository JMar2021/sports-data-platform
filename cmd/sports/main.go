package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JMar2021/sports-data-platform/internal/mlb"
)

func main() {
	httpClient := &http.Client{}
	client := mlb.NewClient(httpClient)

	scheduleResponse, err := client.GetSchedule(time.Now().Format("2006-01-02"))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, game := range scheduleResponse.Dates[0].Games {
		fmt.Printf("%s %d @ %s %d\n", game.Teams.Away.Team.Name, game.Teams.Away.Score, game.Teams.Home.Team.Name, game.Teams.Home.Score)
	}
}
