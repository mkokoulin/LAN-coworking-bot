package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)

func Events(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	msg.ParseMode = "html"
	
	if args.Storage.Language == Languages[0].Lang {
		msg.Text = "We have a large number of different events, we publish announcements of events on our social networks: <a href='https://www.instagram.com/lan_yerevan/'>Instagram</a> and <a href='https://t.me/lan_yerevan'>Telegram</a>. Subscribe to keep up to date with cool events. An up-to-date list of events and reservations is maintained via <a href='https://lettersandnumbers.am/events?eventId=events'>taplink</a>"
	} else if args.Storage.Language == Languages[1].Lang {
		msg.Text = "У нас проходит большое количество разнообразных мероприятий, анонсы событий мы публикуем в наших социальных сетях: <a href='https://www.instagram.com/lan_yerevan/'>Instagram</a> и <a href='https://t.me/lan_yerevan'>Telegram</a>. Подписывайтесь, чтобы быть в курсе классных событий 🎉. Актуальный список мероприятий и бронирование ведется через <a href='https://lettersandnumbers.am/events?eventId=events'>taplink</a>"
	}

	args.Storage.CurrentCommand = ""

	_, err := bot.Send(msg)
		
	return err
}