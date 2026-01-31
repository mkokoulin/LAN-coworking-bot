package types

import (
	"context"
)

type BotLog struct {
	Telegram string `json:"telegram" mapstructure:"telegram"`
	Command  string `json:"command" mapstructure:"command"`
	Datetime string `json:"datetime" mapstructure:"datetime"`
}

type Guest struct {
	Telegram  string `json:"telegram"  mapstructure:"telegram"`
	FirstName string `json:"firstName" mapstructure:"firstName"`
	LastName  string `json:"lastName"  mapstructure:"lastName"`
	Datetime  string `json:"datetime"  mapstructure:"datetime"`
}

type CoworkersSheetsService interface {
    ValidateSecret(ctx context.Context, code string) (bool, error)
    GetUnusedSecrets(ctx context.Context) ([]string, error)
}

type BotLogsSheetsService interface {
	Log(ctx context.Context, rangeName string, log BotLog) error
}

type EventsService interface {
	ListUpcoming(ctx context.Context) ([]Event, error)
}

type SubscriptionsService interface {
	SetWeeklyEvents(ctx context.Context, chatID int64, enabled bool) error
	ListWeeklyEventsSubscribers(ctx context.Context) ([]int64, error)
}

type GuestSheetsService interface {
	GetGuests(ctx context.Context) ([]Guest, error)
	GetGuest(ctx context.Context, telegram string) (Guest, error)

	// историческое имя (оставим для совместимости)
	CreateGuest(ctx context.Context, readRange string, guest Guest) error

	// ожидаемое остальным кодом имя — алиас к CreateGuest
	AddGuest(ctx context.Context, readRange string, guest Guest) error
}

type BarCategory struct {
    ID   string
    Name string
}

type BarProduct struct {
    ID          int
    CategoryID  string
    Name        string
    ShortName   string
    PriceAMD    int
    ImageURL    string
    Description string
    Balance     int
}

type BarCatalogService interface {
    ListCategories(ctx context.Context) ([]BarCategory, error)
    ListProducts(ctx context.Context, ids []int) ([]BarProduct, error)
}
