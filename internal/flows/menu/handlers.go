package flows

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

func send(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	// Canva link (public)
	const menuURL = "https://www.canva.com/design/DAG_yoyNcU0/Iqe7aXbnk4U3T17u4HjExQ/view?utm_content=DAG_yoyNcU0&utm_campaign=designshare&utm_medium=link2&utm_source=uniquelinks&utlId=h71d79c93bb"

	text := fmt.Sprintf("%s\n<a href=\"%s\">%s</a>",
		p.Sprintf("menu_link_title"),
		menuURL,
		p.Sprintf("menu_open_btn"),
	)

	msg := tgbotapi.NewMessage(s.ChatID, text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = false // true если не хочешь превью

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(p.Sprintf("menu_open_btn"), menuURL),
		),
	)

	if _, err := d.Bot.Send(msg); err != nil {
		// for logging purpose
	}

	s.Flow, s.Step = "", ""
	return MenuDone, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return MenuDone, nil
}
