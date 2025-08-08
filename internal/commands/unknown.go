package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/locales"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

func UnknownCommand(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, services types.Services, state *types.ChatStorage) error {
	p := locales.Printer(state.Language)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, p.Sprintf("unknown_command"))
	_, err := bot.Send(msg)
	return err
}
