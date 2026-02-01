package botengine

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func userRefFromTg(u *tgbotapi.User) UserRef {
	if u == nil {
		return UserRef{}
	}
	return UserRef{
		ID:        int64(u.ID),
		Username:  u.UserName,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func Classify(u tgbotapi.Update) Event {
	var ev Event

	// 1) Callback
	if u.CallbackQuery != nil {
		ev.Kind = EventCallback
		if u.CallbackQuery.Message != nil {
			ev.ChatID = u.CallbackQuery.Message.Chat.ID
			ev.MessageID = u.CallbackQuery.Message.MessageID
		}
		if u.CallbackQuery.InlineMessageID != "" {
			ev.InlineMessageID = u.CallbackQuery.InlineMessageID
		}
		ev.CallbackData = u.CallbackQuery.Data
		ev.CallbackQueryID = u.CallbackQuery.ID
		if u.CallbackQuery.From != nil {
			ev.FromUserName = u.CallbackQuery.From.UserName
			ev.FromUserID = u.CallbackQuery.From.ID
		}
		// Если callback содержит "/команду", сохраним её в ev.Command (ACK сделаем ниже в FSM).
		if strings.HasPrefix(ev.CallbackData, "/") {
			ev.Command = normalizeCommand(ev.CallbackData)
		}

		ev.From = userRefFromTg(u.CallbackQuery.From)
		return ev
	}

	// 2) Message
	if u.Message != nil {
		ev.ChatID = u.Message.Chat.ID
		ev.Text = u.Message.Text
		if u.Message.From != nil {
			ev.FromUserName = u.Message.From.UserName
			ev.FromUserID = u.Message.From.ID
		}
		if strings.HasPrefix(ev.Text, "/") {
			ev.Kind = EventCommand
			ev.Command = normalizeCommand(ev.Text)
		} else {
			ev.Kind = EventText
		}

		ev.From = userRefFromTg(u.Message.From)

		return ev
	}

	// 3) MyChatMember и пр. — если нужно, дополни
	return ev
}

func normalizeCommand(s string) string {
	// "/start@Bot arg1 arg2" -> "start"
	tok := firstToken(s)              // "/start@Bot"
	tok = strings.TrimPrefix(tok, "/") // "start@Bot"
	if i := strings.IndexByte(tok, '@'); i >= 0 {
		tok = tok[:i]
	}
	return tok
}

func firstToken(s string) string {
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return s
	}
	return parts[0]
}
