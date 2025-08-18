package flow

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

// -------- Хендлер --------

func coworkingHome(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	// Если пришли по коллбэку не сюда — мягко переадресуемся
	if ev.Kind == botengine.EventCallback && !strings.HasPrefix(ev.CallbackData, "/coworking") {
		s.Step = CoworkingHome
		return botengine.InternalContinue, nil
	}

	text := p.Sprintf("coworking_intro") + "\n\n" +
		p.Sprintf("coworking_prices") + "\n\n" +
		p.Sprintf("coworking_meeting") + "\n\n" +
		p.Sprintf("coworking_options")

	msg := tgbotapi.NewMessage(ev.ChatID, text)
	// Без сложной разметки, чтобы не экранировать; переносы строк уже готовы в переводах
	kb := ui.Inline(
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
	msg.ReplyMarkup = kb

	if _, err := d.Bot.Send(msg); err != nil {
		return CoworkingHome, err
	}
	s.Step = CoworkingHome
	return botengine.InternalContinue, nil
}
