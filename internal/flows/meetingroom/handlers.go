package flows

import (
	"context"
	"fmt"
	"time"
	"strings"

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
	keyDate  = "meeting.date"   // YYYY-MM-DD
	keyStart = "meeting.start"  // HH:MM
	keyEnd   = "meeting.end"    // HH:MM
)

func prompt(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)
	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_prompt"))

	// показать календарь на неделю
	kb := calendarKeyboard(time.Now())
	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_pick_date"), kb)

	return MeetWaitInterval, nil // используем существующий стейт как "ожидание выбора"
}

func waitInterval(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	// Ждём нажатий инлайн-кнопок
	if ev.Kind != botengine.EventCallback {
		return MeetWaitInterval, nil
	}

	p := d.Printer(s.Lang)
	data := strings.TrimSpace(ev.CallbackData)

	switch {
	case strings.HasPrefix(data, "meet:date:"):
		date := strings.TrimPrefix(data, "meet:date:")
		if s.Data == nil { s.Data = map[string]interface{}{} }
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

		// Клавиатура конца — строго позже начала и не позже 22:00
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

		// Серверная валидация: 10:00–22:00, end > start, шаг 30 минут
		if err := validateInterval(start, end); err != nil {
			_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_invalid_interval", err.Error()))
			endKB := endTimeKeyboard(start)
			_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_pick_end_time", s.Data[keyDate], start), endKB)
			return MeetWaitInterval, nil
		}

		s.Data[keyEnd] = end

		// Готово: уведомляем и подтверждаем
		date := fmt.Sprint(s.Data[keyDate])
		intervalHuman := fmt.Sprintf("%s %s–%s", date, start, end)

		adminText := p.Sprintf("meeting_request_admin", intervalHuman)
		_, _ = d.Bot.Send(tgbotapi.NewMessage(d.Cfg.AdminChatId, adminText))

		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("meeting_confirm_interval", intervalHuman))

		// очистка
		delete(s.Data, keyDate)
		delete(s.Data, keyStart)
		delete(s.Data, keyEnd)

		s.Flow, s.Step = "", ""
		return MeetDone, nil
	}

	return MeetWaitInterval, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return MeetDone, nil
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
	// конец — строго после старта
	minEnd := start.Add(time.Minute * slotStepMin)
	slots := generateSlotsFrom(minEnd, slotEndHour, slotStepMin)
	return timesToKeyboard("meet:end:", slots)
}

func timesToKeyboard(prefix string, times []string) tgbotapi.InlineKeyboardMarkup {
	const perRow = 4
	var rows [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < len(times); i += perRow {
		end := i + perRow
		if end > len(times) { end = len(times) }
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
		// последний показываем 22:00, чтобы можно было завершать ровно в 22:00
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
	// сдвинем на кратность stepMin
	minute := (from.Minute()/stepMin + 0) * stepMin
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
	if err != nil { return fmt.Errorf("неверное время начала") }
	end, err := time.Parse("15:04", endHHMM)
	if err != nil { return fmt.Errorf("неверное время окончания") }

	open := time.Date(0,1,1,slotStartHour,0,0,0,time.UTC)
	close := time.Date(0,1,1,slotEndHour,0,0,0,time.UTC)

	s := time.Date(0,1,1,start.Hour(),start.Minute(),0,0,time.UTC)
	e := time.Date(0,1,1,end.Hour(),end.Minute(),0,0,time.UTC)

	if s.Before(open) || e.After(close) {
		return fmt.Errorf("время должно быть в диапазоне %02d:00–%02d:00", slotStartHour, slotEndHour)
	}
	if !e.After(s) {
		return fmt.Errorf("окончание должно быть позже начала")
	}
	// кратность шагу
	if start.Minute()%slotStepMin != 0 || end.Minute()%slotStepMin != 0 {
		return fmt.Errorf("используйте шаг %d минут", slotStepMin)
	}
	return nil
}
