package flows

import (
	"context"
	"strconv"
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
	if err := ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("wifi_select_network"), kb); err != nil {
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
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("wifi_ask_confirmation"))
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
		// считаем попытки в s.Data["wifi_attempts"]
		if s.Data == nil {
			s.Data = map[string]interface{}{}
		}
		attempts := toInt(s.Data["wifi_attempts"])
		attempts++
		s.Data["wifi_attempts"] = attempts

		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("wrong_secret"))
		if attempts >= 3 {
			s.Flow, s.Step = "", ""
			return WifiDone, nil
		}
		return WifiWaitCode, nil
	}

	// успех — сбрасываем попытки и отмечаем авторизацию
	if s.Data == nil {
		s.Data = map[string]interface{}{}
	}
	delete(s.Data, "wifi_attempts")
	s.Data["is_authorized"] = "true"

	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("wifi_coworking", d.Cfg.CoworkingWifiPassword), ui.RemoveKeyboard())
	s.Flow, s.Step = "", ""
	return WifiDone, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return WifiDone, nil
}

// --- helpers ---

// toInt безопасно приводит interface{} к int (учитывая возможные типы после BSON)
func toInt(v interface{}) int {
	switch t := v.(type) {
	case int:
		return t
	case int32:
		return int(t)
	case int64:
		return int(t)
	case float64:
		return int(t)
	case string:
		if n, err := strconv.Atoi(t); err == nil {
			return n
		}
	}
	return 0
}
