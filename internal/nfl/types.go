package nfl

type Team struct {
	ID   string `json:"id"`
	Name string `json:"displayName"`
}

type Venue struct {
	Name string `json:"fullName"`
}

type ScheduleResponse struct {
	Events []Event `json:"events"`
}

type Event struct {
	ID           string        `json:"id"`
	Date         string        `json:"date"`
	Name         string        `json:"name"`
	Season       Season        `json:"season"`
	Week         Week          `json:"week"`
	Competitions []Competition `json:"competitions"`
}
type Competition struct {
	ID          string       `json:"id"`
	Date        string       `json:"date"`
	Venue       Venue        `json:"venue"`
	Competitors []Competitor `json:"competitors"`
	Status      Status       `json:"status"`
}
type Competitor struct {
	ID       string `json:"id"`
	Team     Team   `json:"team"`
	Score    string `json:"score"`
	HomeAway string `json:"homeAway"`
}
type Season struct {
	Year int    `json:"year"`
	Slug string `json:"slug"`
}
type Week struct {
	Number int `json:"number"`
}
type Status struct {
	Clock float32 `json:"clock"`
	Type  Type    `json:"type"`
}
type Type struct {
	Completed bool   `json:"completed"`
	Name      string `json:"name"`
	State     string `json:"state"`
}
