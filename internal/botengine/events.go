package botengine

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func Classify(u tgbotapi.Update) Event {
	ev := Event{}

	// 1) CallbackQuery — первым
	if cq := u.CallbackQuery; cq != nil {
		ev.Kind = EventCallback
		ev.CallbackData = cq.Data
		ev.CallbackQueryID = cq.ID

		if cq.Message != nil && cq.Message.Chat != nil {
			ev.ChatID = cq.Message.Chat.ID
			ev.MessageID = cq.Message.MessageID
		}
		if cq.InlineMessageID != "" {
			ev.InlineMessageID = cq.InlineMessageID
		}
		if cq.From != nil {
			ev.FromUserName = cq.From.UserName
			ev.FromUserID = cq.From.ID
		}
		return ev
	}

	// 2) Message
	if m := u.Message; m != nil {
		ev.ChatID = m.Chat.ID
		if m.IsCommand() {
			ev.Kind = EventCommand
			ev.Command = m.Command()
		} else {
			ev.Kind = EventText
			ev.Text = m.Text
		}
		if m.From != nil {
			ev.FromUserName = m.From.UserName
			ev.FromUserID = m.From.ID
		}
		return ev
	}

	// 3) my_chat_member
	if u.MyChatMember != nil {
		ev.ChatID = u.MyChatMember.Chat.ID // Chat — значение, не указатель
	}
	return ev
}

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
