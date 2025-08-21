package botengine

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func ResolveChatID(u tgbotapi.Update) int64 {
	if u.Message != nil && u.Message.Chat != nil {
		return u.Message.Chat.ID
	}
	if u.CallbackQuery != nil && u.CallbackQuery.Message != nil && u.CallbackQuery.Message.Chat != nil {
		return u.CallbackQuery.Message.Chat.ID
	}
	if u.MyChatMember != nil {
		return u.MyChatMember.Chat.ID
	}
	return 0
}
