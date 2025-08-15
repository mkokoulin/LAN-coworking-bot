package flow

import (
	"context"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

func info(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)
	if err := ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("booking_text")); err != nil {
		return BookInfo, err
	}
	s.Flow, s.Step = "", ""
	return BookDone, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return BookDone, nil
}
