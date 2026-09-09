package nfl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("parse date %q: %w", date, err)
	}

	espnDate := parsedDate.Format("20060102")
	url := "https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard?dates=" + espnDate
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
