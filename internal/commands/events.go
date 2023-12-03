package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)

func Events(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	msg.ParseMode = "html"
	msg.Text = "У нас проходит большое количество разнообразных мероприятий, анонсы событий мы публикуем в наших социальных сетях: <a href='https://www.instagram.com/lan_yerevan/'>Instagram</a> и <a href='https://t.me/lan_yerevan'>Telegram</a>. Подписывайтесь, чтобы быть в курсе классных событий 🎉. Актуальный список мероприятий и бронирование ведется через <a href='https://taplink.cc/lan_yerevan'>taplink</a>"
	
	_, err := bot.Send(msg)
		
	return err
}