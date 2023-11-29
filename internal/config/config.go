package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mkokoulin/LAN-coworking-bot/internal/services"
)

type Config struct {
	Scope string `env:"SCOPE" json:"SCOPE"`
	TelegramToken string `env:"TELEGRAM_TOKEN" json:"TELEGRAM_TOKEN"`
	CoworkersSpreadsheetId string `env:"COWORKERS_SPREADSHEET_ID" json:"COWORKERS_SPREADSHEET_ID"`
	CoworkersReadRange string `env:"COWORKERS_READ_RANGE" json:"COWORKERS_READ_RANGE"`
	GoogleCloudConfig services.GoogleCloudConfig `env:"GOOGLE_CLOUD_CONFIG" json:"GOOGLE_CLOUD_CONFIG"`
}

func New() (*Config, error) {
	cfg := Config{}

	cfg.Scope = os.Getenv("SCOPE")
	if cfg.Scope == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "SCOPE")
	}
	log.Default().Printf("[LAN-COWORKING-BOT] SCOPE: %v", cfg.Scope)

	cfg.TelegramToken = os.Getenv("TELEGRAM_TOKEN")
	if cfg.TelegramToken == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "TELEGRAM_TOKEN")
	}
	log.Default().Printf("[LAN-COWORKING-BOT] TELEGRAM_TOKEN: %v", cfg.TelegramToken)

	cfg.CoworkersSpreadsheetId = os.Getenv("COWORKERS_SPREADSHEET_ID")
	if cfg.CoworkersSpreadsheetId == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "COWORKERS_SPREADSHEET_ID")
	}
	log.Default().Printf("[LAN-COWORKING-BOT] COWORKERS_SPREADSHEET_ID: %v", cfg.CoworkersSpreadsheetId)

	cfg.CoworkersReadRange = os.Getenv("COWORKERS_READ_RANGE")
	if cfg.CoworkersReadRange == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "COWORKERS_READ_RANGE")
	}
	log.Default().Printf("[LAN-COWORKING-BOT] COWORKERS_READ_RANGE: %v", cfg.CoworkersReadRange)

	googleCloudConfigString := os.Getenv("GOOGLE_CLOUD_CONFIG")
	if googleCloudConfigString == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "GOOGLE_CLOUD_CONFIG")
	}

	var googleCloudConfig services.GoogleCloudConfig
	if err := json.Unmarshal([]byte(googleCloudConfigString), &googleCloudConfig); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}
	log.Default().Printf("[LAN-TG-BOT] GOOGLE_CLOUD_CONFIG: %v", cfg.GoogleCloudConfig)
	cfg.GoogleCloudConfig = googleCloudConfig;

	return &cfg, nil
}