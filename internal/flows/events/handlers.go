package flows

import (
	"bytes"
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

// ---------- список событий ----------

func intro(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	s.Step = EventsList
	return botengine.InternalContinue, nil
}

func list(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	// 1) загрузка
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

	// 2) фильтр
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

	// 3) сортировка и топ-5
	sort.Slice(filtered, func(i, j int) bool {
		di, _ := parseAnyEventDate(filtered[i].Date)
		dj, _ := parseAnyEventDate(filtered[j].Date)
		return di.Before(dj)
	})
	if len(filtered) > 5 {
		filtered = filtered[:5]
	}

	// 4) счётчики зарегистрированных
	counts, _ := fetchEntriesCounts(ctx, entriesUniqueURL)

	var sb strings.Builder
	sb.WriteString(p.Sprintf("events_intro"))
	sb.WriteString("\n\n")

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, e := range filtered {
		tm, _ := parseAnyEventDate(e.Date)
		date := tm.Format("02.01.2006")
		wd := weekdayShort(tm.Weekday(), s.Lang)

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
		used := counts[eventID(e)]
		left := 0
		if e.Capacity > 0 {
			left = e.Capacity - used
			if left < 0 {
				left = 0
			}
		}

		// текст блока
		sb.WriteString(fmt.Sprintf("• <b>%s</b> <i>(%s)</i> — <b>%s</b>\n", htmlEscape(date), htmlEscape(wd), htmlEscape(name)))
		if desc != "" {
			sb.WriteString(htmlEscape(desc))
			sb.WriteString("\n")
		}
		if e.Capacity > 0 {
			sb.WriteString(fmt.Sprintf("Места: %d/%d (осталось %d)\n", used, e.Capacity, left))
		}
		sb.WriteString(fmt.Sprintf("<a href=\"%s\">Подробнее →</a>\n\n", htmlEscape(url)))

		// кнопка регистрации / плейсхолдер
		if e.Capacity > 0 && left == 0 {
			rows = append(rows, ui.Row(ui.Cb("⛔ Мест нет", "noop")))
		} else {
			dt := dateShort(tm, s.Lang)
			base := fmt.Sprintf("📝 %s — %s", dt, name) // пример: 📝 21.08 18:00 — Pop-up Smoky BBQ
			lbl  := shortRunes(base, 60)               // Telegram ограничивает длину текста кнопки
			if e.Capacity > 0 {
				lbl = shortRunes(fmt.Sprintf("%s • %d", base, left), 60) // добавим «осталось N», если влезает
			}

			rows = append(rows, ui.Row(
				ui.Cb(lbl, "events:reg:"+eventID(e)),
			))
		}
	}

	// блок подписки
	if s.IsSubscribed {
		rows = append(rows,
			ui.Row(ui.Cb("📬 Подписка активна", "noop")),
			ui.Row(ui.Cb("⚙️ Изменить расписание", "events:edit"), ui.Cb("🛑 Отписаться", "events:unsubscribe")),
		)
	} else {
		rows = append(rows, ui.Row(ui.Cb("📬 Подписаться на еженедельные анонсы", "events:subscribe")))
	}

	kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	
	if err := ui.SendHTML(d.Bot, s.ChatID, sb.String(), kb); err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("[events] send error: %v", err))
	}
	s.Flow, s.Step = "", ""
	return EventsDone, nil
}

// ---------- подписка: мастер выбора дня/времени ----------

func subscribe(_ context.Context, ev botengine.Event, _ botengine.Deps, s *types.Session) (types.Step, error) {
	if ev.CallbackData == "events:subscribe" || strings.HasPrefix(ev.CallbackData, "events:sub:day:") {
		s.Step = EventsSubPickDay
		return botengine.InternalContinue, nil
	}
	return EventsSub, nil
}

func subPickDay(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
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
	txt := "Когда присылать анонсы? Выбери день недели:"
	kb := daysKB(s.Lang)
	if err := ui.SendHTML(d.Bot, s.ChatID, htmlEscape(txt), kb); err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("[events] subPickDay send error: %v", err))
	}
	return EventsSubPickDay, nil
}

func subPickTime(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
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
	txt := fmt.Sprintf("Ок, день: <b>%s</b>.\nТеперь выбери время отправки:",
		htmlEscape(weekdayHuman(time.Weekday(s.EventsSubDOW), s.Lang)))
	kb := timeKB()
	if err := ui.SendHTML(d.Bot, s.ChatID, txt, kb); err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("[events] subPickTime send error: %v", err))
	}
	return EventsSubPickTime, nil
}

func subAwaitTimeText(_ context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
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

func subConfirm(ctx context.Context, _ botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	loc := userLoc(s)
	next := computeNextRunUTC(s.EventsSubHour, s.EventsSubMinute, time.Weekday(s.EventsSubDOW), loc)
	s.IsSubscribed = true
	s.NextDigestAt = next

	if _, err := list(ctx, botengine.Event{}, d, s); err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("[events] preview send error: %v", err))
	}

	msg := fmt.Sprintf(
		"Готово! Будем присылать анонсы каждую <b>%s</b> в <b>%02d:%02d</b> (%s).\n"+
			"Следующая отправка по расписанию: <i>%s</i>.\n\n"+
			"Чтобы изменить расписание — /events_time, чтобы отписаться — /unsubscribe_events.",
		htmlEscape(weekdayHuman(time.Weekday(s.EventsSubDOW), s.Lang)),
		s.EventsSubHour, s.EventsSubMinute, loc.String(),
		next.In(loc).Format("02.01.2006 15:04"),
	)
	kb := ui.Inline(
		ui.Row(ui.Cb("⚙️ Изменить расписание", "events:edit"), ui.Cb("🛑 Отписаться", "events:unsubscribe")),
	)
	if err := ui.SendHTML(d.Bot, s.ChatID, msg, kb); err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("[events] subConfirm send error: %v", err))
	}
	s.Flow, s.Step = "", ""
	return EventsDone, nil
}

// ---------- изменение расписания / отписка ----------

func editSchedule(_ context.Context, _ botengine.Event, _ botengine.Deps, s *types.Session) (types.Step, error) {
	s.Step = EventsSubPickDay
	return botengine.InternalContinue, nil
}

func unsubscribe(_ context.Context, _ botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	s.IsSubscribed = false
	_ = ui.SendText(d.Bot, s.ChatID, "Вы отписаны от еженедельных анонсов. Мы не обиделись — просто будем скучать 🐈‍⬛")
	s.Flow, s.Step = "", ""
	return EventsDone, nil
}

// ---------- helpers общие ----------

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
	r := strings.NewReplacer(`&`, "&amp;", `<`, "&lt;", `>`, "&gt;", `"`, "&quot;", `'`, "&#39;")
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
		ui.Row(ui.Cb(lbl(time.Monday), "events:sub:day:1"), ui.Cb(lbl(time.Tuesday), "events:sub:day:2"), ui.Cb(lbl(time.Wednesday), "events:sub:day:3")),
		ui.Row(ui.Cb(lbl(time.Thursday), "events:sub:day:4"), ui.Cb(lbl(time.Friday), "events:sub:day:5"), ui.Cb(lbl(time.Saturday), "events:sub:day:6")),
		ui.Row(ui.Cb(lbl(time.Sunday), "events:sub:day:0")),
	)
}

func timeKB() tgbotapi.InlineKeyboardMarkup {
	return ui.Inline(
		ui.Row(ui.Cb("09:00", "events:sub:time:09:00"), ui.Cb("12:00", "events:sub:time:12:00"), ui.Cb("15:00", "events:sub:time:15:00")),
		ui.Row(ui.Cb("18:00", "events:sub:time:18:00"), ui.Cb("21:00", "events:sub:time:21:00"), ui.Cb("Другое…", "events:sub:time:custom")),
	)
}

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

func handleListCallbacks(_ context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	switch ev.CallbackData {
	case "events:edit":
		s.Step = EventsSubPickDay
		return botengine.InternalContinue, nil
	case "events:unsubscribe":
		return unsubscribe(context.Background(), ev, d, s)
	default:
		return EventsDone, nil
	}
}

func done(_ context.Context, _ botengine.Event, _ botengine.Deps, _ *types.Session) (types.Step, error) {
	return EventsDone, nil
}

// --- entries counters ---

func fetchEntriesCounts(ctx context.Context, url string) (map[string]int, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var m map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

// ---------- Регистрация: шаги ----------

func ackCB(d botengine.Deps, ev botengine.Event) {
	if ev.CallbackQueryID == "" {
		return
	}
	_, _ = d.Bot.Request(tgbotapi.NewCallback(ev.CallbackQueryID, ""))
}

func regStart(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
    ackCB(d, ev)
    if !strings.HasPrefix(ev.CallbackData, "events:reg:") {
        return EventsDone, nil
    }
    id := strings.TrimPrefix(ev.CallbackData, "events:reg:")
    if id == "" {
        _ = ui.SendText(d.Bot, s.ChatID, "Не удалось распознать мероприятие. Попробуй ещё раз 🙏")
        return EventsDone, nil
    }

    var e *types.Event
    if d.Svcs.Events != nil {
        if events, _ := d.Svcs.Events.ListUpcoming(ctx); len(events) > 0 {
            for i := range events {
                if eventID(events[i]) == id {
                    e = &events[i]
                    break
                }
            }
        }
    }
    var t time.Time
    if e != nil {
        t, _ = parseAnyEventDate(e.Date)
        _ = stSet(ctx, d, s.ChatID, keyRegCapacity, strconv.Itoa(e.Capacity))
    }
    _ = stSet(ctx, d, s.ChatID, keyRegEventID, id)
    if !t.IsZero() {
        _ = stSet(ctx, d, s.ChatID, keyRegEventDate, t.UTC().Format(time.RFC3339))
    } else {
        _ = stDel(ctx, d, s.ChatID, keyRegEventDate)
    }
    _ = stSet(ctx, d, s.ChatID, keyRegGuests, "1")

    // ранняя проверка переполнения
    if e != nil && e.Capacity > 0 {
        if counts, err := fetchEntriesCounts(ctx, entriesUniqueURL); err == nil {
            used := counts[id]
            if used >= e.Capacity {
                _ = ui.SendText(d.Bot, s.ChatID, "Уф… мест уже нет на это событие 😿 Посмотри другие через /events.")
                return EventsDone, nil
            }
        }
    }

    // 👇 ВСТАВЬ ЭТОТ БЛОК: заголовок «на что регистрируемся»
    if e != nil {
        tstr := dateShort(t, s.Lang)
        header := fmt.Sprintf("Регистрация: %s — %s", strings.TrimSpace(e.Name), tstr)
        _ = ui.SendText(d.Bot, s.ChatID, header)
    }

    _ = ui.SendText(d.Bot, s.ChatID, "Начнём регистрацию.\nКак к вам обращаться? (Имя обязательно)")
    return EventsRegAskName, nil
}


func regAskName(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	txt := strings.TrimSpace(ev.Text)
	if txt == "" {
		_ = ui.SendText(d.Bot, s.ChatID, "Имя — обязательное поле 🙏 Введите имя:")
		return EventsRegAskName, nil
	}
	if len([]rune(txt)) < 2 {
		_ = ui.SendText(d.Bot, s.ChatID, "Слишком короткое имя. Давай хотя бы 2 буквы 😊")
		return EventsRegAskName, nil
	}
	_ = stSet(ctx, d, s.ChatID, keyRegName, txt)
	_ = ui.SendText(d.Bot, s.ChatID, "Отлично! Теперь email (мы пришлём подтверждение).")
	return EventsRegAskEmail, nil
}

func regAskEmail(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	txt := strings.TrimSpace(ev.Text)
	if !reEmail.MatchString(txt) {
		_ = ui.SendText(d.Bot, s.ChatID, "Похоже, это не похоже на email 🙂 Введите, пожалуйста, корректный e-mail:")
		return EventsRegAskEmail, nil
	}
	_ = stSet(ctx, d, s.ChatID, keyRegEmail, txt)
	_ = ui.SendText(d.Bot, s.ChatID, "Телефон (только цифры, можно с +):")
	return EventsRegAskPhone, nil
}

func regAskPhone(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	txt := strings.ReplaceAll(strings.TrimSpace(ev.Text), " ", "")
	if !rePhone.MatchString(txt) {
		_ = ui.SendText(d.Bot, s.ChatID, "Телефон обязателен. Введите номер (7–15 цифр, можно с +):")
		return EventsRegAskPhone, nil
	}
	_ = stSet(ctx, d, s.ChatID, keyRegPhone, txt)
	_ = ui.SendText(d.Bot, s.ChatID, "Сколько гостей придёт? (число, по умолчанию 1)")
	return EventsRegAskGuests, nil
}

func regAskGuests(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	txt := strings.TrimSpace(ev.Text)
	if txt == "" {
		txt = "1"
	}
	n, err := strconv.Atoi(txt)
	if err != nil || n < 1 || n > 20 {
		_ = ui.SendText(d.Bot, s.ChatID, "Введите число от 1 до 20:")
		return EventsRegAskGuests, nil
	}

	// если известен capacity — сверяемся с остатком
	if capStr, ok := stGet(ctx, d, s.ChatID, keyRegCapacity); ok {
		if capVal, _ := strconv.Atoi(capStr); capVal > 0 {
			if evID, ok2 := stGet(ctx, d, s.ChatID, keyRegEventID); ok2 {
				if counts, err := fetchEntriesCounts(ctx, entriesUniqueURL); err == nil {
					used := counts[evID]
					left := capVal - used
					if left < 0 {
						left = 0
					}
					if n > left {
						if left == 0 {
							_ = ui.SendText(d.Bot, s.ChatID, "К сожалению, места уже закончились 😿 Выберите другое событие: /events")
							return EventsDone, nil
						}
						_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("Осталось только %d мест(а). Введите число не больше %d:", left, left))
						return EventsRegAskGuests, nil
					}
				}
			}
		}
	}

	_ = stSet(ctx, d, s.ChatID, keyRegGuests, strconv.Itoa(n))

	handle := ev.FromUserName
	if strings.TrimSpace(handle) == "" {
		handle = "@" + strconv.FormatInt(ev.FromUserID, 10)
	}
	_ = stSet(ctx, d, s.ChatID, keyRegTelegram, handle)
	_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("Укажите Telegram (или оставьте как есть):\n%s", handle))
	return EventsRegAskTelegram, nil
}

func regAskTelegram(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	txt := strings.TrimSpace(ev.Text)
	if txt != "" {
		_ = stSet(ctx, d, s.ChatID, keyRegTelegram, txt)
	}
	_ = ui.SendText(d.Bot, s.ChatID, "Комментарий (необязательно). Если нечего добавить — отправьте «-».")
	return EventsRegAskComment, nil
}

func regAskComment(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	txt := strings.TrimSpace(ev.Text)
	if txt == "-" {
		txt = ""
	}
	_ = stSet(ctx, d, s.ChatID, keyRegComment, txt)

	// подтверждение
	name, _ := stGet(ctx, d, s.ChatID, keyRegName)
	email, _ := stGet(ctx, d, s.ChatID, keyRegEmail)
	phone, _ := stGet(ctx, d, s.ChatID, keyRegPhone)
	guests, _ := stGet(ctx, d, s.ChatID, keyRegGuests)
	tg, _ := stGet(ctx, d, s.ChatID, keyRegTelegram)

	dateStr := humanEventDate(ctx, d, s)
	summary := fmt.Sprintf(
		"Проверьте данные:\n\nИмя: <b>%s</b>\nEmail: <b>%s</b>\nТелефон: <b>%s</b>\nГостей: <b>%s</b>\nTelegram: <b>%s</b>\nДата: <b>%s</b>\n",
		htmlEscape(name), htmlEscape(email), htmlEscape(phone), htmlEscape(guests), htmlEscape(tg), htmlEscape(dateStr),
	)
	kb := ui.Inline(
		ui.Row(ui.Cb("✅ Подтвердить", "events:reg:confirm"), ui.Cb("✏️ Исправить имя", "events:reg:edit:name")),
		ui.Row(ui.Cb("❌ Отменить регистрацию", "events:reg_cancel")),
	)
	if err := ui.SendHTML(d.Bot, s.ChatID, summary, kb); err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, "Не удалось отправить подтверждение, попробуйте ещё раз.")
	}
	return EventsRegConfirm, nil
}

func regConfirm(_ context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ackCB(d, ev)
	switch ev.CallbackData {
	case "events:reg:confirm":
		return EventsRegSubmit, nil
	case "events:reg:edit:name":
		_ = ui.SendText(d.Bot, s.ChatID, "Введите имя заново:")
		return EventsRegAskName, nil
	default:
		return EventsRegConfirm, nil
	}
}

func regSubmit(ctx context.Context, _ botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	name, _ := stGet(ctx, d, s.ChatID, keyRegName)
	email, _ := stGet(ctx, d, s.ChatID, keyRegEmail)
	phone, _ := stGet(ctx, d, s.ChatID, keyRegPhone)
	guestsStr, _ := stGet(ctx, d, s.ChatID, keyRegGuests)
	tg, _ := stGet(ctx, d, s.ChatID, keyRegTelegram)
	comment, _ := stGet(ctx, d, s.ChatID, keyRegComment)
	eventID, _ := stGet(ctx, d, s.ChatID, keyRegEventID)
	dateHuman := humanEventDate(ctx, d, s)

	// обязательные
	if name == "" || !reEmail.MatchString(email) || !rePhone.MatchString(phone) {
		_ = ui.SendText(d.Bot, s.ChatID, "Кажется, не все обязательные поля заполнены. Давай начнём заново: /events")
		return EventsDone, nil
	}

	// финальный double-check capacity
	if capStr, ok := stGet(ctx, d, s.ChatID, keyRegCapacity); ok {
		if capVal, _ := strconv.Atoi(capStr); capVal > 0 {
			if counts, err := fetchEntriesCounts(ctx, entriesUniqueURL); err == nil {
				used := counts[eventID]
				left := capVal - used
				if left < 0 {
					left = 0
				}
				need, _ := strconv.Atoi(guestsStr)
				if need > left {
					_ = ui.SendText(d.Bot, s.ChatID, "Пока мы заполняли форму, места закончились 😿\nПопробуйте другое событие — /events")
					return EventsDone, nil
				}
			}
		}
	}

	// POST
	body := regPayload{
		Name:            name,
		Email:           email,
		Phone:           phone,
		NumberOfPersons: guestsStr,
		Telegram:        tg,
		Date:            dateHuman,
		EventID:         eventID,
		Comment:         comment,
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, registrationEndpoint, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, "Не удалось отправить форму, попробуйте ещё раз чуть позже 🙏")
		return EventsDone, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("Сервер ответил %d. Попробуйте позже или напишите нам сюда.", resp.StatusCode))
		return EventsDone, nil
	}

	// финалка
	text := "Спасибо за регистрацию! 🎉\n\n" +
		"Пожалуйста, не закрывайте и не удаляйте бота — иначе мы не сможем прислать напоминание и важные детали мероприятия.\n" +
		"Если что-то изменится — просто напишите нам сюда в чат.\n\n" +
		"До встречи!"
	kb := ui.Inline(ui.Row(ui.Cb("❌ Отменить регистрацию", "events:reg_cancel")))
	_ = ui.SendHTML(d.Bot, s.ChatID, htmlEscape(text), kb)

	// можно выставить напоминание (если нужна фоновая джоба — установите её тут)
	// _ = stSet(ctx, d, s.ChatID, keyRegReminderAt, time.Now().UTC().Format(time.RFC3339))

	s.Flow, s.Step = "", ""
	return EventsDone, nil
}

// --- человекочитаемая дата ---

func humanEventDate(ctx context.Context, d botengine.Deps, s *types.Session) string {
	if raw, ok := stGet(ctx, d, s.ChatID, keyRegEventDate); ok {
		if t, err := time.Parse(time.RFC3339, raw); err == nil {
			return formatRuHuman(t)
		}
	}
	return "дата будет уточнена"
}

func formatRuHuman(t time.Time) string {
	w := map[time.Weekday]string{
		time.Monday: "пн", time.Tuesday: "вт", time.Wednesday: "ср",
		time.Thursday: "чт", time.Friday: "пт", time.Saturday: "сб", time.Sunday: "вс",
	}[t.Weekday()]
	months := []string{"января", "февраля", "марта", "апреля", "мая", "июня", "июля", "августа", "сентября", "октября", "ноября", "декабря"}
	return fmt.Sprintf("%s %d %s. %02d:%02d", w, t.Day(), months[int(t.Month())-1], t.Hour(), t.Minute())
}

func regCancelAsk(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ackCB(d, ev)
	msg := "Ой… Нам очень жаль 😿 Мы готовимся к каждому гостю и бережём места.\n" +
		"Точно отменяем? (можно просто прийти на другое событие — мы будем рады!)"
	kb := ui.Inline(
		ui.Row(ui.Cb("Да, отменить", "events:reg_cancel:yes"), ui.Cb("Оставить регистрацию", "events:reg_cancel:no")),
	)
	_ = ui.SendHTML(d.Bot, s.ChatID, htmlEscape(msg), kb)
	return EventsRegCancelDo, nil
}

func regCancelDo(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ackCB(d, ev)
	switch ev.CallbackData {
	case "events:reg_cancel:yes":
		// помечаем как отменённую локально; при появлении backend-эндпойнта — дернуть его здесь
		_ = stDel(ctx, d, s.ChatID, keyRegReminderAt)
		_ = ui.SendText(d.Bot, s.ChatID, "Окей, мы отметили отмену. Если передумаете — снова жмякните /events ❤️")
	case "events:reg_cancel:no":
		_ = ui.SendText(d.Bot, s.ChatID, "Ура! Мы вас ждём 🥳")
	}
	s.Flow, s.Step = "", ""
	return EventsDone, nil
}

// dd.mm [HH:MM] (RU) / 02 Jan [15:04] (EN)
func dateShort(t time.Time, lang string) string {
	lang = strings.ToLower(lang)
	if strings.HasPrefix(lang, "ru") {
		if t.Hour() == 0 && t.Minute() == 0 {
			return fmt.Sprintf("%02d.%02d", t.Day(), int(t.Month()))
		}
		return fmt.Sprintf("%02d.%02d %02d:%02d", t.Day(), int(t.Month()), t.Hour(), t.Minute())
	}
	if t.Hour() == 0 && t.Minute() == 0 {
		return t.Format("02 Jan")
	}
	return t.Format("02 Jan 15:04")
}

func shortRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n { return s }
	return string(r[:n-1]) + "…"
}
