// package ui
package ui

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AnswerCallback — короткий тост (или alert) по клику инлайн-кнопки.
// text можно оставить пустым (""), чтобы просто погасить спиннер.
func AnswerCallback(bot *tgbotapi.BotAPI, callbackQueryID, text string, showAlert ...bool) error {
	cfg := tgbotapi.NewCallback(callbackQueryID, text)
	if len(showAlert) > 0 && showAlert[0] {
		cfg.ShowAlert = true // показывать модальный alert вместо тоста
	}
	_, err := bot.Request(cfg)
	return err
}

// Удобные синонимы, если нравится более говорящая семантика:
func Toast(bot *tgbotapi.BotAPI, callbackQueryID, text string) error {
	return AnswerCallback(bot, callbackQueryID, text)
}
func Alert(bot *tgbotapi.BotAPI, callbackQueryID, text string) error {
	return AnswerCallback(bot, callbackQueryID, text, true)
}
