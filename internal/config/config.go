package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/mkokoulin/LAN-coworking-bot/internal/services"
)

type Config struct {
	Scope string `env:"SCOPE" json:"SCOPE"`
	TelegramToken string `env:"TELEGRAM_TOKEN" json:"TELEGRAM_TOKEN"`
	CoworkersSpreadsheetId string `env:"COWORKERS_SPREADSHEET_ID" json:"COWORKERS_SPREADSHEET_ID"`
	CoworkersReadRange string `env:"COWORKERS_READ_RANGE" json:"COWORKERS_READ_RANGE"`
	GoogleCloudConfig services.GoogleCloudConfig `env:"GOOGLE_CLOUD_CONFIG" json:"GOOGLE_CLOUD_CONFIG"`
	GuestWifiPassword string `env:"GUEST_WIFI_PASSWORD" json:"GUEST_WIFI_PASSWORD"`
	CoworkingWifiPassword string `env:"COWORKING_WIFI_PASSWORD" json:"COWORKING_WIFI_PASSWORD"`
	AdminChatId int64 `env:"ADMIN_CHAT_ID" json:"ADMIN_CHAT_ID"`
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

	cfg.GuestWifiPassword = os.Getenv("GUEST_WIFI_PASSWORD")
	if cfg.GuestWifiPassword == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "GUEST_WIFI_PASSWORD")
	}
	log.Default().Printf("[LAN-COWORKING-BOT] GUEST_WIFI_PASSWORD: %v", cfg.GuestWifiPassword)

	cfg.CoworkingWifiPassword = os.Getenv("COWORKING_WIFI_PASSWORD")
	if cfg.CoworkingWifiPassword == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "COWORKING_WIFI_PASSWORD")
	}
	log.Default().Printf("[LAN-COWORKING-BOT] COWORKING_WIFI_PASSWORD: %v", cfg.CoworkingWifiPassword)

	adminChatIdString := os.Getenv("ADMIN_CHAT_ID")
	if adminChatIdString == "" {
		return nil, fmt.Errorf("environment variable %v is not set or empty", "ADMIN_CHAT_ID")
	}

	i, err := strconv.Atoi(adminChatIdString)
    if err != nil {
        return nil, fmt.Errorf("error parsing int: %v", err)
    }

	cfg.AdminChatId = int64(i)

	log.Default().Printf("[LAN-COWORKING-BOT] ADMIN_CHAT_ID: %v", cfg.AdminChatId)

	return &cfg, nil
}