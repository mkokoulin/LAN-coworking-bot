package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/services"
)

const (
	START = "start"
	WIFI = "wifi"
	MEETINGROOM = "meetingroom"
	PRINTOUT = "printout"
	EVENTS = "events"
	ABOUT = "about"
	LANGUAGE = "language"
)

type Command func (ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error

var commandList = map[string]Command {
	START: Start,
	WIFI: Wifi,
	MEETINGROOM: Meetingroom,
	PRINTOUT: Printout,
	EVENTS: Events,
	ABOUT: About,
	LANGUAGE: Language,
}

type CommandsHandlerArgs struct {
	Language *string
	CurrentCommand *string
	IsBookingProcess *bool
	IsAwaitingConfirmation *bool
	IsAuthorized *bool
	CoworkersSheets *services.CoworkersSheetService
	BotLogsSheets *services.BotLogsSheetService
	IsWifiProcess *bool
	GuestSheets *services.GuestsSheetService
}

func CommandsHandler(ctx context.Context, cfg *config.Config, update tgbotapi.Update, bot *tgbotapi.BotAPI, args CommandsHandlerArgs) error {
	if *args.CurrentCommand != "" {
		if _, ok := commandList[*args.CurrentCommand]; !ok {
			return Unknown(ctx, update, bot, cfg, args)
		}

		err := commandList[*args.CurrentCommand](ctx, update, bot, cfg, args)
		if err != nil {
			return err
		}

		err = args.BotLogsSheets.Log(ctx, cfg.BotLogsReadRange, services.BotLog{
			Telegram: update.Message.Chat.UserName,
			Command: *args.CurrentCommand,
		})
		if err != nil {
			return err
		}
	} else {
		return Unknown(ctx, update, bot, cfg, args)
	}		

	return nil
}