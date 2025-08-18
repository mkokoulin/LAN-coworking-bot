package flows

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

func intro(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	s.Step = EventsList
	return botengine.InternalContinue, nil
}

func list(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	// 1) Загрузка
	var items []types.Event
	var err error
	if d.Svcs.Events != nil {
		items, err = d.Svcs.Events.ListUpcoming(ctx)
	} else {
		items, err = fetchEventsFallback(ctx, "https://shark-app-wrcei.ondigitalocean.app/api/events")
	}
	if err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("[events] load error: %v", err))
		s.Flow, s.Step = "", ""
		return EventsDone, nil
	}
	if len(items) == 0 {
		_ = ui.SendText(d.Bot, s.ChatID, "[events] no items")
		s.Flow, s.Step = "", ""
		return EventsDone, nil
	}

	// 2) Фильтр
	filtered := make([]types.Event, 0, len(items))
	for _, e := range items {
		if !e.ShowForm && hasShowFormField(e) {
			continue
		}
		if _, err := parseAnyEventDate(e.Date); err != nil {
			continue
		}
		filtered = append(filtered, e)
	}
	if len(filtered) == 0 {
		_ = ui.SendText(d.Bot, s.ChatID, "[events] nothing after filter (date/showForm)")
		s.Flow, s.Step = "", ""
		return EventsDone, nil
	}

	// 3) Сортировка + top-5
	sort.Slice(filtered, func(i, j int) bool {
		di, _ := parseAnyEventDate(filtered[i].Date)
		dj, _ := parseAnyEventDate(filtered[j].Date)
		return di.Before(dj)
	})
	if len(filtered) > 5 {
		filtered = filtered[:5]
	}

	var sb strings.Builder
	sb.WriteString(p.Sprintf("events_intro"))
	sb.WriteString("\n\n")

	for _, e := range filtered {
		t, _ := parseAnyEventDate(e.Date)
		date := t.Format("02.01.2006")
		wd := weekdayShort(t.Weekday(), s.Lang)

		name := strings.TrimSpace(e.Name)
		if name == "" {
			name = "Untitled"
		}
		if len([]rune(name)) > 120 {
			name = string([]rune(name)[:117]) + "…"
		}

		desc := strings.TrimSpace(stripHTML(e.Description))
		if len([]rune(desc)) > 200 {
			desc = string([]rune(desc)[:197]) + "…"
		}

		url := fmt.Sprintf("https://lettersandnumbers.am/events/%s", eventID(e))

		sb.WriteString(fmt.Sprintf("• <b>%s</b> <i>(%s)</i> — <b>%s</b>\n",
			htmlEscape(date), htmlEscape(wd), htmlEscape(name)))
		if desc != "" {
			sb.WriteString(htmlEscape(desc))
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("<a href=\"%s\">Registration →</a>\n\n", htmlEscape(url)))
	}

	// одна отправка: HTML + inline-кнопки
	kb := ui.Inline(
		ui.Row(
			ui.Cb("📬 Подписаться на еженедельные анонсы", "events:subscribe"),
		),
	)
	// если уже подписан — сразу покажем кнопку редактирования
	if s.IsSubscribed {
		kb = ui.Inline(
			ui.Row(
				ui.Cb("📬 Подписка активна", "noop"),
			),
			ui.Row(
				ui.Cb("⚙️ Изменить расписание", "events:edit"),
				ui.Cb("🛑 Отписаться", "events:unsubscribe"),
			),
		)
	}

	if err := ui.SendHTML(d.Bot, s.ChatID, sb.String(), kb); err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("[events] send error: %v", err))
	}

	s.Flow, s.Step = "", ""
	return EventsDone, nil
}

// ---------- подписка: мастер выбора дня/времени ----------

func subscribe(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if ev.CallbackData == "events:subscribe" {
		s.Step = EventsSubPickDay
		return botengine.InternalContinue, nil
	}
	// fallback на случай широкого роутинга
	if strings.HasPrefix(ev.CallbackData, "events:sub:day:") {
		s.Step = EventsSubPickDay
		return botengine.InternalContinue, nil
	}
	return EventsSub, nil
}

func subPickDay(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	// обработка выбора дня
	if strings.HasPrefix(ev.CallbackData, "events:sub:day:") {
		part := strings.TrimPrefix(ev.CallbackData, "events:sub:day:")
		wd, _ := strconv.Atoi(part)
		if wd < 0 || wd > 6 {
			_ = ui.SendText(d.Bot, s.ChatID, "Хмм, день недели не распознан. Попробуем ещё раз?")
			return EventsSubPickDay, nil
		}
		s.EventsSubDOW = wd
		s.Step = EventsSubPickTime
		return botengine.InternalContinue, nil
	}

	// показать клавиатуру с днями
	txt := "Когда присылать анонсы? Выбери день недели:"
	kb := daysKB(s.Lang)
	if err := ui.SendHTML(d.Bot, s.ChatID, htmlEscape(txt), kb); err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("[events] subPickDay send error: %v", err))
	}
	return EventsSubPickDay, nil
}

func subPickTime(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	// обработка готовых слотов
	if strings.HasPrefix(ev.CallbackData, "events:sub:time:") {
		val := strings.TrimPrefix(ev.CallbackData, "events:sub:time:")
		if val == "custom" {
			_ = ui.SendText(d.Bot, s.ChatID, "Введи время в формате HH:MM (например, 10:30)")
			return EventsSubAwaitInput, nil
		}
		hh, mm, ok := parseHHMM(val)
		if !ok {
			_ = ui.SendText(d.Bot, s.ChatID, "Не понял время. Давай ещё раз?")
			return EventsSubPickTime, nil
		}
		s.EventsSubHour, s.EventsSubMinute = hh, mm
		s.Step = EventsSubConfirm
		return botengine.InternalContinue, nil
	}

	// показать пресеты
	txt := fmt.Sprintf("Ок, день: <b>%s</b>.\nТеперь выбери время отправки:",
		htmlEscape(weekdayHuman(time.Weekday(s.EventsSubDOW), s.Lang)))
	kb := timeKB()
	if err := ui.SendHTML(d.Bot, s.ChatID, txt, kb); err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("[events] subPickTime send error: %v", err))
	}
	return EventsSubPickTime, nil
}

func subAwaitTimeText(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	text := strings.TrimSpace(ev.Text)
	if text == "" {
		return EventsSubAwaitInput, nil
	}
	hh, mm, ok := parseHHMM(text)
	if !ok {
		_ = ui.SendText(d.Bot, s.ChatID, "Формат времени — HH:MM (00–23:00–59). Попробуй ещё раз 🙏")
		return EventsSubAwaitInput, nil
	}
	s.EventsSubHour, s.EventsSubMinute = hh, mm
	s.Step = EventsSubConfirm
	return botengine.InternalContinue, nil
}

func subConfirm(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	loc := userLoc(s)

	// 1) Считаем ближайший запуск и сохраняем в сессию
	next := computeNextRunUTC(s.EventsSubHour, s.EventsSubMinute, time.Weekday(s.EventsSubDOW), loc)
	s.IsSubscribed = true
	s.NextDigestAt = next // сохраняем UTC-дату для джобы

	// 2) Шлём превью СЕЙЧАС (чтобы пользователь получил список сразу)
	if _, err := list(ctx, ev, d, s); err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("[events] preview send error: %v", err))
	}

	// 3) Сообщение-подтверждение (с датой следующей отправки)
	msg := fmt.Sprintf(
		"Готово! Будем присылать анонсы каждую <b>%s</b> в <b>%02d:%02d</b> (%s).\n"+
			"Следующая отправка по расписанию: <i>%s</i>.\n\n"+
			"Чтобы изменить расписание — /events_time, чтобы отписаться — /unsubscribe_events.",
		htmlEscape(weekdayHuman(time.Weekday(s.EventsSubDOW), s.Lang)),
		s.EventsSubHour, s.EventsSubMinute, loc.String(),
		next.In(loc).Format("02.01.2006 15:04"),
	)

	kb := ui.Inline(
		ui.Row(
			ui.Cb("⚙️ Изменить расписание", "events:edit"),
			ui.Cb("🛑 Отписаться", "events:unsubscribe"),
		),
	)
	if err := ui.SendHTML(d.Bot, s.ChatID, msg, kb); err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("[events] subConfirm send error: %v", err))
	}

	// 4) Закрываем флоу
	s.Flow, s.Step = "", ""
	return EventsDone, nil
}

// ---------- изменение расписания ----------

func editSchedule(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	// вызывается по коллбэку "events:edit" или командой /events_time
	_ = ev // неважно, откуда пришли — ведём в выбор дня
	s.Step = EventsSubPickDay
	return botengine.InternalContinue, nil
}

// ---------- отписка ----------

func unsubscribe(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	s.IsSubscribed = false
	_ = ui.SendText(d.Bot, s.ChatID, "Вы отписаны от еженедельных анонсов. Мы не обиделись — просто будем скучать 🐈‍⬛")
	s.Flow, s.Step = "", ""
	return EventsDone, nil
}

// ---------- helpers (как у тебя + новые) ----------

func parseAnyEventDate(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	for _, f := range []string{"2006-01-02", "02.01.2006"} {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized date: %q", s)
}

func weekdayShort(w time.Weekday, lang string) string {
	if strings.HasPrefix(strings.ToLower(lang), "ru") {
		switch w {
		case time.Monday:
			return "Пн"
		case time.Tuesday:
			return "Вт"
		case time.Wednesday:
			return "Ср"
		case time.Thursday:
			return "Чт"
		case time.Friday:
			return "Пт"
		case time.Saturday:
			return "Сб"
		default:
			return "Вс"
		}
	}
	switch w {
	case time.Monday:
		return "Mon"
	case time.Tuesday:
		return "Tue"
	case time.Wednesday:
		return "Wed"
	case time.Thursday:
		return "Thu"
	case time.Friday:
		return "Fri"
	case time.Saturday:
		return "Sat"
	default:
		return "Sun"
	}
}

func weekdayHuman(w time.Weekday, lang string) string {
	if strings.HasPrefix(strings.ToLower(lang), "ru") {
		switch w {
		case time.Monday:
			return "понедельник"
		case time.Tuesday:
			return "вторник"
		case time.Wednesday:
			return "среду"
		case time.Thursday:
			return "четверг"
		case time.Friday:
			return "пятницу"
		case time.Saturday:
			return "субботу"
		default:
			return "воскресенье"
		}
	}
	switch w {
	case time.Monday:
		return "Monday"
	case time.Tuesday:
		return "Tuesday"
	case time.Wednesday:
		return "Wednesday"
	case time.Thursday:
		return "Thursday"
	case time.Friday:
		return "Friday"
	case time.Saturday:
		return "Saturday"
	default:
		return "Sunday"
	}
}

func stripHTML(input string) string {
	re := regexp.MustCompile(`<.*?>`)
	return strings.TrimSpace(re.ReplaceAllString(input, ""))
}

func htmlEscape(s string) string {
	r := strings.NewReplacer(
		`&`, "&amp;",
		`<`, "&lt;",
		`>`, "&gt;",
		`"`, "&quot;",
		`'`, "&#39;",
	)
	return r.Replace(s)
}

func eventID(e types.Event) string { return e.ID }
func hasShowFormField(_ types.Event) bool { return true }

func fetchEventsFallback(ctx context.Context, url string) ([]types.Event, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw []types.Event
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// --- UI клавиатуры ---

func daysKB(lang string) tgbotapi.InlineKeyboardMarkup {
	// time.Weekday: 0=Sunday ... 6=Saturday
	lbl := func(w time.Weekday) string {
		if strings.HasPrefix(strings.ToLower(lang), "ru") {
			switch w {
			case time.Monday:
				return "Пн"
			case time.Tuesday:
				return "Вт"
			case time.Wednesday:
				return "Ср"
			case time.Thursday:
				return "Чт"
			case time.Friday:
				return "Пт"
			case time.Saturday:
				return "Сб"
			default:
				return "Вс"
			}
		}
		return weekdayShort(w, lang)
	}
	return ui.Inline(
		ui.Row(ui.Cb(lbl(time.Monday), "events:sub:day:1"),
			ui.Cb(lbl(time.Tuesday), "events:sub:day:2"),
			ui.Cb(lbl(time.Wednesday), "events:sub:day:3")),
		ui.Row(ui.Cb(lbl(time.Thursday), "events:sub:day:4"),
			ui.Cb(lbl(time.Friday), "events:sub:day:5"),
			ui.Cb(lbl(time.Saturday), "events:sub:day:6")),
		ui.Row(ui.Cb(lbl(time.Sunday), "events:sub:day:0")),
	)
}

func timeKB() tgbotapi.InlineKeyboardMarkup {
	return ui.Inline(
		ui.Row(ui.Cb("09:00", "events:sub:time:09:00"),
			ui.Cb("12:00", "events:sub:time:12:00"),
			ui.Cb("15:00", "events:sub:time:15:00")),
		ui.Row(ui.Cb("18:00", "events:sub:time:18:00"),
			ui.Cb("21:00", "events:sub:time:21:00"),
			ui.Cb("Другое…", "events:sub:time:custom")),
	)
}

var reHHMM = regexp.MustCompile(`^(?:[01]?\d|2[0-3]):[0-5]\d$`)

func parseHHMM(s string) (int, int, bool) {
	s = strings.TrimSpace(s)
	if !reHHMM.MatchString(s) {
		return 0, 0, false
	}
	parts := strings.SplitN(s, ":", 2)
	hh, _ := strconv.Atoi(parts[0])
	mm, _ := strconv.Atoi(parts[1])
	return hh, mm, true
}

// --- расчёт следующего запуска ---

func userLoc(_ *types.Session) *time.Location {
	// Если у тебя есть поле с TZ — используй его.
	// По умолчанию — Ереван.
	loc, err := time.LoadLocation("Asia/Yerevan")
	if err != nil {
		return time.FixedZone("Asia/Yerevan", 4*3600)
	}
	return loc
}

func computeNextRunUTC(hh, mm int, dow time.Weekday, loc *time.Location) time.Time {
	now := time.Now().In(loc)
	shift := (int(dow) - int(now.Weekday()) + 7) % 7
	cand := time.Date(now.Year(), now.Month(), now.Day(), hh, mm, 0, 0, loc).AddDate(0, 0, shift)
	if !cand.After(now) {
		cand = cand.AddDate(0, 0, 7)
	}
	return cand.UTC()
}

// ---------- edit & unsubscribe callbacks from list() ----------

func handleListCallbacks(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	switch ev.CallbackData {
	case "events:edit":
		s.Step = EventsSubPickDay
		return botengine.InternalContinue, nil
	case "events:unsubscribe":
		return unsubscribe(ctx, ev, d, s)
	default:
		return EventsDone, nil
	}
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return EventsDone, nil
}
