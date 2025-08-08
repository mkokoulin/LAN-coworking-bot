package botengine

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mkokoulin/LAN-coworking-bot/internal/commands"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/services"
)

type CommandHandler func(
	ctx context.Context,
	update tgbotapi.Update,
	bot *tgbotapi.BotAPI,
	cfg *config.Config,
	services types.Services,
	state *types.ChatStorage,
) error

type CommandRouter struct {
	handlers       map[string]CommandHandler
	defaultHandler CommandHandler
}

func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		handlers: map[string]CommandHandler{
			"start":       commands.StartCommand,
			"wifi":        commands.WifiCommand,
			"language":    commands.LanguageCommand,
			"booking":     commands.BookingCommand,
			"meetingroom": commands.MeetingroomCommand,
			"printout":    commands.PrintoutCommand,
			"events":      commands.EventsCommand,
			"about":       commands.AboutCommand,
			"menu":        commands.MenuCommand,
		},
		defaultHandler: commands.UnknownCommand,
	}
}

func (r *CommandRouter) Handle(
	ctx context.Context,
	update tgbotapi.Update,
	bot *tgbotapi.BotAPI,
	cfg *config.Config,
	services types.Services,
	state *types.ChatStorage,
) error {
	// 1. Если язык не установлен — предлагаем выбрать
	if state.Language == "" {
		state.CurrentCommand = "language"
		_ = commands.LanguageCommand(ctx, update, bot, cfg, services, state)

		if state.Language != "" {
			state.CurrentCommand = "start"
			return commands.StartCommand(ctx, update, bot, cfg, services, state)
		}
		return nil
	}

	// 2. Если команды нет — вызываем unknown
	if state.CurrentCommand == "" {
		return r.defaultHandler(ctx, update, bot, cfg, services, state)
	}

	// 3. Получаем хендлер
	handler, ok := r.handlers[state.CurrentCommand]
	if !ok || handler == nil {
		handler = r.defaultHandler
	}

	// 4. Выполняем хендлер
	err := handler(ctx, update, bot, cfg, services, state)
	if err != nil {
		return err
	}

	// 5. Логируем команду (если есть лог-сервис)
	if services.BotLogsSheets != nil {
		_ = services.BotLogsSheets.Log(ctx, cfg.BotLogsReadRange, services.BotLog{
			Telegram: update.Message.Chat.UserName,
			Command:  state.CurrentCommand,
		})
	}

	return nil
}
