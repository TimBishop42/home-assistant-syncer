package home

import (
	"TimBishop42/home-assistant-syncer/internal/config"
	"bytes"
	"context"
	"fmt"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	apiUrl     string
	config     *config.Config
}

type Request struct {
	Status     string     `json:"state"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	PriorMonthSpend   float32 `json:"last_month_spend"`
	CurrentMonthSpend float32 `json:"current_month_spend"`
}

type SimpleRequest struct {
	Status any `json:"state"`
}

func NewHomeClient(url string, config *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{},
		apiUrl:     url,
		config:     config,
	}
}

func (c *Client) UpdateHomeEntityStatus(ctx context.Context, body *bytes.Buffer, entity string) (*http.Response,
	error) {
	urlWithEntity := fmt.Sprintf("%s/%s", c.apiUrl, entity)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlWithEntity, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	//Add bearer token header
	req.Header.Set("Authorization", "Bearer "+c.config.HomeKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	// Handle the response (simplified)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, err: %v", resp.StatusCode, resp.Status)
	}

	return resp, nil
}
