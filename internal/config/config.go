package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mkokoulin/LAN-coworking-bot/internal/services"
)

type Config struct {
	Scope                  string                     `env:"SCOPE" json:"SCOPE"`
	TelegramToken          string                     `env:"TELEGRAM_TOKEN" json:"TELEGRAM_TOKEN"`
	CoworkersSpreadsheetId string                     `env:"COWORKERS_SPREADSHEET_ID" json:"COWORKERS_SPREADSHEET_ID"`
	CoworkersReadRange     string                     `env:"COWORKERS_READ_RANGE" json:"COWORKERS_READ_RANGE"`
	GoogleCloudConfig      services.GoogleCloudConfig `env:"GOOGLE_CLOUD_CONFIG" json:"GOOGLE_CLOUD_CONFIG"`
	GuestWifiPassword      string                     `env:"GUEST_WIFI_PASSWORD" json:"GUEST_WIFI_PASSWORD"`
	CoworkingWifiPassword  string                     `env:"COWORKING_WIFI_PASSWORD" json:"COWORKING_WIFI_PASSWORD"`
	AdminChatId            int64                      `env:"ADMIN_CHAT_ID" json:"ADMIN_CHAT_ID"`
	BotLogsReadRange       string                     `env:"BOT_LOGS_READ_RANGE" json:"BOT_LOGS_READ_RANGE"`
	GuestsReadRange        string                     `env:"GUESTS_READ_RANGE" json:"GUESTS_READ_RANGE"`
	MongoURI               string                     `env:"MONGO_URI" json:"MONGO_URI"`
	MongoDB                string                     `env:"MONGO_DB" json:"MONGO_DB"`
	MongoLocksColl         string                     `env:"MONGO_LOCKS_COLL" json:"MONGO_LOCKS_COLL"`

	AdminUserID   int64 `env:"ADMIN_USER_ID" json:"ADMIN_USER_ID"`
	OrdersChatId  int64 `env:"ORDERS_CHAT_ID" json:"ORDERS_CHAT_ID"`
	OrdersTopicId int   `env:"ORDERS_TOPIC_ID" json:"ORDERS_TOPIC_ID"`

	HaysellBaseURL string   `env:"HAYSELL_BASE_URL" json:"HAYSELL_BASE_URL"`
	HaysellAPIKey string   `env:"HAYSELL_API_KEY" json:"HAYSELL_API_KEY"`
}

func New() (*Config, error) {
	// .env (локально)
	if err := godotenv.Load(); err != nil {
		log.Println("[boot] .env not found (ok if using real env)")
	}

	get := func(key string) (string, error) {
		val := os.Getenv(key)
		if val == "" {
			return "", fmt.Errorf("environment variable %s is not set or empty", key)
		}
		log.Printf("[LAN-COWORKING-BOT] %s: %s", key, val)
		return val, nil
	}

	getOpt := func(key string) string {
		val := os.Getenv(key)
		if val != "" {
			log.Printf("[LAN-COWORKING-BOT] %s: %s", key, val)
		}
		return val
	}

	var cfg Config
	var err error

	if cfg.Scope, err = get("SCOPE"); err != nil { return nil, err }
	if cfg.TelegramToken, err = get("TELEGRAM_TOKEN"); err != nil { return nil, err }
	if cfg.CoworkersSpreadsheetId, err = get("COWORKERS_SPREADSHEET_ID"); err != nil { return nil, err }
	if cfg.CoworkersReadRange, err = get("COWORKERS_READ_RANGE"); err != nil { return nil, err }
	if cfg.GuestWifiPassword, err = get("GUEST_WIFI_PASSWORD"); err != nil { return nil, err }
	if cfg.CoworkingWifiPassword, err = get("COWORKING_WIFI_PASSWORD"); err != nil { return nil, err }
	if cfg.BotLogsReadRange, err = get("BOT_LOGS_READ_RANGE"); err != nil { return nil, err }
	if cfg.GuestsReadRange, err = get("GUESTS_READ_RANGE"); err != nil { return nil, err }
	if cfg.MongoURI, err = get("MONGO_URI"); err != nil { return nil, err }
	if cfg.MongoDB, err = get("MONGO_DB"); err != nil { return nil, err }
	if cfg.MongoLocksColl, err = get("MONGO_LOCKS_COLL"); err != nil { return nil, err }

	if cfg.HaysellBaseURL, err = get("HAYSELL_BASE_URL"); err != nil { return nil, err }
	if cfg.HaysellAPIKey, err = get("HAYSELL_API_KEY"); err != nil { return nil, err }

	// AdminChatId
	if v, err := get("ADMIN_CHAT_ID"); err != nil {
		return nil, err
	} else if n, e := strconv.ParseInt(v, 10, 64); e != nil {
		return nil, fmt.Errorf("invalid ADMIN_CHAT_ID: %w", e)
	} else {
		cfg.AdminChatId = n
	}

	// AdminUserID (опционально)
	if v := getOpt("ADMIN_USER_ID"); v != "" {
		if n, e := strconv.ParseInt(v, 10, 64); e != nil {
			return nil, fmt.Errorf("invalid ADMIN_USER_ID: %w", e)
		} else {
			cfg.AdminUserID = n
		}
	}

	// OrdersChatId (ВНИМАНИЕ: для группы/канала должен быть ОТРИЦАТЕЛЬНЫЙ chat_id вида -100xxxxxxxxxx)
	if v := getOpt("ORDERS_CHAT_ID"); v != "" {
		if n, e := strconv.ParseInt(v, 10, 64); e != nil {
			return nil, fmt.Errorf("invalid ORDERS_CHAT_ID: %w", e)
		} else {
			cfg.OrdersChatId = n
		}
	}

	// OrdersTopicId (если в группе включены темы)
	if v := getOpt("ORDERS_TOPIC_ID"); v != "" {
		if n, e := strconv.Atoi(v); e != nil {
			return nil, fmt.Errorf("invalid ORDERS_TOPIC_ID: %w", e)
		} else {
			cfg.OrdersTopicId = n
		}
	}

	// Google creds
	if gcfg, err := get("GOOGLE_CLOUD_CONFIG"); err != nil {
		return nil, err
	} else if err := json.Unmarshal([]byte(gcfg), &cfg.GoogleCloudConfig); err != nil {
		return nil, fmt.Errorf("failed to parse GOOGLE_CLOUD_CONFIG: %w", err)
	}

	// Санити-чек: чат должен быть отрицательный, иначе бот, как правило, не сможет писать
	if cfg.OrdersChatId > 0 {
		log.Printf("[warn] ORDERS_CHAT_ID выглядит как положительный (%d). Для групп/каналов нужен отрицательный (-100...)", cfg.OrdersChatId)
	}

	return &cfg, nil
}
