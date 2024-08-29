package service

import (
	"TimBishop42/home-assistant-syncer/internal/api/finance"
	"TimBishop42/home-assistant-syncer/internal/api/home"
	"TimBishop42/home-assistant-syncer/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

const (
	budget_status_entity = "sensor.financial_status_raw"
	last_month_spend     = "sensor.last_month_spend_raw"
	current_month_spend  = "sensor.current_month_spend_raw"
)

type FinanceService struct {
	financeClient *finance.Client
	HomeClient    *home.Client
	Config        *config.Config
	logger        *zap.Logger
}

type RequestWrapper struct {
	request *bytes.Buffer
	entity  string
}

func NewService(config *config.Config, logger *zap.Logger) *FinanceService {
	return &FinanceService{
		financeClient: finance.NewFinanceClient(config.FinanceTrackerUrl),
		HomeClient:    home.NewHomeClient(config.HomeAssistantUrl, config),
		Config:        config,
		logger:        logger,
	}
}

func RegisterHooks(lc fx.Lifecycle, s *FinanceService) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go s.run(ctx) // Run the service in a separate goroutine
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Cleanup logic if needed
			return nil
		},
	})
}
func (s *FinanceService) run(ctx context.Context) {
	ticker := time.NewTicker(s.Config.RefreshPeriod)
	defer ticker.Stop()

	tickerContext := context.Background()

	for {
		select {
		case <-ticker.C:
			resp, err := s.financeClient.CallFinanceStore(tickerContext)
			if err != nil {
				s.logger.Error("error calling finance client.",
					zap.Error(err))
				continue
			}

			//Now we need to push 3 entity status's to HA
			allRequests, err := getHomeRequests(resp)

			if err != nil {
				s.logger.Error("error building request for home API",
					zap.Error(err))
				continue
			}

			for _, r := range allRequests {
				homeResp, err := s.HomeClient.UpdateHomeEntityStatus(tickerContext, r.request, r.entity)

				if err != nil {
					s.logger.Error("error calling home API",
						zap.Error(err))
					continue
				}

				s.logger.Info("successfully updated Home entity",
					zap.Int("response_code:", homeResp.StatusCode),
					zap.String("entity", r.entity))
			}

		}
	}
}

func getHomeRequest[T any](input T) (*bytes.Buffer, error) {
	req := home.SimpleRequest{
		Status: input,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return nil, err
	}

	fmt.Println("Json payload: ", string(jsonData))
	body := bytes.NewBuffer(jsonData)

	return body, nil
}

func getHomeRequests(finResp *finance.Response) ([]RequestWrapper, error) {
	requests := make([]RequestWrapper, 3)
	budgetReq, err := getHomeRequest[string](finResp.Status)
	if err != nil {
		return nil, err
	}
	requests[0] = RequestWrapper{
		request: budgetReq,
		entity:  budget_status_entity,
	}

	lastMonthReq, err := getHomeRequest[int](finResp.PriorMonth)
	if err != nil {
		return nil, err
	}
	requests[1] = RequestWrapper{
		request: lastMonthReq,
		entity:  last_month_spend,
	}

	currentMonthReq, err := getHomeRequest[int](finResp.CurrentMonth)
	if err != nil {
		return nil, err
	}
	requests[2] = RequestWrapper{
		request: currentMonthReq,
		entity:  current_month_spend,
	}

	return requests, nil
}
