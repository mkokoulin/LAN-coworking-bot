package main

import (
	"context"
	"log"
	"net/http"

	"github.com/mkokoulin/LAN-coworking-bot/internal/commands"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	START       = "start"
	WIFI        = "wifi"
	MEETINGROOM = "meetingroom"
	PRINTOUT    = "printout"
	EVENTS      = "events"
	ABOUT       = "about"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatalln(err)
		return
	}

	gc, err := services.NewGoogleClient(ctx, cfg.GoogleCloudConfig, cfg.Scope)
	if err != nil {
		log.Fatalf("fatal error %v", err)
	}

	coworkersSheets, err := services.NewCoworkersSheets(ctx, gc, cfg.CoworkersSpreadsheetId, cfg.CoworkersReadRange)
	if err != nil {
		log.Fatalf("fatal error %v", err)
	}

	botLogsSheets, err := services.NewBotLogsSheets(ctx, gc, cfg.CoworkersSpreadsheetId, cfg.BotLogsReadRange)
	if err != nil {
		log.Fatalf("fatal error %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("fatal error %v", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	var currentCommand string
	var isAwaitingConfirmation bool
	var isAuthorized bool
	var language string
	var isBookingProcess bool
	var isWifiProcess bool

	go func() {
		_ = http.ListenAndServe(":8080", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("ok"))
			},
		))
	}()

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				currentCommand = update.Message.Command()
			}

			err := commands.CommandsHandler(ctx, cfg, update, bot, commands.CommandsHandlerArgs{
				Language:               &language,
				CurrentCommand:         &currentCommand,
				IsBookingProcess:       &isBookingProcess,
				IsAwaitingConfirmation: &isAwaitingConfirmation,
				IsAuthorized:           &isAuthorized,
				CoworkersSheets:        coworkersSheets,
				BotLogsSheets:		    botLogsSheets,
				IsWifiProcess:          &isWifiProcess,
			})
			if err != nil {
				log.Fatalf("fatal error %v", err)
			}
		}
	}
}
