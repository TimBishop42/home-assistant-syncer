package finance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	apiUrl     string
}

type Response struct {
	Status       string `json:"status"`
	CurrentMonth int    `json:"currentMonth"`
	PriorMonth   int    `json:"priorMonth"`
}

func NewFinanceClient(url string) *Client {
	return &Client{
		httpClient: &http.Client{},
		apiUrl:     url,
	}
}

func (c *Client) CallFinanceStore(ctx context.Context) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiUrl, nil)
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

	var financeResponse Response
	if err := json.NewDecoder(resp.Body).Decode(&financeResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &financeResponse, nil
}
