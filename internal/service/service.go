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

type FinanceService struct {
	financeClient *finance.Client
	HomeClient    *home.Client
	Config        *config.Config
	logger        *zap.Logger
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

	for {
		select {
		case <-ticker.C:
			resp, err := s.financeClient.CallFinanceStore(ctx)
			if err != nil {
				s.logger.Error("error calling finance client.",
					zap.Error(err))
				continue
			}

			homeRequest, err := getHomeRequest(resp)

			if err != nil {
				s.logger.Error("error building request for home API",
					zap.Error(err))
				continue
			}

			homeResp, err := s.HomeClient.UpdateHomeEntityStatus(ctx, homeRequest)

			if err != nil {
				s.logger.Error("error calling home API",
					zap.Error(err))
				continue
			}

			s.logger.Info("successfully updated Home entity",
				zap.Any("response", homeResp))
		case <-ctx.Done():
			return
		}
	}
}

func getHomeRequest(resp *finance.Response) (*bytes.Reader, error) {
	req := home.Request{
		Status: resp.Status,
		Attributes: home.Attributes{
			PriorMonthSpend:   resp.PriorMonth,
			CurrentMonthSpend: resp.CurrentMonth,
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return nil, err
	}

	body := bytes.NewReader(jsonData)

	return body, nil
}
