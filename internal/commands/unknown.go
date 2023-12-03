package commands

import (
	"context"
	
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)



func Unknown(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	msg.Text = "Я не знаю этой команды 😔"
	
	_, err := bot.Send(msg)
		
	return err
}