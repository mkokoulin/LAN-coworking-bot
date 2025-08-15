
package flows

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

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
		// если ShowForm есть и false — пропускаем
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
        if name == "" { name = "Untitled" }
        if len([]rune(name)) > 120 { name = string([]rune(name)[:117]) + "…" }

        desc := strings.TrimSpace(stripHTML(e.Description))
        if len([]rune(desc)) > 200 { desc = string([]rune(desc)[:197]) + "…" }

        url := fmt.Sprintf("https://lettersandnumbers.am/events/%s", eventID(e))

        sb.WriteString(fmt.Sprintf("• <b>%s</b> <i>(%s)</i> — <b>%s</b>\n", 
            htmlEscape(date), htmlEscape(wd), htmlEscape(name)))
        if desc != "" {
            sb.WriteString(htmlEscape(desc))
            sb.WriteString("\n")
        }
        sb.WriteString(fmt.Sprintf("<a href=\"%s\">Registration →</a>\n\n", htmlEscape(url)))
    }

    // одна отправка: HTML + inline-кнопка
    kb := ui.Inline(
        ui.Row(
            ui.Cb("📬 Подписаться на еженедельные анонсы", "events:subscribe"),
        ),
    )
    if err := ui.SendHTML(d.Bot, s.ChatID, sb.String(), kb); err != nil {
        _ = ui.SendText(d.Bot, s.ChatID, fmt.Sprintf("[events] send error: %v", err))
    }

    s.Flow, s.Step = "", ""
    return EventsDone, nil
}

// --- helpers ---

func parseAnyEventDate(s string) (time.Time, error) {
	// Примеры: "2025-08-12", "2025-08-12T18:00:00Z", "12.08.2025"
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

func eventID(e types.Event) string {
	// Если в структуре только ID — используем его.
	// Если у тебя ещё есть поле Id (с маленькой буквы), добавь сюда fallback.
	return e.ID
}

func hasShowFormField(e types.Event) bool {
	// Если ShowForm гарантированно есть — оставь true.
	// Если нет — убери проверку выше (или детектируй наличия поля).
	return true
}

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


func subscribe(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
    if ev.Kind != botengine.EventCallback || ev.CallbackData != "events:subscribe" {
        return EventsSub, nil
    }
    s.IsSubscribed = true  // 👈 главное действие
    _ = ui.SendText(d.Bot, s.ChatID, "Готово! Будем присылать список мероприятий раз в неделю. Чтобы отписаться — /unsubscribe_events")
    s.Flow, s.Step = "", ""
    return EventsDone, nil
}

func unsubscribe(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
    s.IsSubscribed = false // 👈 снять флаг
    _ = ui.SendText(d.Bot, s.ChatID, "Вы отписаны от еженедельных анонсов.")
    s.Flow, s.Step = "", ""
    return EventsDone, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return EventsDone, nil
}
