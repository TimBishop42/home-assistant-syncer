package home

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	apiUrl     string
}

type Request struct {
	Status     string     `json:"status"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	PriorMonthSpend   float32 `json:"prior_month_spend"`
	CurrentMonthSpend float32 `json:"current_month_spend"`
}

func NewHomeClient(url string) *Client {
	return &Client{
		httpClient: &http.Client{},
		apiUrl:     url,
	}
}

func (c *Client) UpdateHomeEntityStatus(ctx context.Context, body *bytes.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiUrl, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	// Handle the response (simplified)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}
