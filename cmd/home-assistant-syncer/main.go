package main

import (
	"TimBishop42/home-assistant-syncer/internal/config"
	"TimBishop42/home-assistant-syncer/internal/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.NewConfig,
			service.NewService,
			NewLogger,
		),
		fx.Invoke(service.RegisterHooks),
	)
	app.Run()
}

func NewLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
