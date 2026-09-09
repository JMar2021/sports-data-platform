package domain

import "time"

type Game struct {
	Sport       Sport
	ExternalID  string
	ScheduledAt time.Time
	HomeTeam    Team
	AwayTeam    Team
	HomeScore   int
	AwayScore   int
	Status      string
	Venue       string
}
