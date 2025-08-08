package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/locales"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

var LanguageOptions = []types.LangOption{
	{Code: "en", Label: "🇺🇸 English"},
	{Code: "ru", Label: "🇷🇺 Русский"},
}

func LanguageCommand(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, services types.Services, state *types.ChatStorage) error {
	userInput := update.Message.Text
	chatID := update.Message.Chat.ID

	// 1. Обработка выбора языка
	for _, opt := range LanguageOptions {
		if userInput == opt.Label {
			state.Language = opt.Code
			state.CurrentCommand = ""

			p := locales.Printer(opt.Code)
			confirm := tgbotapi.NewMessage(chatID, p.Sprintf("language_selected", opt.Label))
			confirm.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}

			_, err := bot.Send(confirm)	
			return err
		}
	}

	// 2. Показываем выбор языка
	var buttons []tgbotapi.KeyboardButton
	for _, opt := range LanguageOptions {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(opt.Label))
	}

	p := locales.Printer(state.Language)
	msg := tgbotapi.NewMessage(chatID, p.Sprintf("language_prompt"))
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)

	_, err := bot.Send(msg)
	return err
}
