package commands

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)

func Meetingroom(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if !*args.IsBookingProcess {	
		if *args.Language == Languages[0].Lang {
			msg.Text = "Write the date and time interval for which you want to book a meeting room in the format yyyy-mm-dd hh:mm - hh:mm"
		} else if *args.Language == Languages[1].Lang {
			msg.Text = "Напишите дату и интервал времени, на который вы хотите забронировать комнату для переговоров в формате yyyy-mm-dd hh:mm - hh:mm"
		}
		
		*args.IsBookingProcess = true

		_, err := bot.Send(msg)
			
		return err
	} else {
		if update.Message.Text == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.Text = "Сообщение не может быть пустым"
			bot.Send(msg)
			
			return nil
		}
	
		msgToAdmin := tgbotapi.NewMessage(cfg.AdminChatId, fmt.Sprintf("Пользователь @%s просит забронировать переговорку - %s", update.Message.Chat.UserName, update.Message.Text))
		bot.Send(msgToAdmin)

		if *args.Language == Languages[0].Lang {
			msg.Text = "Our administrator will contact you soon 🧑‍💼"
		} else if *args.Language == Languages[1].Lang {
			msg.Text = "В ближайшее время с вами свяжется наш администратор 🧑‍💼"
		}
		
		bot.Send(msg)

		*args.IsBookingProcess = false
		*args.CurrentCommand = ""
	}

	return nil
}