package mlb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	HTTPClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		HTTPClient: httpClient,
	}
}

func (c *Client) GetSchedule(ctx context.Context, date string) (*ScheduleResponse, error) {
	url := "https://statsapi.mlb.com/api/v1/schedule?date=" + date + "&sportId=1"
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 300 || response.StatusCode < 200 {
		return nil, fmt.Errorf("error response code: %d", response.StatusCode)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var scheduleResponse ScheduleResponse
	err = json.Unmarshal(body, &scheduleResponse)
	if err != nil {
		return nil, err
	}
	return &scheduleResponse, nil
}

func (c *Client) GetStandings(ctx context.Context, season string) (*StandingsResponse, error) {
	url := fmt.Sprintf(
		"https://statsapi.mlb.com/api/v1/standings?leagueId=103,104&season=%s&standingsTypes=regularSeason",
		season,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MLB API returned status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var standingsResponse StandingsResponse

	if err := json.Unmarshal(body, &standingsResponse); err != nil {
		return nil, err
	}

	return &standingsResponse, nil
}
