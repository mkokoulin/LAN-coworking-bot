package flows

import (
	"context"
	"regexp"
	"sync"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

// --- внешние API ---
const (
	entriesUniqueURL     = "https://shark-app-wrcei.ondigitalocean.app/api/entries/unique"
	registrationEndpoint = "https://shark-app-wrcei.ondigitalocean.app/api/entries"
)

// --- Flow / Steps ---

const (
	FlowEvents types.Flow = "events"

	EventsIntro types.Step = "events:intro"
	EventsList  types.Step = "events:list"
	EventsDone  types.Step = "events:done"

	// Подписка / отписка
	EventsSub   types.Step = "events:subscribe"
	EventsUnsub types.Step = "events:unsubscribe"

	// Мастер выбора расписания
	EventsSubPickDay    types.Step = "events:sub_pick_day"
	EventsSubPickTime   types.Step = "events:sub_pick_time"
	EventsSubAwaitInput types.Step = "events:sub_await_time_text"
	EventsSubConfirm    types.Step = "events:sub_confirm"

	// Редактирование расписания
	EventsEditStart types.Step = "events:edit_start"
)

// Регистрация
const (
	EventsRegStart       types.Step = "events:reg_start"
	EventsRegAskName     types.Step = "events:reg_name"
	EventsRegAskEmail    types.Step = "events:reg_email"
	EventsRegAskPhone    types.Step = "events:reg_phone"
	EventsRegAskGuests   types.Step = "events:reg_guests"
	EventsRegAskTelegram types.Step = "events:reg_telegram"
	EventsRegAskComment  types.Step = "events:reg_comment"
	EventsRegConfirm     types.Step = "events:reg_confirm"
	EventsRegSubmit      types.Step = "events:reg_submit"
	EventsRegCancelAsk   types.Step = "events:reg_cancel_ask"
	EventsRegCancelDo    types.Step = "events:reg_cancel_do"
)

// Ключи KV (in-memory)
const (
	keyRegCapacity   = "events:reg:capacity"
	keyRegEventID    = "events:reg:event_id"
	keyRegEventDate  = "events:reg:event_date" // RFC3339
	keyRegName       = "events:reg:name"
	keyRegEmail      = "events:reg:email"
	keyRegPhone      = "events:reg:phone"
	keyRegGuests     = "events:reg:guests"
	keyRegTelegram   = "events:reg:telegram"
	keyRegComment    = "events:reg:comment"
	keyRegReminderAt = "events:reg:reminder_at" // RFC3339 — напоминание в день события
)

// --- простая in-memory KV на процесс (т.к. d.Store нет) ---
var memKV = struct {
	mu   sync.RWMutex
	data map[int64]map[string]string
}{data: make(map[int64]map[string]string)}

func stSet(_ context.Context, _ botengine.Deps, chatID int64, key, val string) error {
	memKV.mu.Lock()
	defer memKV.mu.Unlock()
	if _, ok := memKV.data[chatID]; !ok {
		memKV.data[chatID] = make(map[string]string)
	}
	memKV.data[chatID][key] = val
	return nil
}
func stGet(_ context.Context, _ botengine.Deps, chatID int64, key string) (string, bool) {
	memKV.mu.RLock()
	defer memKV.mu.RUnlock()
	if m, ok := memKV.data[chatID]; ok {
		v, ok2 := m[key]
		return v, ok2
	}
	return "", false
}
func stDel(_ context.Context, _ botengine.Deps, chatID int64, key string) error {
	memKV.mu.Lock()
	defer memKV.mu.Unlock()
	if m, ok := memKV.data[chatID]; ok {
		delete(m, key)
	}
	return nil
}

// --- модели/валидация для POST ---
type regPayload struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	NumberOfPersons string `json:"numberOfPersons"`
	Telegram        string `json:"telegram"`
	Date            string `json:"date"` // «ср 20 августа. 19:30»
	EventID         string `json:"eventId"`
	Comment         string `json:"comment"`
}

var (
	reEmail = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	rePhone = regexp.MustCompile(`^\+?\d{7,15}$`)
	reHHMM  = regexp.MustCompile(`^(?:[01]?\d|2[0-3]):[0-5]\d$`)
)

// --- регистрация в botengine.Registry ---

func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowEvents, map[types.Step]botengine.StepHandler{
		// список
		EventsIntro: intro,
		EventsList:  list,
		EventsDone:  done,

		// подписка
		EventsSub:           subscribe,
		EventsUnsub:         unsubscribe,
		EventsSubPickDay:    subPickDay,
		EventsSubPickTime:   subPickTime,
		EventsSubAwaitInput: subAwaitTimeText,
		EventsSubConfirm:    subConfirm,
		EventsEditStart:     editSchedule,

		// регистрация
		EventsRegStart:       regStart,
		EventsRegAskName:     regAskName,
		EventsRegAskEmail:    regAskEmail,
		EventsRegAskPhone:    regAskPhone,
		EventsRegAskGuests:   regAskGuests,
		EventsRegAskTelegram: regAskTelegram,
		EventsRegAskComment:  regAskComment,
		EventsRegConfirm:     regConfirm,
		EventsRegSubmit:      regSubmit,
		EventsRegCancelAsk:   regCancelAsk,
		EventsRegCancelDo:    regCancelDo,
	})

	 // Команды
    reg.RegisterCommand("events",              botengine.FlowEntry{Flow: FlowEvents, Step: EventsList})
    reg.RegisterCommand("events_time",         botengine.FlowEntry{Flow: FlowEvents, Step: EventsEditStart})
    reg.RegisterCommand("unsubscribe_events",  botengine.FlowEntry{Flow: FlowEvents, Step: EventsUnsub})

    // Подписка
    reg.RegisterCallbackPrefix("events:subscribe",   botengine.FlowEntry{Flow: FlowEvents, Step: EventsSub})
    reg.RegisterCallbackPrefix("events:edit",        botengine.FlowEntry{Flow: FlowEvents, Step: EventsEditStart})
    reg.RegisterCallbackPrefix("events:unsubscribe", botengine.FlowEntry{Flow: FlowEvents, Step: EventsUnsub})
    reg.RegisterCallbackPrefix("events:sub:day:",    botengine.FlowEntry{Flow: FlowEvents, Step: EventsSubPickDay})
    reg.RegisterCallbackPrefix("events:sub:time:",   botengine.FlowEntry{Flow: FlowEvents, Step: EventsSubPickTime})

    // ✅ Регистрация — СНАЧАЛА точные префиксы
    reg.RegisterCallbackPrefix("events:reg:confirm", botengine.FlowEntry{Flow: FlowEvents, Step: EventsRegConfirm})
    reg.RegisterCallbackPrefix("events:reg:edit:",   botengine.FlowEntry{Flow: FlowEvents, Step: EventsRegConfirm})
    reg.RegisterCallbackPrefix("events:reg_cancel",  botengine.FlowEntry{Flow: FlowEvents, Step: EventsRegCancelAsk})

    // 👇 А общий — В САМОМ КОНЦЕ, чтобы не перехватывал confirm/edit
    reg.RegisterCallbackPrefix("events:reg:",        botengine.FlowEntry{Flow: FlowEvents, Step: EventsRegStart})
}
