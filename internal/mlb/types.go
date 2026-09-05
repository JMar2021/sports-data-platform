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

type Venue struct {
	Name string `json:"name"`
}

type GameStatus struct {
	AbstractGameState string `json:"abstractGameState"`
}

type Game struct {
	GamePk   int        `json:"gamePk"`
	GameDate string     `json:"gameDate"`
	Status   GameStatus `json:"status"`
	Teams    GameTeams  `json:"teams"`
	Venue    Venue      `json:"venue"`
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

type StandingsTeam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LeagueRecord struct {
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
	Pct    string `json:"pct"`
}

type StandingsRecord struct {
	Team              StandingsTeam `json:"team"`
	Season            string        `json:"season"`
	DivisionRank      string        `json:"divisionRank"`
	LeagueRank        string        `json:"leagueRank"`
	GamesPlayed       int           `json:"gamesPlayed"`
	GamesBack         string        `json:"gamesBack"`
	LeagueRecord      LeagueRecord  `json:"leagueRecord"`
	RunsScored        int           `json:"runsScored"`
	RunsAllowed       int           `json:"runsAllowed"`
	RunDifferential   int           `json:"runDifferential"`
	Wins              int           `json:"wins"`
	Losses            int           `json:"losses"`
	WinningPercentage string        `json:"winningPercentage"`
	LastUpdated       string        `json:"lastUpdated"`
}

type StandingsGroup struct {
	League      any               `json:"league"`
	Division    any               `json:"division"`
	TeamRecords []StandingsRecord `json:"teamRecords"`
}

type StandingsResponse struct {
	Records []StandingsGroup `json:"records"`
}
