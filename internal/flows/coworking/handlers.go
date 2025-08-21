package flows

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

func ackCallback(d botengine.Deps, ev botengine.Event) {
	if ev.CallbackQueryID == "" { // имя поля может быть CallbackID/CallbackQueryID — возьми то, что у тебя в Event
		return
	}
	_, _ = d.Bot.Request(tgbotapi.NewCallback(ev.CallbackQueryID, "")) // AnswerCallbackQuery
}

func coworkingHome(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	// Навигация по инлайн-кнопкам вида "/bar", "/events", ...
	if ev.Kind == botengine.EventCallback && strings.HasPrefix(ev.CallbackData, "/") {
		ackCallback(d, ev)

		target := strings.TrimPrefix(ev.CallbackData, "/") // "bar", "events", "meetingroom", ...

		// Маппинг на реальные флоу/стартовые шаги
		flowName := target
		step := target + ":home"

		switch target {
		case "bar":
			step = "bar:home"
		case "events":
			step = "events:list"
		case "meetingroom":
			flowName = "meeting"
			step = "meeting:prompt"
		case "booking":
			step = "booking:info"
		case "about":
			step = "about:send"
		case "language":
			step = "language:prompt"
		}

		s.Flow = types.Flow(flowName)
		s.Step = types.Step(step)
		return botengine.InternalContinue, nil
	}

	// Рендер текста /coworking
	text := p.Sprintf("coworking_intro") + "\n\n" +
		p.Sprintf("coworking_prices") + "\n\n" +
		p.Sprintf("coworking_meeting") + "\n\n" +
		p.Sprintf("coworking_options")

	msg := tgbotapi.NewMessage(ev.ChatID, text)
	msg.ReplyMarkup = ui.Inline(
		ui.Row(
			ui.Cb(p.Sprintf("coworking_btn_booking"), "/booking"),
			ui.Cb(p.Sprintf("coworking_btn_meetingroom"), "/meetingroom"),
		),
		ui.Row(
			ui.Cb(p.Sprintf("coworking_btn_events"), "/events"),
			ui.Cb(p.Sprintf("coworking_btn_bar"), "/bar"),
		),
		ui.Row(
			ui.Cb(p.Sprintf("coworking_btn_about"), "/about"),
			ui.Cb(p.Sprintf("coworking_btn_language"), "/language"),
		),
	)

	if _, err := d.Bot.Send(msg); err != nil {
		return CoworkingHome, err
	}
	s.Step = CoworkingHome
	return botengine.InternalContinue, nil
}
