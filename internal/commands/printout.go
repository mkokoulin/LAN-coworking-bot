package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)

func Printout(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if *args.Language == Languages[0].Lang {
		msg.Text = "Send the documents for printing to the account @lan_yerevan (administrator) and check with him the cost of the service"
	} else if *args.Language == Languages[1].Lang {
		msg.Text = "Отправьте документы для распечатки в аккаунт @lan_yerevan (администратору) и уточните у него стоимость услуги"
	}

	_, err := bot.Send(msg)
		
	return err
}