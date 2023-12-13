package commands

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/helpers/stack"
	"github.com/mkokoulin/LAN-coworking-bot/internal/services"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
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
	CoworkersSheets *services.CoworkersSheetService
	BotLogsSheets *services.BotLogsSheetService
	GuestSheets *services.GuestsSheetService
	ChatState *types.ChatStorage
	FirebaseStore *services.Store
	CommandsStack *stack.CommandsStack
}

func CommandsHandler(ctx context.Context, cfg *config.Config, update tgbotapi.Update, bot *tgbotapi.BotAPI, args CommandsHandlerArgs) error {
	if update.Message.IsCommand() {
		args.CommandsStack.Push(update.Message.Command())
	}

	if args.ChatState.Language == "" {
		// args.ChatState.PreviousCommand = args.ChatState.CurrentCommand
		// args.ChatState.CurrentCommand = LANGUAGE

		Language(ctx, update, bot, cfg, args)
	}

	for args.CommandsStack.Top > 0 {

	}
	
	// if update.Message.IsCommand() {
	// 	// args.ChatState.PreviousCommand = args.ChatState.CurrentCommand
	// 	// args.ChatState.CurrentCommand = update.Message.Command()

	// 	// args.

	// 	err := args.FirebaseStore.Update(update.Message.Chat.ID, *args.ChatState)
	// 	if err != nil {
	// 		log.Default().Println("[firebase] failed to update a chat")
	// 	}
	// }

	// if args.ChatState.Language == "" {
	// 	args.ChatState.PreviousCommand = args.ChatState.CurrentCommand
	// 	args.ChatState.CurrentCommand = LANGUAGE

	// 	Language(ctx, update, bot, cfg, args, isFirstTime)
	// }

	// if args.Storage.Language == "" {
	// 	args.Storage.PreviousCommand = args.Storage.CurrentCommand
	// 	args.Storage.CurrentCommand = LANGUAGE

	// 	Language(ctx, update, bot, cfg, args)

	// 	if args.Storage.Language != "" {
	// 		args.Storage.CurrentCommand = args.Storage.PreviousCommand
	// 		args.Storage.PreviousCommand = LANGUAGE
	// 		// return Start(ctx, update, bot, cfg, args)
	// 	}
	// }
	
	// if args.Storage.Language != "" {
	// 	if args.Storage.CurrentCommand != "" {
	// 		if _, ok := commandList[args.Storage.CurrentCommand]; !ok {
	// 			return Unknown(ctx, update, bot, cfg, args)
	// 		}
	
	// 		err := commandList[args.Storage.CurrentCommand](ctx, update, bot, cfg, args)
	// 		if err != nil {
	// 			return err
	// 		}
	
	// 		err = args.BotLogsSheets.Log(ctx, cfg.BotLogsReadRange, services.BotLog{
	// 			Telegram: update.Message.Chat.UserName,
	// 			Command: args.Storage.CurrentCommand,
	// 		})
	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else {
	// 		return Unknown(ctx, update, bot, cfg, args)
	// 	}	
	// }

	return nil
}