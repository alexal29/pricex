package handlers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

type CoinInfo struct {
	ID         string `json:"id,omitempty"`
	Symbol     string `json:"symbol,omitempty"`
	Name       string `json:"name,omitempty"`
	MetricDesc *prometheus.Desc
}

type Config struct {
	Coins      []*CoinInfo `json:"coins,omitempty"`
	Currencies string      `json:"currencies,omitempty"`
}

func LoadConfig() (*Config, error) {
	b, err := os.ReadFile("./config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read file config.json: %v", err)
	}

	var cfg Config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse coins info")
	}

	return &cfg, nil
}
