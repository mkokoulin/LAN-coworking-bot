package flows

import (
	"context"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

func prompt(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)
	kb := ui.Inline(
		ui.Row(
			ui.Cb("🇺🇸 English", "lang:en"),
			ui.Cb("🇷🇺 Русский", "lang:ru"),
		),
	)
	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("language_prompt"), kb)
	return LangWaitChoice, nil
}

func waitChoice(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if ev.Kind != botengine.EventCallback {
		return LangWaitChoice, nil
	}

	var (
		newLang string
		label   string
	)
	switch ev.CallbackData {
	case "lang:en":
		newLang, label = "en", "🇺🇸 English"
	case "lang:ru":
		newLang, label = "ru", "🇷🇺 Русский"
	default:
		return LangWaitChoice, nil
	}

	// Сохраняем язык и подтверждаем
	s.Lang = newLang
	p := d.Printer(s.Lang)
	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("language_selected", label))

	// // Куда возвращаемся после выбора языка
	// if s.PendingCmd == "" {
	// 	s.PendingCmd = "start"
	// }

	// Сбрасываем текущий флоу, просим FSM продолжить сразу (он сам дернёт PendingCmd)
	s.Flow, s.Step = "", ""
	return botengine.InternalContinue, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return LangDone, nil
}
