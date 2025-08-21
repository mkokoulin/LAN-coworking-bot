package ui

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Inline(...) — сборка InlineKeyboardMarkup из рядов
func Inline(rows ...[]tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

// Row(...) — один ряд инлайн-кнопок
func Row(btns ...tgbotapi.InlineKeyboardButton) []tgbotapi.InlineKeyboardButton {
	return btns
}

// Cb(text, data) — inline-кнопка с callback data
func Cb(text, data string) tgbotapi.InlineKeyboardButton {
	btn := tgbotapi.NewInlineKeyboardButtonData(text, data)
	return btn
}

// RemoveKeyboard — убрать reply-клавиатуру
func RemoveKeyboard() tgbotapi.ReplyKeyboardRemove {
	return tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
}
