package commands

import (
	"context"
	
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)



func Unknown(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if *args.Language == Languages[0].Lang {
		msg.Text = "I do not know this command 😔 use the /start command"
	} else if *args.Language == Languages[1].Lang {
		msg.Text = "Я не знаю этой команды 😔 воспользуйтесь командой /start"
	} else {
		msg.Text = "I do not know this command 😔 use the /start command"
	}
	
	_, err := bot.Send(msg)
		
	return err
}