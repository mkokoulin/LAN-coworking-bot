package commands

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)

type L struct {
	Lang string
	Desc string
}

var Languages = []L{
	{
		Lang: "🇺🇸 English",
		Desc: "Choose the interface language",
	},
	{
		Lang: "🇷🇺 Русский",
		Desc: "Выберите язык интерфейса",
	},
}

func Language(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if update.Message.Text != "" {
		for _, v := range Languages {
			if update.Message.Text == v.Lang {
				*args.Language = v.Lang
				*args.CurrentCommand = ""

				if v.Lang == Languages[0].Lang {
					msg.Text = fmt.Sprintf("Selected %s", v.Lang)
				} else if v.Lang == Languages[1].Lang {
					msg.Text = fmt.Sprintf("Выбран %s", v.Lang)
				}

				msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
					RemoveKeyboard: true,
					Selective: false,
				}

				_, err := bot.Send(msg)
				return err
			}
		}
	} else {
		return nil
	}

	if *args.CurrentCommand == LANGUAGE || *args.CurrentCommand == START {
		descs := []string {}
	
		for _, v := range Languages {
			descs = append(descs, v.Desc)
		}
	
		msg.Text = fmt.Sprintf("%v 🌎", strings.Join(descs, " / "))
		
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("🇺🇸 English"),
				tgbotapi.NewKeyboardButton("🇷🇺 Русский"),
			),
		)	
		_, err := bot.Send(msg)
		return err
			
	}

	return nil
}