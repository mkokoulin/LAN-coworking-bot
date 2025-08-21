package flows

import (
	"context"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

func send(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	// выбираем файл по языку
	base := "menu_eng"
	if s.Lang == "ru" {
		base = "menu_rus"
	}
	path := fmt.Sprintf("internal/assets/%s.pdf", base)

	f, err := os.Open(path)
	if err != nil {
		// если файла нет — сообщаем и выходим мягко
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("menu_unavailable"))
		s.Flow, s.Step = "", ""
		return MenuDone, nil
	}
	defer f.Close()

	doc := tgbotapi.NewDocument(s.ChatID, tgbotapi.FileReader{
		Name:   "menu.pdf",
		Reader: f,
	})
	_, _ = d.Bot.Send(doc)

	s.Flow, s.Step = "", ""
	return MenuDone, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return MenuDone, nil
}
