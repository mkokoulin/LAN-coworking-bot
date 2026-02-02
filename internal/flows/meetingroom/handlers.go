// internal/flows/meeting.go
package flows

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

const (
	slotStartHour = 10
	slotEndHour   = 22
	slotStepMin   = 30
)

const (
	keyDate    = "meeting.date"    // YYYY-MM-DD
	keyStart   = "meeting.start"   // HH:MM
	keyEnd     = "meeting.end"     // HH:MM
	keyContact = "meeting.contact" // fallback/entered contact
)

func prompt(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	kb := calendarKeyboard(time.Now())
	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_pick_date"), kb)

	return MeetWaitInterval, nil
}

func waitInterval(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	// wait for inline button presses
	if ev.Kind != botengine.EventCallback {
		return MeetWaitInterval, nil
	}

	p := d.Printer(s.Lang)
	data := strings.TrimSpace(ev.CallbackData)

	switch {
	case strings.HasPrefix(data, "meet:date:"):
		date := strings.TrimPrefix(data, "meet:date:")
		ensureData(s)
		s.Data[keyDate] = date

		kb := startTimeKeyboard()
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_pick_start_time", date), kb)

	case strings.HasPrefix(data, "meet:start:"):
		start := strings.TrimPrefix(data, "meet:start:")
		if s.Data == nil || s.Data[keyDate] == nil {
			_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_select_date_first"))
			kb := calendarKeyboard(time.Now())
			_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_pick_date"), kb)
			return MeetWaitInterval, nil
		}
		s.Data[keyStart] = start

		endKB := endTimeKeyboard(start)
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_pick_end_time", s.Data[keyDate], start), endKB)

	case strings.HasPrefix(data, "meet:end:"):
		end := strings.TrimPrefix(data, "meet:end:")
		if s.Data == nil || s.Data[keyDate] == nil || s.Data[keyStart] == nil {
			_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_flow_broken"))
			kb := calendarKeyboard(time.Now())
			_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_pick_date"), kb)
			return MeetWaitInterval, nil
		}

		start := fmt.Sprint(s.Data[keyStart])

		if err := validateInterval(start, end); err != nil {
			_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_invalid_interval", err.Error()))
			endKB := endTimeKeyboard(start)
			_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_pick_end_time", s.Data[keyDate], start), endKB)
			return MeetWaitInterval, nil
		}

		s.Data[keyEnd] = end

		date := fmt.Sprint(s.Data[keyDate])
		intervalHuman := fmt.Sprintf("%s %s–%s", date, start, end)

		contactLabel, hasUsername := eventContactLabel(ev)

		// If username is hidden / missing -> ask how to contact
		if !hasUsername {
			ensureData(s)
			s.Data[keyContact] = contactLabel // fallback (id / name if you have it somewhere)

			_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_need_contact", intervalHuman))
			return MeetWaitContact, nil
		}

		adminText := p.Sprintf("meeting_request_admin", intervalHuman, contactLabel)
		_, _ = d.Bot.Send(tgbotapi.NewMessage(d.Cfg.AdminChatId, adminText))

		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_confirm_interval", intervalHuman))

		cleanupMeeting(s)
		s.Flow, s.Step = "", ""
		return MeetDone, nil
	}

	return MeetWaitInterval, nil
}

func waitContact(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	// Need plain text message from user (not callbacks)
	if ev.Kind == botengine.EventCallback {
		return MeetWaitContact, nil
	}

	text := strings.TrimSpace(ev.Text)
	if text == "" {
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_empty"))
		return MeetWaitContact, nil
	}
	if len([]rune(text)) > 64 {
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_contact_too_long"))
		return MeetWaitContact, nil
	}

	if s.Data == nil || s.Data[keyDate] == nil || s.Data[keyStart] == nil || s.Data[keyEnd] == nil {
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_flow_broken"))
		kb := calendarKeyboard(time.Now())
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_pick_date"), kb)
		return MeetWaitInterval, nil
	}

	date := fmt.Sprint(s.Data[keyDate])
	start := fmt.Sprint(s.Data[keyStart])
	end := fmt.Sprint(s.Data[keyEnd])
	intervalHuman := fmt.Sprintf("%s %s–%s", date, start, end)

	// Fallback identifier (at least id:123). If you later add FirstName/LastName — it will become prettier.
	fallbackLabel, _ := eventContactLabel(ev)
	contact := text

	adminText := p.Sprintf("meeting_request_admin", intervalHuman, fmt.Sprintf("%s | %s", contact, fallbackLabel))
	_, _ = d.Bot.Send(tgbotapi.NewMessage(d.Cfg.AdminChatId, adminText))

	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_confirm_interval", intervalHuman))

	cleanupMeeting(s)
	s.Flow, s.Step = "", ""
	return MeetDone, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return MeetDone, nil
}

// ---------- helpers ----------

func ensureData(s *types.Session) {
	if s.Data == nil {
		s.Data = map[string]interface{}{}
	}
}

func cleanupMeeting(s *types.Session) {
	if s.Data == nil {
		return
	}
	delete(s.Data, keyDate)
	delete(s.Data, keyStart)
	delete(s.Data, keyEnd)
	delete(s.Data, keyContact)
}

// Minimal-invasive user label based ONLY on your current Event fields.
func eventContactLabel(ev botengine.Event) (label string, hasUsername bool) {
	if strings.TrimSpace(ev.FromUserName) != "" {
		return "@" + ev.FromUserName, true
	}

	// Username hidden: still show stable id to admin for fallback.
	if ev.FromUserID != 0 {
		return "id:" + strconv.FormatInt(ev.FromUserID, 10), false
	}

	return "unknown", false
}

// ---------- UI builders ----------

func calendarKeyboard(now time.Time) tgbotapi.InlineKeyboardMarkup {
	loc := now.Location()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	var rows [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < 7; i++ {
		day := today.AddDate(0, 0, i)
		label := day.Format("Mon 02.01")
		cb := "meet:date:" + day.Format("2006-01-02")
		rows = append(rows, ui.Row(ui.Cb(label, cb)))
	}
	return ui.Inline(rows...)
}

func startTimeKeyboard() tgbotapi.InlineKeyboardMarkup {
	slots := generateSlots(slotStartHour, slotEndHour, slotStepMin)
	return timesToKeyboard("meet:start:", slots)
}

func endTimeKeyboard(startHHMM string) tgbotapi.InlineKeyboardMarkup {
	start, _ := time.Parse("15:04", startHHMM)
	minEnd := start.Add(time.Minute * slotStepMin)
	slots := generateSlotsFrom(minEnd, slotEndHour, slotStepMin)
	return timesToKeyboard("meet:end:", slots)
}

func timesToKeyboard(prefix string, times []string) tgbotapi.InlineKeyboardMarkup {
	const perRow = 4
	var rows [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < len(times); i += perRow {
		end := i + perRow
		if end > len(times) {
			end = len(times)
		}
		var row []tgbotapi.InlineKeyboardButton
		for _, t := range times[i:end] {
			row = append(row, ui.Cb(t, prefix+t))
		}
		rows = append(rows, row)
	}
	return ui.Inline(rows...)
}

// ---------- Validation & slots ----------

func generateSlots(startHour, endHour, stepMin int) []string {
	loc := time.Now().Location()
	base := time.Date(0, 1, 1, startHour, 0, 0, 0, loc)
	var out []string
	for t := base; t.Hour() < endHour || (t.Hour() == endHour && t.Minute() == 0); t = t.Add(time.Duration(stepMin) * time.Minute) {
		out = append(out, t.Format("15:04"))
		if t.Hour() == endHour && t.Minute() == 0 {
			break
		}
	}
	return out
}

func generateSlotsFrom(from time.Time, endHour, stepMin int) []string {
	loc := from.Location()
	limit := time.Date(from.Year(), from.Month(), from.Day(), endHour, 0, 0, 0, loc)

	minute := (from.Minute() / stepMin) * stepMin
	aligned := time.Date(from.Year(), from.Month(), from.Day(), from.Hour(), minute, 0, 0, loc)
	if aligned.Before(from) {
		aligned = aligned.Add(time.Duration(stepMin) * time.Minute)
	}

	var out []string
	for t := aligned; !t.After(limit); t = t.Add(time.Duration(stepMin) * time.Minute) {
		out = append(out, t.Format("15:04"))
	}
	return out
}

func validateInterval(startHHMM, endHHMM string) error {
	start, err := time.Parse("15:04", startHHMM)
	if err != nil {
		return fmt.Errorf("неверное время начала")
	}
	end, err := time.Parse("15:04", endHHMM)
	if err != nil {
		return fmt.Errorf("неверное время окончания")
	}

	open := time.Date(0, 1, 1, slotStartHour, 0, 0, 0, time.UTC)
	close := time.Date(0, 1, 1, slotEndHour, 0, 0, 0, time.UTC)

	s := time.Date(0, 1, 1, start.Hour(), start.Minute(), 0, 0, time.UTC)
	e := time.Date(0, 1, 1, end.Hour(), end.Minute(), 0, 0, time.UTC)

	if s.Before(open) || e.After(close) {
		return fmt.Errorf("время должно быть в диапазоне %02d:00–%02d:00", slotStartHour, slotEndHour)
	}
	if !e.After(s) {
		return fmt.Errorf("окончание должно быть позже начала")
	}
	if start.Minute()%slotStepMin != 0 || end.Minute()%slotStepMin != 0 {
		return fmt.Errorf("используйте шаг %d минут", slotStepMin)
	}
	return nil
}
