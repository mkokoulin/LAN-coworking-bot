package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
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
	MENU = "menu"
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
	MENU: Menu,
}

type CommandsHandlerArgs struct {
	CoworkersSheets *services.CoworkersSheetService
	BotLogsSheets *services.BotLogsSheetService
	Storage *types.ChatStorage
	GuestSheets *services.GuestsSheetService
}

func CommandsHandler(ctx context.Context, cfg *config.Config, update tgbotapi.Update, bot *tgbotapi.BotAPI, args CommandsHandlerArgs) error {
	if args.Storage.Language == "" {
		args.Storage.CurrentCommand = LANGUAGE

		Language(ctx, update, bot, cfg, args)

		if args.Storage.Language != "" {
			args.Storage.CurrentCommand = START
			return Start(ctx, update, bot, cfg, args)
		}
		
		return nil
	}
	
	if args.Storage.CurrentCommand != "" {
		if _, ok := commandList[args.Storage.CurrentCommand]; !ok {
			return Unknown(ctx, update, bot, cfg, args)
		}

		err := commandList[args.Storage.CurrentCommand](ctx, update, bot, cfg, args)
		if err != nil {
			return err
		}

		err = args.BotLogsSheets.Log(ctx, cfg.BotLogsReadRange, services.BotLog{
			Telegram: update.Message.Chat.UserName,
			Command: args.Storage.CurrentCommand,
		})
		if err != nil {
			return err
		}
	} else {
		return Unknown(ctx, update, bot, cfg, args)
	}		

	return nil
}
