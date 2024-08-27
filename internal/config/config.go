package config

import (
	"os"
	"time"
)

type Config struct {
	FinanceTrackerUrl string
	HomeAssistantUrl  string
	RefreshPeriod     time.Duration
	HomeKey           string
}

func NewConfig() *Config {
	return &Config{
		FinanceTrackerUrl: getEnvString("FINANCE_URL", "http://localhost:8080/api/finance/get-home-data"),
		HomeAssistantUrl:  getEnvString("HOME_URL", "http://192.168.0.73:8123/api/states/finance.data"),
		RefreshPeriod:     getEnvDuration("REFRESH_PERIOD", time.Second*10),
		HomeKey:           getEnvString("HOME_KEY", ""),
	}
}

func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		duration, _ := time.ParseDuration(value)
		return duration
	}
	return defaultVal
}
