package commands

import (
	"context"
	// "fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)

func Start(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	args.Storage.IsAuthorized = false

	if args.Storage.Language == Languages[0].Lang {
		msg.Text =
			"The Letters and Numbers space contains:\n" +
			"💻 coworking,\n" +
			"☕️ coffee shop and \n" +
			"✨ event venue.\n\n" +
			"Be sure to check out the /about section — there you will find information about our locations and the rules of conduct in them.\n\n" +
			"Select the command to continue the dialog:\n\n" +
			"commands:\n" +
			"/start – restart\n" +
			"/booking – book your first visit 🎁✨\n" +
			"/wifi – get a password from wifi\n" +
			"/meetingroom – book a meeting\n" +
			"/printout – send documents for printing\n" +
			"/events – information about events\n" +
			"/menu – bar menu 🍷\n" +
			"/about – information about the site and the scheme\n" +
			"/language – changing the interface language\n"
	} else if args.Storage.Language == Languages[1].Lang {
		msg.Text =
			"В пространстве Letters and Numbers размещаются:\n" +
			"💻 коворкинг,\n" +
			"☕️ кофейня и\n" +
			"✨ площадка для мероприятий.\n\n" +
			"Обязательно ознакомьтесь с разделом /about — там вы найдете информацию о наших локациях и правилах поведения в них.\n\n" +
			"Выберите команду про продолжения диалога:\n\n" +
			"команды:\n" +
			"/start – перезапуск\n" +
			"/booking – забронировать своё первое посещение 🎁✨\n" +
			"/wifi – получить пароль от вайфай\n" +
			"/meetingroom – забронировать переговорку\n" +
			"/printout – отправить документы на печать\n" +
			"/events – информация о мероприятиях\n" +
			"/menu – меню бара 🍷\n" +
			"/about – информация о площадке и схема\n" +
			"/language – смена языка интерфейса\n"
	}
	
	args.Storage.CurrentCommand = ""

	_, err := bot.Send(msg)
		
	return err
}
