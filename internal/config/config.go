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
	Scope                  string                         `env:"SCOPE" json:"SCOPE"`
	TelegramToken          string                         `env:"TELEGRAM_TOKEN" json:"TELEGRAM_TOKEN"`
	CoworkersSpreadsheetId string                         `env:"COWORKERS_SPREADSHEET_ID" json:"COWORKERS_SPREADSHEET_ID"`
	CoworkersReadRange     string                         `env:"COWORKERS_READ_RANGE" json:"COWORKERS_READ_RANGE"`
	GoogleCloudConfig      services.GoogleCloudConfig     `env:"GOOGLE_CLOUD_CONFIG" json:"GOOGLE_CLOUD_CONFIG"`
	GuestWifiPassword      string                         `env:"GUEST_WIFI_PASSWORD" json:"GUEST_WIFI_PASSWORD"`
	CoworkingWifiPassword  string                         `env:"COWORKING_WIFI_PASSWORD" json:"COWORKING_WIFI_PASSWORD"`
	AdminChatId            int64                          `env:"ADMIN_CHAT_ID" json:"ADMIN_CHAT_ID"`
	BotLogsReadRange       string                         `env:"BOT_LOGS_READ_RANGE" json:"BOT_LOGS_READ_RANGE"`
	GuestsReadRange        string                         `env:"GUESTS_READ_RANGE" json:"GUESTS_READ_RANGE"`
	MongoURI               string                         `env:"MONGO_URI" json:"MONGO_URI"`
}

func New() (*Config, error) {
	get := func(key string) (string, error) {
		val := os.Getenv(key)
		if val == "" {
			return "", fmt.Errorf("environment variable %s is not set or empty", key)
		}
		log.Printf("[LAN-COWORKING-BOT] %s: %s", key, val)
		return val, nil
	}

	cfg := Config{}

	var err error

	if cfg.Scope, err = get("SCOPE"); err != nil {
		return nil, err
	}
	if cfg.TelegramToken, err = get("TELEGRAM_TOKEN"); err != nil {
		return nil, err
	}
	if cfg.CoworkersSpreadsheetId, err = get("COWORKERS_SPREADSHEET_ID"); err != nil {
		return nil, err
	}
	if cfg.CoworkersReadRange, err = get("COWORKERS_READ_RANGE"); err != nil {
		return nil, err
	}
	if cfg.GuestWifiPassword, err = get("GUEST_WIFI_PASSWORD"); err != nil {
		return nil, err
	}
	if cfg.CoworkingWifiPassword, err = get("COWORKING_WIFI_PASSWORD"); err != nil {
		return nil, err
	}
	if cfg.BotLogsReadRange, err = get("BOT_LOGS_READ_RANGE"); err != nil {
		return nil, err
	}
	if cfg.GuestsReadRange, err = get("GUESTS_READ_RANGE"); err != nil {
		return nil, err
	}
	if cfg.MongoURI, err = get("MONGO_URI"); err != nil {
		return nil, err
	}

	adminID, err := get("ADMIN_CHAT_ID")
	if err != nil {
		return nil, err
	}
	numID, err := strconv.ParseInt(adminID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid ADMIN_CHAT_ID: %w", err)
	}
	cfg.AdminChatId = numID

	gcfg, err := get("GOOGLE_CLOUD_CONFIG")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(gcfg), &cfg.GoogleCloudConfig); err != nil {
		return nil, fmt.Errorf("failed to parse GOOGLE_CLOUD_CONFIG: %w", err)
	}

	return &cfg, nil
}
