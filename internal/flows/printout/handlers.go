package flows

import (
	"context"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

func showInfo(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)
	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("printout_info"))

	// завершаем сразу
	s.Flow, s.Step = "", ""
	return PrintoutDone, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return PrintoutDone, nil
}
