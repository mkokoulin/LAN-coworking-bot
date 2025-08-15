package flows

import (
	"context"
	"strings"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

func start(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	kb := ui.Inline(
		ui.Row(
			ui.Cb("Guest", "wifi:guest"),
			ui.Cb("Coworking", "wifi:coworking"),
		),
	)
	if err := ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("select_network"), kb); err != nil {
		return WifiStart, err
	}
	return WifiWaitChoice, nil
}

func waitChoice(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	if ev.Kind != botengine.EventCallback {
		return WifiWaitChoice, nil
	}

	switch ev.CallbackData {
	case "wifi:guest":
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("wifi_guest", d.Cfg.GuestWifiPassword), ui.RemoveKeyboard())
		s.Flow, s.Step = "", ""
		return WifiDone, nil

	case "wifi:coworking":
		if s.Data["is_authorized"] == "true" {
			_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("wifi_coworking", d.Cfg.CoworkingWifiPassword), ui.RemoveKeyboard())
			s.Flow, s.Step = "", ""
			return WifiDone, nil
		}
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("ask_confirmation"))
		return WifiWaitCode, nil
	}

	return WifiWaitChoice, nil
}

func waitCode(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	if ev.Kind != botengine.EventText {
		return WifiWaitCode, nil
	}
	code := strings.TrimSpace(ev.Text)

	ok, err := d.Svcs.CoworkersSheets.ValidateSecret(ctx, code)
	if err != nil {
		return WifiWaitCode, err
	}
	if !ok {
		s.Attempts++
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("wrong_secret"))
		if s.Attempts >= 3 {
			s.Flow, s.Step = "", ""
			return WifiDone, nil
		}
		return WifiWaitCode, nil
	}

	if s.Data == nil {
		s.Data = map[string]interface{}{}
	}
	s.Data["is_authorized"] = "true"

	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("wifi_coworking", d.Cfg.CoworkingWifiPassword), ui.RemoveKeyboard())
	s.Flow, s.Step = "", ""
	return WifiDone, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return WifiDone, nil
}
