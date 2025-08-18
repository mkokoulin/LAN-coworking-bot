package flow

import (
	"context"
	"reflect"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// -------- Константы флоу --------

// Выбери реальную логику определения языка (например, из сессии).
func userLang(_ *types.Session) language.Tag {
	// TODO: верни language.English или language.Russian из твоего s.Locale / s.Lang
	return language.Russian
}

// -------- Хендлеры --------

func donationHome(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	// Роутинг: все donation:* приходят сюда
	if ev.Kind == botengine.EventCallback && strings.HasPrefix(ev.CallbackData, "donation:") {
		switch ev.CallbackData {
		case "donation:card":
			return donationCard(ctx, ev, d, s)
		case "donation:copied":
			return donationCopied(ctx, ev, d, s)
		case "donation:done":
			return donationDone(ctx, ev, d, s)
		case "donation:home":
			// просто отрисуем домашний экран ниже
		default:
			ackCallback(d, ev)
		}
	}

	p := message.NewPrinter(userLang(s))

	text := strings.Join([]string{
		p.Sprintf("Letters & Numbers is an independent project. We exist thanks to your support ❤️"),
		"",
		p.Sprintf("How you can support:"),
		"• " + p.Sprintf("Attend our 🎟 events"),
		"• " + p.Sprintf("Grab a coffee and desserts at the ☕ bar"),
		"• " + p.Sprintf("Work from our 💻 coworking"),
		"• " + p.Sprintf("Or send a 💳 card donation (add note “lan cats”)"),
		"",
		p.Sprintf("Choose an option:"),
	}, "\n")

	kb := ui.Inline(
		ui.Row(
			ui.Cb(p.Sprintf("💳 Card donation"), "donation:card"),
		),
		ui.Row(
			ui.Cb(p.Sprintf("🎟 Events"), "/events"),
			ui.Cb(p.Sprintf("☕ Bar"), "/bar"),
		),
		ui.Row(
			ui.Cb(p.Sprintf("💻 Coworking"), "/booking"),
			ui.Cb(p.Sprintf("⬅️ Home"), "/start"),
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
	p := message.NewPrinter(userLang(s))

	text := p.Sprintf("💳 Card donation") + "\n\n" +
		p.Sprintf("Card number:") + "\n`" + cardNumber + "`\n\n" +
		p.Sprintf("Important: add **lan cats** in payment note — this helps us understand the purpose.") + "\n\n" +
		p.Sprintf("Thank you for your support! 🐱")

	kb := ui.Inline(
		ui.Row(
			ui.Cb(p.Sprintf("📋 Copy number"), "donation:copied"),
			ui.Cb(p.Sprintf("⬅️ Back"), "donation:home"),
		),
		ui.Row(
			ui.Cb(p.Sprintf("✅ Done"), "donation:done"),
		),
	)

	msg := tgbotapi.NewMessage(ev.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = kb
	_, _ = d.Bot.Send(msg)

	return DonationCard, nil
}

func donationCopied(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if ev.Kind == botengine.EventCallback {
		ackCallback(d, ev)
	}
	p := message.NewPrinter(userLang(s))

	notice := p.Sprintf("Copy the card number from the message above:") + "\n`" + cardNumber + "`"
	n := tgbotapi.NewMessage(ev.ChatID, notice)
	n.ParseMode = "Markdown"
	_, _ = d.Bot.Send(n)

	return DonationCard, nil
}

func donationDone(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if ev.Kind == botengine.EventCallback {
		ackCallback(d, ev)
	}
	p := message.NewPrinter(userLang(s))
	msg := tgbotapi.NewMessage(ev.ChatID, p.Sprintf("Thank you! /donation is always available."))
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
	// 1) Если в Event есть метод CallbackID() string
	if id := tryMethodString(ev, "CallbackID"); id != "" {
		return id
	}
	if id := tryMethodString(ev, "GetCallbackID"); id != "" {
		return id
	}

	// 2) Пробуем поля по распространённым именам
	v := reflect.ValueOf(ev)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return ""
	}

	for _, name := range []string{
		"CallbackQueryID", "CallbackID", "CallbackQueryId",
	} {
		if f := v.FieldByName(name); f.IsValid() && f.Kind() == reflect.String {
			return f.String()
		}
	}

	// 3) Если внутри лежит tgbotapi.CallbackQuery — достаем ID оттуда
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
