package main

import (
	"context"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mkokoulin/LAN-coworking-bot/internal/commands"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/services"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

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

	guestSheets, err := services.NewGuestSheets(ctx, gc, cfg.CoworkersSpreadsheetId, cfg.GuestsReadRange)
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

	var storage = map[int64] *types.ChatStorage {}

	go func() {
		_ = http.ListenAndServe(":8080", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("ok"))
			},
		))
	}()

	for update := range updates {
		_, ok := storage[update.Message.Chat.ID]

		if !ok {
			storage[update.Message.Chat.ID] = &types.ChatStorage{}
		}

		s := storage[update.Message.Chat.ID]

		if update.Message != nil {
			go func() {
				guest := services.Guest{}
				
				guest.FirstName = update.Message.Chat.FirstName
				guest.LastName = update.Message.Chat.LastName
				guest.Telegram = update.Message.Chat.UserName
	
				err := guestSheets.CreateGuest(ctx, cfg.GuestsReadRange, guest)
				if err != nil {
					log.Default().Println("failed to save the guest")
				}
			}()

			if update.Message.IsCommand() {
				s.IsAwaitingConfirmation = false
				s.IsBookingProcess = false
				s.IsWifiProcess = false

				s.CurrentCommand = update.Message.Command()
			}

			err := commands.CommandsHandler(ctx, cfg, update, bot, commands.CommandsHandlerArgs{
				CoworkersSheets:        coworkersSheets,
				BotLogsSheets:		    botLogsSheets,
				GuestSheets: 			guestSheets,
				Storage: 				s,
			})
			if err != nil {
				log.Fatalf("fatal error %v", err)
			}
		}
	}
}
