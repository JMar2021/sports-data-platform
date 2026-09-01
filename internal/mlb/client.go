package mlb

import (
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

func (c *Client) GetSchedule(date string) (*ScheduleResponse, error) {
	url := "https://statsapi.mlb.com/api/v1/schedule?date=" + date + "&sportId=1"
	request, err := http.NewRequest("GET", url, nil)
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
