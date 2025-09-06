package flows

import (
	"regexp"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

// --- внешние API ---
const (
	entriesUniqueURL     = "https://shark-app-wrcei.ondigitalocean.app/api/entries/unique"
	eventsURLFallback    = "https://shark-app-wrcei.ondigitalocean.app/api/events"
	registrationEndpoint = "https://shark-app-wrcei.ondigitalocean.app/api/entries"
	updateEntryEndpoint  = "https://shark-app-wrcei.ondigitalocean.app/api/entries/update"
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

	// Регистрация
	EventsRegStart       types.Step = "events:reg_start"
	EventsRegAskName     types.Step = "events:reg_name"
	EventsRegAskEmail    types.Step = "events:reg_email"
	EventsRegAskPhone    types.Step = "events:reg_phone"
	EventsRegAskTelegram types.Step = "events:reg_telegram"
	EventsRegAskGuests   types.Step = "events:reg_guests"
	EventsRegAskComment  types.Step = "events:reg_comment"
	EventsRegConfirm     types.Step = "events:reg_confirm"
	EventsRegSubmit      types.Step = "events:reg_submit"
	EventsRegCancelAsk   types.Step = "events:reg_cancel_ask"
	EventsRegCancelDo    types.Step = "events:reg_cancel_do"

	EventsRemindHandle types.Step = "events:rem_handle"
)

// --- Ключи KV / Session Data ---

// Профиль пользователя (долгоживущие данные в s.Data)
const (
	keyProfName     = "profile:name"
	keyProfEmail    = "profile:email"
	keyProfPhone    = "profile:phone"
	keyProfTelegram = "profile:telegram"
)

// Временные ключи регистрации (in-memory)
const (
	keyRegCapacity  = "events:reg:capacity"
	keyRegEventID   = "events:reg:event_id"
	keyRegEventDate = "events:reg:event_date" // RFC3339
	keyRegGuests    = "events:reg:guests"
	keyRegComment   = "events:reg:comment"
	keyRegEntryID   = "events:reg:entry_id" // id записи от бэкенда, если вернёт
)

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

// --- Регистрация обработчиков в Registry ---

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
		EventsRegAskTelegram: regAskTelegram,
		EventsRegAskGuests:   regAskGuests,
		EventsRegAskComment:  regAskComment,
		EventsRegConfirm:     regConfirm,
		EventsRegSubmit:      regSubmit,
		EventsRegCancelAsk:   regCancelAsk,
		EventsRegCancelDo:    regCancelDo,

		// напоминания
		EventsRemindHandle: remindHandle,
	})

	// Команды
	reg.RegisterCommand("events",             botengine.FlowEntry{Flow: FlowEvents, Step: EventsList})
	reg.RegisterCommand("events_time",        botengine.FlowEntry{Flow: FlowEvents, Step: EventsEditStart})
	reg.RegisterCommand("unsubscribe_events", botengine.FlowEntry{Flow: FlowEvents, Step: EventsUnsub})

	// Подписка
	reg.RegisterCallbackPrefix("events:subscribe",   botengine.FlowEntry{Flow: FlowEvents, Step: EventsSub})
	reg.RegisterCallbackPrefix("events:edit",        botengine.FlowEntry{Flow: FlowEvents, Step: EventsEditStart})
	reg.RegisterCallbackPrefix("events:unsubscribe", botengine.FlowEntry{Flow: FlowEvents, Step: EventsUnsub}) // ✅ убрал лишний аргумент
	reg.RegisterCallbackPrefix("events:sub:day:",    botengine.FlowEntry{Flow: FlowEvents, Step: EventsSubPickDay})
	reg.RegisterCallbackPrefix("events:sub:time:",   botengine.FlowEntry{Flow: FlowEvents, Step: EventsSubPickTime})

	// старт регистрации по событию
	reg.RegisterCallbackPrefix("events:regstart:", botengine.FlowEntry{Flow: FlowEvents, Step: EventsRegStart})

	// подтверждение / изменение гостей (+/-) / правка комментария
	reg.RegisterCallbackPrefix("events:reg:confirm",      botengine.FlowEntry{Flow: FlowEvents, Step: EventsRegConfirm})
	reg.RegisterCallbackPrefix("events:reg:g:+",          botengine.FlowEntry{Flow: FlowEvents, Step: EventsRegConfirm})
	reg.RegisterCallbackPrefix("events:reg:g:-",          botengine.FlowEntry{Flow: FlowEvents, Step: EventsRegConfirm})
	reg.RegisterCallbackPrefix("events:reg:edit:comment", botengine.FlowEntry{Flow: FlowEvents, Step: EventsRegConfirm})

	// отмена регистрации
	reg.RegisterCallbackPrefix("events:rc:ask", botengine.FlowEntry{Flow: FlowEvents, Step: EventsRegCancelAsk})
	reg.RegisterCallbackPrefix("events:rc:yes", botengine.FlowEntry{Flow: FlowEvents, Step: EventsRegCancelDo})
	reg.RegisterCallbackPrefix("events:rc:no",  botengine.FlowEntry{Flow: FlowEvents, Step: EventsRegCancelDo})

	reg.RegisterCallbackPrefix("events:rem:c:", botengine.FlowEntry{Flow: FlowEvents, Step: EventsRemindHandle}) // confirm
	reg.RegisterCallbackPrefix("events:rem:x:", botengine.FlowEntry{Flow: FlowEvents, Step: EventsRemindHandle}) // cancel
}
