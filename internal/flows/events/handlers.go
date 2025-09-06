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
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

// ====================== in-memory KV ======================

var memKV = struct {
	data map[int64]map[string]string
}{data: make(map[int64]map[string]string)}

func stSet(_ context.Context, _ botengine.Deps, chatID int64, key, val string) error {
	if _, ok := memKV.data[chatID]; !ok {
		memKV.data[chatID] = make(map[string]string)
	}
	memKV.data[chatID][key] = val
	return nil
}
func stGet(_ context.Context, _ botengine.Deps, chatID int64, key string) (string, bool) {
	if m, ok := memKV.data[chatID]; ok {
		v, ok2 := m[key]
		return v, ok2
	}
	return "", false
}
func stDel(_ context.Context, _ botengine.Deps, chatID int64, key string) error {
	if m, ok := memKV.data[chatID]; ok {
		delete(m, key)
	}
	return nil
}

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

// ====================== Список событий ======================

func intro(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	s.Step = EventsList
	return botengine.InternalContinue, nil
}

func list(ctx context.Context, _ botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	// загрузка
	var items []types.Event
	var err error
	if d.Svcs.Events != nil {
		items, err = d.Svcs.Events.ListUpcoming(ctx)
	} else {
		items, err = fetchEventsFallback(ctx, eventsURLFallback)
	}
	if err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("[events] load error: %v", err))
		s.Flow, s.Step = "", ""
		return EventsDone, nil
	}

	// фильтр/сорт/топ-5
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
	sort.Slice(filtered, func(i, j int) bool {
		di, _ := parseAnyEventDate(filtered[i].Date)
		dj, _ := parseAnyEventDate(filtered[j].Date)
		return di.Before(dj)
	})
	if len(filtered) > 5 {
		filtered = filtered[:5]
	}

	// счётчики
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

		sb.WriteString(fmt.Sprintf("• <b>%s</b> <i>(%s)</i> — <b>%s</b>\n", htmlEscape(date), htmlEscape(wd), htmlEscape(name)))
		if desc != "" {
			sb.WriteString(htmlEscape(desc))
			sb.WriteString("\n")
		}
		if e.Capacity > 0 {
			sb.WriteString(fmt.Sprintf("Места: %d/%d (осталось %d)\n", used, e.Capacity, left))
		}
		sb.WriteString(fmt.Sprintf("<a href=\"%s\">Подробнее →</a>\n\n", htmlEscape(url)))

		if e.Capacity > 0 && left == 0 {
			rows = append(rows, ui.Row(ui.Cb("⛔ Мест нет", "noop")))
		} else {
			dt := dateShort(tm, s.Lang)
			base := fmt.Sprintf("📝 %s — %s", dt, name)
			lbl := shortRunes(base, 60)
			if e.Capacity > 0 {
				lbl = shortRunes(fmt.Sprintf("%s • %d", base, left), 60)
			}
			rows = append(rows, ui.Row(ui.Cb(lbl, "events:regstart:"+eventID(e))))
		}
	}

	// подписка
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

// ====================== Подписка ======================

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

// правка расписания / отписка
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

// ====================== Регистрация с профилем ======================

func ackCB(d botengine.Deps, ev botengine.Event) {
	if ev.CallbackQueryID == "" {
		return
	}
	_, _ = d.Bot.Request(tgbotapi.NewCallback(ev.CallbackQueryID, ""))
}

func profGet(s *types.Session, key string) (string, bool) {
	if s == nil || s.Data == nil {
		return "", false
	}
	if v, ok := s.Data[key]; ok {
		if str, ok2 := v.(string); ok2 {
			return str, true
		}
	}
	return "", false
}
func profSet(ctx context.Context, d botengine.Deps, s *types.Session, key, val string) {
	if s.Data == nil {
		s.Data = map[string]interface{}{}
	}
	s.Data[key] = val
	if d.State != nil {
		d.State.Set(s.ChatID, s)
	}
}
func profileComplete(s *types.Session) bool {
	name, _ := profGet(s, keyProfName)
	email, _ := profGet(s, keyProfEmail)
	phone, _ := profGet(s, keyProfPhone)
	return strings.TrimSpace(name) != "" && reEmail.MatchString(email) && rePhone.MatchString(phone)
}

func regStart(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ackCB(d, ev)
	if !strings.HasPrefix(ev.CallbackData, "events:regstart:") {
		return EventsDone, nil
	}
	id := strings.TrimPrefix(ev.CallbackData, "events:regstart:")
	if id == "" {
		_ = ui.SendText(d.Bot, s.ChatID, "Не удалось распознать мероприятие. Попробуй ещё раз 🙏")
		return EventsDone, nil
	}

	// найдём событие (для даты/вместимости)
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
	_ = stDel(ctx, d, s.ChatID, keyRegComment)

	// Ранняя проверка переполнения
	if e != nil && e.Capacity > 0 {
		if counts, err := fetchEntriesCounts(ctx, entriesUniqueURL); err == nil {
			used := counts[id]
			if used >= e.Capacity {
				_ = ui.SendText(d.Bot, s.ChatID, "Уф… мест уже нет на это событие 😿 Посмотри другие через /events.")
				return EventsDone, nil
			}
		}
	}

	// Заголовок блока
	if e != nil {
		tstr := dateShort(t, s.Lang)
		header := fmt.Sprintf("Регистрация: %s — %s", strings.TrimSpace(e.Name), tstr)
		_ = ui.SendText(d.Bot, s.ChatID, header)
	}

	// Если профиль уже есть — спрашиваем только гостей/комментарий
	if profileComplete(s) {
		_ = ui.SendText(d.Bot, s.ChatID, "Сколько гостей придёт? (число, по умолчанию 1)")
		return EventsRegAskGuests, nil
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
	profSet(ctx, d, s, keyProfName, txt)
	_ = ui.SendText(d.Bot, s.ChatID, "Отлично! Теперь email (мы пришлём подтверждение).")
	return EventsRegAskEmail, nil
}

func regAskEmail(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	txt := strings.TrimSpace(ev.Text)
	if !reEmail.MatchString(txt) {
		_ = ui.SendText(d.Bot, s.ChatID, "Похоже, это не похоже на email 🙂 Введите корректный e-mail:")
		return EventsRegAskEmail, nil
	}
	profSet(ctx, d, s, keyProfEmail, txt)
	_ = ui.SendText(d.Bot, s.ChatID, "Телефон (только цифры, можно с +):")
	return EventsRegAskPhone, nil
}

func regAskPhone(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	txt := strings.ReplaceAll(strings.TrimSpace(ev.Text), " ", "")
	if !rePhone.MatchString(txt) {
		_ = ui.SendText(d.Bot, s.ChatID, "Телефон обязателен. Введите номер (7–15 цифр, можно с +):")
		return EventsRegAskPhone, nil
	}
	profSet(ctx, d, s, keyProfPhone, txt)

	// Предлагаем Telegram со значением по умолчанию
	handle := ev.FromUserName
	if strings.TrimSpace(handle) == "" {
		handle = "@" + strconv.FormatInt(ev.FromUserID, 10)
	}
	profSet(ctx, d, s, keyProfTelegram, handle)
	_ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("Укажите Telegram (или оставьте как есть):\n%s", handle))
	return EventsRegAskTelegram, nil
}

func regAskTelegram(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	txt := strings.TrimSpace(ev.Text)
	if txt != "" {
		profSet(ctx, d, s, keyProfTelegram, txt)
	}
	_ = ui.SendText(d.Bot, s.ChatID, "Сколько гостей придёт? (число, по умолчанию 1)")
	return EventsRegAskGuests, nil
}

func currentEventID(ctx context.Context, d botengine.Deps, s *types.Session) string {
	if v, ok := stGet(ctx, d, s.ChatID, keyRegEventID); ok {
		return v
	}
	return ""
}

func sendConfirmUI(ctx context.Context, d botengine.Deps, s *types.Session) {
	name, _ := profGet(s, keyProfName)
	email, _ := profGet(s, keyProfEmail)
	phone, _ := profGet(s, keyProfPhone)
	tg, _ := profGet(s, keyProfTelegram)
	dateStr := humanEventDate(ctx, d, s)

	gstr, _ := stGet(ctx, d, s.ChatID, keyRegGuests)
	if gstr == "" {
		gstr = "1"
	}
	guests, _ := strconv.Atoi(gstr)
	comment, _ := stGet(ctx, d, s.ChatID, keyRegComment)
	if strings.TrimSpace(comment) == "" {
		comment = "—"
	}

	summary := fmt.Sprintf(
		"Проверьте данные:\n\nИмя: <b>%s</b>\nEmail: <b>%s</b>\nТелефон: <b>%s</b>\nTelegram: <b>%s</b>\nДата: <b>%s</b>\nГостей: <b>%d</b>\nКомментарий: %s\n",
		htmlEscape(name),
		htmlEscape(email),
		htmlEscape(phone),
		htmlEscape(tg),
		htmlEscape(dateStr),
		guests,
		htmlEscape(comment),
	)

	kb := ui.Inline(
		ui.Row(ui.Cb("✅ Подтвердить", "events:reg:confirm")),
		ui.Row(ui.Cb("➖ Гостей", "events:reg:g:-"), ui.Cb("➕ Гостей", "events:reg:g:+")),
		ui.Row(ui.Cb("✏️ Комментарий", "events:reg:edit:comment")),
		ui.Row(ui.Cb("❌ Отменить регистрацию", "events:rc:ask")),
	)

	if err := ui.SendHTML(d.Bot, s.ChatID, summary, kb); err != nil {
		_ = ui.SendText(d.Bot, s.ChatID, "Не удалось отправить подтверждение, попробуйте ещё раз.")
	}
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

	evID := currentEventID(ctx, d, s)
	if ok, msg := checkCapacityOK(ctx, d, evID, n); !ok {
		_ = ui.SendText(d.Bot, s.ChatID, msg)
		return EventsRegAskGuests, nil
	}

	_ = stSet(ctx, d, s.ChatID, keyRegGuests, strconv.Itoa(n))
	_ = ui.SendText(d.Bot, s.ChatID, "Комментарий (необязательно). Если нечего добавить — отправьте «-».")
	return EventsRegAskComment, nil
}

func regAskComment(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	txt := strings.TrimSpace(ev.Text)
	if txt == "-" {
		txt = ""
	}
	_ = stSet(ctx, d, s.ChatID, keyRegComment, txt)

	sendConfirmUI(ctx, d, s)
	return EventsRegConfirm, nil
}


func incGuests(ctx context.Context, d botengine.Deps, s *types.Session, delta int) {
	gstr, _ := stGet(ctx, d, s.ChatID, keyRegGuests)
	if gstr == "" {
		gstr = "1"
	}
	cur, _ := strconv.Atoi(gstr)
	next := cur + delta
	if next < 1 {
		next = 1
	}
	if next > 20 {
		next = 20
	}

	evID := currentEventID(ctx, d, s)
	if ok, msg := checkCapacityOK(ctx, d, evID, next); !ok {
		_ = ui.SendText(d.Bot, s.ChatID, msg)
		sendConfirmUI(ctx, d, s)
		return
	}

	_ = stSet(ctx, d, s.ChatID, keyRegGuests, strconv.Itoa(next))
	sendConfirmUI(ctx, d, s)
}

func regConfirm(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ackCB(d, ev)
	switch ev.CallbackData {
	case "events:reg:confirm":
		s.Step = EventsRegSubmit
		return botengine.InternalContinue, nil
	case "events:reg:g:+":
		incGuests(ctx, d, s, +1)
		return EventsRegConfirm, nil
	case "events:reg:g:-":
		incGuests(ctx, d, s, -1)
		return EventsRegConfirm, nil
	case "events:reg:edit:comment":
		_ = ui.SendText(d.Bot, s.ChatID, "Ок, пришлите новый комментарий (или «-», чтобы очистить):")
		return EventsRegAskComment, nil
	}
	return EventsRegConfirm, nil
}

func regSubmit(ctx context.Context, _ botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	name, _ := profGet(s, keyProfName)
	email, _ := profGet(s, keyProfEmail)
	phone, _ := profGet(s, keyProfPhone)
	tg, _ := profGet(s, keyProfTelegram)

	if strings.TrimSpace(name) == "" || !reEmail.MatchString(email) || !rePhone.MatchString(phone) {
		_ = ui.SendText(d.Bot, s.ChatID, "Кажется, не все обязательные поля заполнены. Давай начнём заново: /events")
		return EventsDone, nil
	}
	guestsStr, _ := stGet(ctx, d, s.ChatID, keyRegGuests)
	if strings.TrimSpace(guestsStr) == "" {
		guestsStr = "1"
	}
	evID := currentEventID(ctx, d, s)
	comment, _ := stGet(ctx, d, s.ChatID, keyRegComment)
	dateHuman := humanEventDate(ctx, d, s)

	need, _ := strconv.Atoi(guestsStr)
	if ok, msg := checkCapacityOK(ctx, d, evID, need); !ok {
		_ = ui.SendText(d.Bot, s.ChatID, msg)
		return EventsDone, nil
	}

	body := regPayload{
		Name:            name,
		Email:           email,
		Phone:           phone,
		NumberOfPersons: guestsStr,
		Telegram:        tg,
		Date:            dateHuman,
		EventID:         evID,
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

	var created struct {
		Id string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&created); err == nil && created.Id != "" {
		_ = stSet(ctx, d, s.ChatID, keyRegEntryID, created.Id)
	}

	text := "Спасибо за регистрацию! 🎉\n\n" +
		"Пожалуйста, не закрывайте и не удаляйте бота — иначе мы не сможем прислать напоминание и важные детали мероприятия.\n" +
		"Если что-то изменится — просто напишите нам сюда в чат.\n\n" +
		"До встречи!"
	kb := ui.Inline(ui.Row(ui.Cb("❌ Отменить регистрацию", "events:rc:ask")))
	_ = ui.SendHTML(d.Bot, s.ChatID, htmlEscape(text), kb)


	scheduleReminders(ctx, d, s)

	s.Flow, s.Step = "", ""
	return EventsDone, nil
}

// ====================== Отмена регистрации (PUT willCome=false) ======================

type updateEntriePayload struct {
	Id              string `json:"id,omitempty"`
	CreationDate    string `json:"creationDate,omitempty"`
	Name            string `json:"name,omitempty"`
	Email           string `json:"email,omitempty"`
	Phone           string `json:"phone,omitempty"`
	NumberOfPersons string `json:"numberOfPersons,omitempty"`
	Instagram       string `json:"instagram,omitempty"`
	Telegram        string `json:"telegram,omitempty"`
	Date            string `json:"date,omitempty"`
	EventId         string `json:"eventId,omitempty"`
	Comment         string `json:"comment,omitempty"`
	WillCome        bool   `json:"willCome"`
}

func updateWillCome(ctx context.Context, d botengine.Deps, s *types.Session, eventID string, will bool) error {
	// профайл
	name, _  := profGet(s, keyProfName)
	email, _ := profGet(s, keyProfEmail)
	phone, _ := profGet(s, keyProfPhone)
	tg, _    := profGet(s, keyProfTelegram)

	// гости/коммент
	guests, _ := stGet(ctx, d, s.ChatID, keyRegGuests)
	if strings.TrimSpace(guests) == "" {
		guests = "1"
	}
	comment, _ := stGet(ctx, d, s.ChatID, keyRegComment)

	// дата: сначала из KV (RFC3339), если нет — из списка событий
	var dateStr string
	if raw, ok := stGet(ctx, d, s.ChatID, keyRegEventDate); ok && raw != "" {
		if t, err := time.Parse(time.RFC3339, raw); err == nil {
			dateStr = formatRuHuman(t.In(userLoc(s)))
		}
	}
	if dateStr == "" {
		if _, t := loadEventByID(ctx, d, eventID); !t.IsZero() {
			dateStr = formatRuHuman(t.In(userLoc(s)))
		}
	}

	// возможный id записи
	var entryID string
	if id, ok := stGet(ctx, d, s.ChatID, keyRegEntryID); ok {
		entryID = id
	}

	p := updateEntriePayload{
		Id:              entryID,
		Name:            name,
		Email:           email,
		Phone:           phone,
		NumberOfPersons: guests,
		Telegram:        tg,
		Date:            dateStr,
		EventId:         eventID,
		Comment:         comment,
		WillCome:        will,
	}

	b, _ := json.Marshal(p)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPut, updateEntryEndpoint, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("update returned %d", resp.StatusCode)
	}
	return nil
}

func regCancelAsk(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ackCB(d, ev)
	msg := "Ой… Нам очень жаль 😿 Мы готовимся к каждому гостю и бережём места.\n" +
		"Точно отменяем? (можно просто прийти на другое событие — мы будем рады!)"
	kb := ui.Inline(
		ui.Row(ui.Cb("Да, отменить", "events:rc:yes"), ui.Cb("Оставить регистрацию", "events:rc:no")),
	)
	_ = ui.SendHTML(d.Bot, s.ChatID, htmlEscape(msg), kb)
	return EventsRegCancelDo, nil
}

func regCancelDo(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ackCB(d, ev)
	switch ev.CallbackData {
	case "events:rc:yes":
		// пометим статус и отменим будущие таймеры
		if evID := currentEventID(ctx, d, s); evID != "" {
			_ = stSet(ctx, d, s.ChatID, remStatusKey(evID), "canceled")
			cancelTimers(s.ChatID, evID)

			if err := updateWillCome(ctx, d, s, evID, true); err != nil {
				_ = ui.SendText(d.Bot, s.ChatID, "Не удалось отменить автоматически. Мы отметили у себя, но на всякий случай напишите нам: @lettersandnumbers_am 🙏")
			} else {
				_ = ui.SendText(d.Bot, s.ChatID, "Окей, мы отметили отмену. Если передумаете — снова жмякните /events ❤️")
			}
		}
	case "events:rc:no":
		_ = ui.SendText(d.Bot, s.ChatID, "Ура! Мы вас ждём 🥳")
	}
	s.Flow, s.Step = "", ""
	return EventsDone, nil
}


// ====================== done ======================

func done(_ context.Context, _ botengine.Event, _ botengine.Deps, _ *types.Session) (types.Step, error) {
	return EventsDone, nil
}

// ====================== Хелперы ======================

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
	// фикс кавычек
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", `'`, "&#39;")
	return r.Replace(s)
}

func eventID(e types.Event) string { return e.ID }
func hasShowFormField(_ types.Event) bool { return true }

func fetchEventsFallback(ctx context.Context, baseURL string) ([]types.Event, error) {
	sep := "?"
	if strings.Contains(baseURL, "?") { sep = "&" }
	u := fmt.Sprintf("%s%sts=%d", baseURL, sep, time.Now().UnixNano())

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	req.Header.Set("Cache-Control", "no-cache, no-store, max-age=0")
	req.Header.Set("Pragma", "no-cache")

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

func humanEventDate(ctx context.Context, d botengine.Deps, s *types.Session) string {
	raw, ok := stGet(ctx, d, s.ChatID, keyRegEventDate)
	if ok && raw != "" {
		if t, err := time.Parse(time.RFC3339, raw); err == nil {
			return formatRuHuman(t.In(userLoc(s)))
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

// --- capacity ---

func checkCapacityOK(ctx context.Context, d botengine.Deps, evID string, need int) (bool, string) {
	if evID == "" || need < 1 {
		return true, ""
	}

	// найдём событие
	var e *types.Event
	if d.Svcs.Events != nil {
		if list, _ := d.Svcs.Events.ListUpcoming(ctx); len(list) > 0 {
			for i := range list {
				if eventID(list[i]) == evID {
					e = &list[i]
					break
				}
			}
		}
	} else {
		if list, err := fetchEventsFallback(ctx, eventsURLFallback); err == nil {
			for i := range list {
				if eventID(list[i]) == evID {
					e = &list[i]
					break
				}
			}
		}
	}
	if e == nil || e.Capacity <= 0 { // без capacity — не ограничиваем
		return true, ""
	}

	counts, err := fetchEntriesCounts(ctx, entriesUniqueURL)
	if err != nil {
		return true, ""
	}
	used := counts[evID]
	left := e.Capacity - used
	if left < need {
		if left <= 0 {
			return false, "К сожалению, места уже закончились 😿 Попробуйте другое событие — /events"
		}
		return false, fmt.Sprintf("Осталось только %d мест(а). Введите число не больше %d:", left, left)
	}
	return true, ""
}

func fetchEntriesCounts(ctx context.Context, baseURL string) (map[string]int, error) {
	sep := "?"
	if strings.Contains(baseURL, "?") { sep = "&" }
	u := fmt.Sprintf("%s%sts=%d", baseURL, sep, time.Now().UnixNano())

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	req.Header.Set("Cache-Control", "no-cache, no-store, max-age=0")
	req.Header.Set("Pragma", "no-cache")

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

// ===== Напоминания =====

// В проде поставь false — тогда будут 2 реальных пуша: D-1@12:00 и H-4
const remindersTestMode = true

// Тестовые интервалы от момента регистрации
var (
	testReminder1 = 1 * time.Minute
	testReminder2 = 2 * time.Minute
)

// Ключ статуса уведомлений на конкретное событие: "confirmed"/"canceled"
func remStatusKey(eventID string) string { return "events:rem:status:" + eventID }

// Хранилище таймеров, чтобы уметь отменять второй пуш после ответа
var reminderJobs = struct {
	mu     sync.Mutex
	timers map[string][]*time.Timer // key: "<chatID>:<eventID>"
}{timers: make(map[string][]*time.Timer)}

func rkey(chatID int64, eventID string) string { return fmt.Sprintf("%d:%s", chatID, eventID) }

func rememberTimers(chatID int64, eventID string, ts ...*time.Timer) {
	reminderJobs.mu.Lock()
	defer reminderJobs.mu.Unlock()
	reminderJobs.timers[rkey(chatID, eventID)] = ts
}

func cancelTimers(chatID int64, eventID string) {
	reminderJobs.mu.Lock()
	defer reminderJobs.mu.Unlock()
	if arr, ok := reminderJobs.timers[rkey(chatID, eventID)]; ok {
		for _, t := range arr {
			if t != nil {
				t.Stop()
			}
		}
		delete(reminderJobs.timers, rkey(chatID, eventID))
	}
}

// Ставит таймеры: в тесте 1m и 2m; в проде — D-1@12:00 и H-4
func scheduleReminders(ctx context.Context, d botengine.Deps, s *types.Session) {
	evID, ok := stGet(ctx, d, s.ChatID, keyRegEventID)
	if !ok || evID == "" {
		return
	}
	raw, ok := stGet(ctx, d, s.ChatID, keyRegEventDate)
	if !ok || raw == "" {
		return
	}
	tUTC, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return
	}

	// Если уже подтвердил/отменил — не ставим
	if st, ok := stGet(ctx, d, s.ChatID, remStatusKey(evID)); ok && (st == "confirmed" || st == "canceled") {
		return
	}

	loc := userLoc(s)
	tLocal := tUTC.In(loc)

	var when1, when2 time.Time
	if remindersTestMode {
		when1 = time.Now().Add(testReminder1)
		when2 = time.Now().Add(testReminder2)
	} else {
		// 1) За сутки, в дневное время — возьмём 12:00 локали
		dayBeforeNoon := time.Date(tLocal.Year(), tLocal.Month(), tLocal.Day(), 12, 0, 0, 0, loc).AddDate(0, 0, -1)
		// 2) За 4 часа до начала
		before4h := tLocal.Add(-4 * time.Hour)

		now := time.Now().In(loc)
		if dayBeforeNoon.After(now) {
			when1 = dayBeforeNoon
		}
		if before4h.After(now) {
			when2 = before4h
		}
	}

	// Если обе даты уже в прошлом — нечего планировать
	if when1.IsZero() && when2.IsZero() {
		return
	}

	var timers []*time.Timer
	if !when1.IsZero() {
		dur := time.Until(when1)
		if dur < 0 {
			dur = 0
		}
		t := time.AfterFunc(dur, func() {
			sendReminder(d, s.ChatID, s.Lang, evID, "D-1")
		})
		timers = append(timers, t)
	}
	if !when2.IsZero() {
		dur := time.Until(when2)
		if dur < 0 {
			dur = 0
		}
		t := time.AfterFunc(dur, func() {
			sendReminder(d, s.ChatID, s.Lang, evID, "H-4")
		})
		timers = append(timers, t)
	}

	rememberTimers(s.ChatID, evID, timers...)
}

func sendReminder(d botengine.Deps, chatID int64, lang, eventID, tag string) {
	// Если уже подтвердил/отменил — не шлём
	if st, ok := stGet(context.Background(), d, chatID, remStatusKey(eventID)); ok && (st == "confirmed" || st == "canceled") {
		return
	}

	// Попытаемся найти событие (для названия/времени)
	name := "мероприятие"
	dateStr := "скоро"
	if e, t := loadEventByID(context.Background(), d, eventID); e != nil {
		if strings.TrimSpace(e.Name) != "" {
			name = strings.TrimSpace(e.Name)
		}
		dateStr = formatRuHuman(t.In(userLoc(&types.Session{Lang: lang})))
	}

	// Текст
	var prefix string
	switch tag {
	case "D-1":
		if remindersTestMode {
			prefix = "Тест-напоминание"
		} else {
			prefix = "Напоминание: завтра"
		}
	case "H-4":
		if remindersTestMode {
			prefix = "Тест-напоминание №2"
		} else {
			prefix = "Напоминание: через ~4 часа"
		}
	default:
		prefix = "Напоминание"
	}

	msg := fmt.Sprintf("%s о событии «%s» — <b>%s</b>.\n\nПодтверди участие или, если планы поменялись, отмени пожалуйста 🙏",
		prefix, htmlEscape(name), htmlEscape(dateStr))

	kb := ui.Inline(
		ui.Row(
			ui.Cb("✅ Я приду", "events:rem:c:"+eventID),
			ui.Cb("❌ Отмениться", "events:rem:x:"+eventID),
		),
	)
	_ = ui.SendHTML(d.Bot, chatID, msg, kb)
}

func loadEventByID(ctx context.Context, d botengine.Deps, id string) (*types.Event, time.Time) {
	var events []types.Event
	var err error
	if d.Svcs.Events != nil {
		events, err = d.Svcs.Events.ListUpcoming(ctx)
	} else {
		events, err = fetchEventsFallback(ctx, eventsURLFallback)
	}
	if err != nil {
		return nil, time.Time{}
	}
	for i := range events {
		if eventID(events[i]) == id {
			t, _ := parseAnyEventDate(events[i].Date)
			return &events[i], t
		}
	}
	return nil, time.Time{}
}

// Обработчик кликов из напоминаний
func remindHandle(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ackCB(d, ev)

	var action, evID string
	switch {
	case strings.HasPrefix(ev.CallbackData, "events:rem:c:"):
		action = "confirm"
		evID = strings.TrimPrefix(ev.CallbackData, "events:rem:c:")
	case strings.HasPrefix(ev.CallbackData, "events:rem:x:"):
		action = "cancel"
		evID = strings.TrimPrefix(ev.CallbackData, "events:rem:x:")
	default:
		return EventsDone, nil
	}
	if evID == "" {
		return EventsDone, nil
	}

	switch action {
	case "confirm":
		_ = stSet(ctx, d, s.ChatID, remStatusKey(evID), "confirmed")
		cancelTimers(s.ChatID, evID)
		if err := updateWillCome(ctx, d, s, evID, true); err != nil {
			_ = ui.SendText(d.Bot, s.ChatID, "Подтвердили у нас ✅ Но сервер сейчас недоступен, мы попробуем ещё раз позже.")
		} else {
			_ = ui.SendText(d.Bot, s.ChatID, "Ура! Отметили, что вы придёте 🥳 До встречи!")
		}
	case "cancel":
		_ = stSet(ctx, d, s.ChatID, remStatusKey(evID), "canceled")
		cancelTimers(s.ChatID, evID)
		if err := updateWillCome(ctx, d, s, evID, false); err != nil {
			_ = ui.SendText(d.Bot, s.ChatID, "Мы отменили локально ❌ Но сервер сейчас недоступен, на всякий случай напишите нам: @lettersandnumbers_am")
		} else {
			_ = ui.SendText(d.Bot, s.ChatID, "Окей, отменили запись. Если планы изменятся — загляните в /events ❤️")
		}
	}

	s.Flow, s.Step = "", ""
	return EventsDone, nil
}
