package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)

func Booking(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	msg.ParseMode = "html"
	if args.Storage.Language == Languages[0].Lang {
		msg.Text =
			"You can book your first visit to the coworking space for 2 hours for free and a coffee of your choice. Just contact us: <a href='https://t.me/lan_yerevan'>telegram</a>, <a href='tel:+37494601303'>phone</a>."
	} else if args.Storage.Language == Languages[1].Lang {
		msg.Text = 
			"Вы можете забронировать своё первое посещение коворкинга на 2 часа бесплатно и получить кофе на ваш выбор. Просто свяжитесь с нами: <a href='https://t.me/lan_yerevan'>telegram</a>, <a href='tel:+37494601303'>телефон</a>."
	}

	_, err := bot.Send(msg)
		
	return err
}
