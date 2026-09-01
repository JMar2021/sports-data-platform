package mlb

type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GameTeam struct {
	Team  Team `json:"team"`
	Score int  `json:"score"`
}

type GameTeams struct {
	Away GameTeam `json:"away"`
	Home GameTeam `json:"home"`
}

type GameStatus struct {
	AbstractGameState string `json:"abstractGameState"`
}

type Game struct {
	GamePk   int        `json:"gamePk"`
	GameDate string     `json:"gameDate"`
	Status   GameStatus `json:"status"`
	Teams    GameTeams  `json:"teams"`
}

type ScheduleDate struct {
	Date       string `json:"date"`
	TotalGames int    `json:"totalGames"`
	Games      []Game `json:"games"`
}

type ScheduleResponse struct {
	Dates      []ScheduleDate `json:"dates"`
	TotalGames int            `json:"totalGames"`
}
