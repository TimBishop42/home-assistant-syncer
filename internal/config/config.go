package config

import "time"

type Config struct {
	FinanceTrackerUrl string
	HomeAssistantUrl  string
	RefreshPeriod     time.Duration
}

func NewConfig() *Config {
	return &Config{
		FinanceTrackerUrl: "http://192.168.0.67:8080/api/finance/get-home-data",
		HomeAssistantUrl:  "http://192.168.0.73:8123/api/states/finance.data",
		RefreshPeriod:     time.Second * 10,
	}
}
