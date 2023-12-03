package commands

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)

func Start(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if *args.Language == "" {
		if update.Message.Text != "" {
			for _, v := range Languages {
				if update.Message.Text == v.Lang {
					*args.Language = v.Lang

					if v.Lang == Languages[0].Lang {
						msg.Text = fmt.Sprintf("Selected %s", v.Lang)
					} else if v.Lang == Languages[1].Lang {
						msg.Text = fmt.Sprintf("Выбран %s", v.Lang)
					}

					msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
						RemoveKeyboard: true,
						Selective: false,
					}
				}
			}
		} else {
			return nil
		}

		if *args.Language == "" {
			Language(ctx, update, bot, cfg, args)
			return nil
		}
	}

	if *args.Language == Languages[0].Lang {
		msg.Text =
			"В пространстве Letters and Numbers размещаются:\n" +
			"💻 коворкинг,\n" +
			"☕️ кофейня и\n" +
			"✨ площадка для мероприятий.\n\n" +
			"Обязательно ознакомьтесь с разделом /about — там вы найдете информацию о наших локациях и правилах поведения в них.\n\n" +
			"Выберите команду про продолжения диалога:\n\n" +
			"команды:\n" +
			"/start – перезапуск\n" +
			"/wifi – получить пароль от вайфай\n" +
			"/meetingroom – забронировать переговорку\n" +
			"/printout – отправить документы на печать\n" +
			"/events – информация о мероприятиях\n" +
			"/about – информация о площадке и схема\n" +
			"/language – смена языка интерфейса\n"
	} else if *args.Language == Languages[1].Lang {
		msg.Text =
			"В пространстве Letters and Numbers размещаются:\n" +
			"💻 коворкинг,\n" +
			"☕️ кофейня и\n" +
			"✨ площадка для мероприятий.\n\n" +
			"Обязательно ознакомьтесь с разделом /about — там вы найдете информацию о наших локациях и правилах поведения в них.\n\n" +
			"Выберите команду про продолжения диалога:\n\n" +
			"команды:\n" +
			"/start – перезапуск\n" +
			"/wifi – получить пароль от вайфай\n" +
			"/meetingroom – забронировать переговорку\n" +
			"/printout – отправить документы на печать\n" +
			"/events – информация о мероприятиях\n" +
			"/about – информация о площадке и схема\n" +
			"/language – смена языка интерфейса\n"
	}
	
	*args.CurrentCommand = ""

	_, err := bot.Send(msg)
		
	return err
}