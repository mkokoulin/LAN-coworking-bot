package flows

import (
	"context"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

func send(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	// 1) Пытаемся отправить карту (jpg). Ошибки не фейлим — просто идём дальше.
	if b, err := os.ReadFile("internal/assets/coworking_scheme.jpg"); err == nil && len(b) > 0 {
		photo := tgbotapi.NewPhoto(s.ChatID, tgbotapi.FileBytes{
			Name:  "coworking_scheme.jpg",
			Bytes: b,
		})
		_, _ = d.Bot.Send(photo)
	}

	// 2) Отправляем текст (HTML локалью)
	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("about_text"))

	// 3) Завершаем сценарий
	s.Flow, s.Step = "", ""
	return AboutDone, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return AboutDone, nil
}
