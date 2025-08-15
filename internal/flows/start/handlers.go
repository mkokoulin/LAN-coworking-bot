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
	s.IsAuthorized = false
	// оставляем s.Flow/Step — их обновит возврат шага ниже

	p := d.Printer(s.Lang)
	if err := ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("start_message")); err != nil {
		log.Printf("[flow start.show] send error chat=%d: %v", s.ChatID, err)
		return StepShow, err
	}

	// завершаем сценарий
	s.Flow, s.Step = "", ""
	return StepDone, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return StepDone, nil
}
