package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ClientID string `json:"client_id"`
	RedirectURI string `json:"redirect_uri"`
	Scope string `json:"scope"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("config/config.json")

	if err != nil {
		return nil, fmt.Errorf("Failed to read config file: %w", err)
	}

	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("Failed to parse config file: %w", err)
	}

	return &config, nil
}
