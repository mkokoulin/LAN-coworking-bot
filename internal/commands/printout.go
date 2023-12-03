package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)

func Printout(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	msg.Text = "Отправьте документы для распечатки в аккаунт @lan_yerevan (администратору) и уточните у него стоимость услуги"
	
	_, err := bot.Send(msg)
		
	return err
}