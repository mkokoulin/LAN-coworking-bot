package ui

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SendHTML — отправка HTML-сообщения с опциональным markup
// Поддерживает: InlineKeyboardMarkup, ReplyKeyboardMarkup, ReplyKeyboardRemove
func SendHTML(bot *tgbotapi.BotAPI, chatID int64, html string, markup ...interface{}) error {
	msg := tgbotapi.NewMessage(chatID, html)
	msg.ParseMode = "HTML"

	if len(markup) > 0 && markup[0] != nil {
		switch m := markup[0].(type) {
		case tgbotapi.InlineKeyboardMarkup:
			msg.ReplyMarkup = m
		case tgbotapi.ReplyKeyboardMarkup:
			msg.ReplyMarkup = m
		case tgbotapi.ReplyKeyboardRemove:
			msg.ReplyMarkup = m
		}
	}

	_, err := bot.Send(msg)
	return err
}

func SendText(bot *tgbotapi.BotAPI, chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	return err
}
