package start

import (
	"context"
	"log"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

func show(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	// при /start можно «обнулить» авторизацию/контекст, если по логике нужно
	// s.IsAuthorized = false
	// оставляем s.Flow/Step — их обновит возврат шага ниже

	// p := d.Printer(s.Lang)
	// if err := ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("start_message")); err != nil {
	// 	log.Printf("[flow start.show] send error chat=%d: %v", s.ChatID, err)
	// 	return StepShow, err
	// }

	p := d.Printer(s.Lang)
	// kb := ui.Inline( // компактные быстрые действия
	// 	ui.Row(
	// 		ui.Cb("🎁 First visit", "/booking"),
	// 		ui.Cb("📅 Meeting room", "/meetingroom"),
	// 	),
	// 	ui.Row(
	// 		ui.Cb("☕ Bar", "/bar"),
	// 		ui.Cb("🎟 Events", "/events"),
	// 	),
	// 	ui.Row(
	// 		ui.Cb("ℹ️ About", "/about"),
	// 		ui.Cb("🌐 Language", "/language"),
	// 	),
	// )
	if err := ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("start_message")); err != nil {
		log.Printf("[flow start.show] send error chat=%d: %v", s.ChatID, err)
		return StepShow, err
	}
	// _ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("start_message"), kb)

	// завершаем сценарий
	s.Flow, s.Step = "", ""
	return StepDone, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return StepDone, nil
}
