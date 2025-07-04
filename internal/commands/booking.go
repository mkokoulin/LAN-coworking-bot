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
			"You can book your first visit to the coworking space for 2 hours for free or a coffee of your choice. Just contact us: <a href='https://t.me/lan_yerevan'>telegram</a>, phone — +37494601303."
	} else if args.Storage.Language == Languages[1].Lang {
		msg.Text =
			"Вы можете забронировать своё первое посещение коворкинга на 2 часа бесплатно или получить кофе на ваш выбор. Просто свяжитесь с нами: <a href='https://t.me/lan_yerevan'>telegram</a>, телефон — +37494601303."
	}

	_, err := bot.Send(msg)
		
	return err
}
