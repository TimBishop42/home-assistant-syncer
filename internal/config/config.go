package config

import (
	"go.uber.org/zap"
	"os"
	"time"
)

type Config struct {
	FinanceTrackerUrl string
	HomeAssistantUrl  string
	RefreshPeriod     time.Duration
	HomeKey           string
	logger            *zap.Logger
}

func NewConfig(logger *zap.Logger) *Config {
	return &Config{
		FinanceTrackerUrl: getEnvString("FINANCE_URL", "http://localhost:8080/api/finance/get-home-data"),
		HomeAssistantUrl:  getEnvString("HOME_URL", "http://192.168.0.73:8123/api/states"),
		RefreshPeriod:     getEnvDuration("REFRESH_PERIOD", time.Second*10, logger),
		HomeKey:           getEnvString("HOME_KEY", ""),
		logger:            logger,
	}
}

func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvDuration(key string, defaultVal time.Duration, logger *zap.Logger) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		duration, err := time.ParseDuration(value)
		if err != nil {
			logger.Error("Error unmarshaling duration: ", zap.Error(err))
			duration = time.Second * 60
		}
		logger.Info("Duration is: ", zap.Any("duration", duration))
		return duration
	}
	return defaultVal
}
