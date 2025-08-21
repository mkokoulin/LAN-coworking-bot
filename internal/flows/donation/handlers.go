package flows

import (
	"context"
	"reflect"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

func donationHome(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if ok, next := botengine.InterceptSlashNav(ev, func() { ackCallback(d, ev) }); ok {
		return next, nil
	}

	if ev.Kind == botengine.EventCallback && strings.HasPrefix(ev.CallbackData, "donation:") {
		switch ev.CallbackData {
		case "donation:card":
			return donationCard(ctx, ev, d, s)
		case "donation:copied":
			return donationCopied(ctx, ev, d, s)
		case "donation:done":
			return donationDone(ctx, ev, d, s)
		case "donation:home":
			// отрисуем экран ниже
		default:
			ackCallback(d, ev)
		}
	}

	p := d.Printer(s.Lang)

	text := strings.Join([]string{
		p.Sprintf("donation_title"),
		"",
		p.Sprintf("donation_howto"),
		"• " + p.Sprintf("donation_opt_events"),
		"• " + p.Sprintf("donation_opt_bar"),
		"• " + p.Sprintf("donation_opt_cowork"),
		"• " + p.Sprintf("donation_opt_card"),
		"",
		p.Sprintf("donation_choose"),
	}, "\n")

	kb := ui.Inline(
		ui.Row(
			ui.Cb(p.Sprintf("donation_btn_card"), "donation:card"),
		),
		ui.Row(
			ui.CbCmd(p.Sprintf("donation_btn_events"), "events"),
			ui.CbCmd(p.Sprintf("donation_btn_bar"), "bar"),
		),
		ui.Row(
			ui.CbCmd(p.Sprintf("donation_btn_cowork"), "booking"),
			ui.CbCmd(p.Sprintf("donation_btn_home"), "start"),
		),
	)

	msg := tgbotapi.NewMessage(ev.ChatID, text)
	msg.ReplyMarkup = kb
	_, _ = d.Bot.Send(msg)

	return DonationHome, nil
}

func donationCard(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if ev.Kind == botengine.EventCallback {
		ackCallback(d, ev)
	}
	p := d.Printer(s.Lang)

	text := p.Sprintf("donation_btn_card") + "\n\n" +
		p.Sprintf("donation_card_label") + "\n<code>" + cardNumber + "</code>\n\n" +
		p.Sprintf("donation_card_note") + "\n\n" +
		p.Sprintf("donation_thanks")

	kb := ui.Inline(
		ui.Row(
			ui.Cb(p.Sprintf("donation_btn_copy"), "donation:copied"),
			ui.Cb(p.Sprintf("donation_btn_back"), "donation:home"),
		),
		ui.Row(
			ui.Cb(p.Sprintf("donation_btn_done"), "donation:done"),
		),
	)

	msg := tgbotapi.NewMessage(ev.ChatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = kb
	_, _ = d.Bot.Send(msg)

	return DonationCard, nil
}

func donationCopied(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if ev.Kind == botengine.EventCallback {
		ackCallback(d, ev)
	}
	p := d.Printer(s.Lang)

	notice := p.Sprintf("donation_copy_hint") + "\n<code>" + cardNumber + "</code>"
	n := tgbotapi.NewMessage(ev.ChatID, notice)
	n.ParseMode = "HTML"
	_, _ = d.Bot.Send(n)

	return DonationCard, nil
}

func donationDone(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if ev.Kind == botengine.EventCallback {
		ackCallback(d, ev)
	}
	p := d.Printer(s.Lang)
	msg := tgbotapi.NewMessage(ev.ChatID, p.Sprintf("donation_always"))
	_, _ = d.Bot.Send(msg)
	return DonationDone, nil
}

func ackCallback(d botengine.Deps, ev botengine.Event) {
	id := callbackID(ev)
	if id == "" {
		return
	}
	cb := tgbotapi.NewCallback(id, "")
	cb.ShowAlert = false
	_, _ = d.Bot.Request(cb)
}

func callbackID(ev botengine.Event) string {
	if id := tryMethodString(ev, "CallbackID"); id != "" {
		return id
	}
	if id := tryMethodString(ev, "GetCallbackID"); id != "" {
		return id
	}

	v := reflect.ValueOf(ev)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return ""
	}

	for _, name := range []string{"CallbackQueryID", "CallbackID", "CallbackQueryId"} {
		if f := v.FieldByName(name); f.IsValid() && f.Kind() == reflect.String {
			return f.String()
		}
	}

	for _, name := range []string{"CallbackQuery", "TGCallback", "TgCallback"} {
		if f := v.FieldByName(name); f.IsValid() && f.CanInterface() {
			switch q := f.Interface().(type) {
			case *tgbotapi.CallbackQuery:
				if q != nil {
					return q.ID
				}
			case tgbotapi.CallbackQuery:
				return q.ID
			}
		}
	}

	return ""
}

func tryMethodString(ev any, method string) string {
	val := reflect.ValueOf(ev)
	m := val.MethodByName(method)
	if !m.IsValid() || m.Type().NumIn() != 0 || m.Type().NumOut() != 1 || m.Type().Out(0).Kind() != reflect.String {
		return ""
	}
	out := m.Call(nil)
	return out[0].String()
}
