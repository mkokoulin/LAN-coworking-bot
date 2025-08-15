package botengine

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// internal/botengine/adapter.go (пример)
func FromUpdate(u tgbotapi.Update) Event {
    var ev Event
    if u.CallbackQuery != nil {
        ev.Kind = EventCallback
        if u.CallbackQuery.Message != nil {
            ev.ChatID = u.CallbackQuery.Message.Chat.ID
            ev.MessageID = u.CallbackQuery.Message.MessageID // 👈 вот он
        }
        if u.CallbackQuery.InlineMessageID != "" {
            ev.InlineMessageID = u.CallbackQuery.InlineMessageID // 👈 inline-режим
        }
        ev.CallbackData = u.CallbackQuery.Data
        ev.CallbackQueryID = u.CallbackQuery.ID
        if u.CallbackQuery.From != nil {
            ev.FromUserName = u.CallbackQuery.From.UserName
            ev.FromUserID = u.CallbackQuery.From.ID
        }
        return ev
    }

	if u.Message != nil {
		ev.Kind = EventText
		ev.ChatID = u.Message.Chat.ID
		ev.Text = u.Message.Text
		if u.Message.From != nil {
			ev.FromUserName = u.Message.From.UserName
			ev.FromUserID = u.Message.From.ID
		}
	}

	return ev
}
