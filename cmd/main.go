package main

import (
	"context"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/locales"
	"github.com/mkokoulin/LAN-coworking-bot/internal/services"
	"github.com/mkokoulin/LAN-coworking-bot/internal/state"
)

func main() {
	locales.Init()

	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	googleClient, err := services.NewGoogleClient(ctx, cfg.GoogleCloudConfig, cfg.Scope)
	if err != nil {
		log.Fatalf("failed to init google client: %v", err)
	}

	coworkersSheets, err := services.NewCoworkersSheets(ctx, googleClient, cfg.CoworkersSpreadsheetId, cfg.CoworkersReadRange)
	if err != nil {
		log.Fatalf("failed to init coworkers sheet: %v", err)
	}

	guestSheets, err := services.NewGuestSheets(ctx, googleClient, cfg.CoworkersSpreadsheetId, cfg.GuestsReadRange)
	if err != nil {
		log.Fatalf("failed to init guest sheet: %v", err)
	}

	botLogsSheets, err := services.NewBotLogsSheets(ctx, googleClient, cfg.CoworkersSpreadsheetId, cfg.BotLogsReadRange)
	if err != nil {
		log.Fatalf("failed to init bot logs sheet: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("failed to init telegram bot: %v", err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	db := mongoClient.Database("lan_bot")

	stateManager, err := state.NewMongo(cfg.MongoURI, "lan_bot", "user_states")
	if err != nil {
		log.Fatalf("failed to init mongo state manager: %v", err)
	}

	commandLogger := state.NewCommandLogger(db, "command_history")

	services := botengine.Services{
		CoworkersSheets: coworkersSheets,
		GuestSheets:     guestSheets,
		BotLogsSheets:   botLogsSheets,
	}
	dispatcher := botengine.NewDispatcher(bot, cfg, services, stateManager)
	dispatcher.AttachLogger(commandLogger)

	go func() {
		_ = http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("ok"))
		}))
	}()

	dispatcher.Run(ctx)
}
