package flows

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowEvents types.Flow = "events"

	EventsIntro types.Step = "events:intro"
	EventsList  types.Step = "events:list"
	EventsDone  types.Step = "events:done"

	// Подписка / отписка
	EventsSub   types.Step = "events:subscribe"
	EventsUnsub types.Step = "events:unsubscribe"

	// Новый мастер выбора расписания
	EventsSubPickDay    types.Step = "events:sub_pick_day"
	EventsSubPickTime   types.Step = "events:sub_pick_time"
	EventsSubAwaitInput types.Step = "events:sub_await_time_text"
	EventsSubConfirm    types.Step = "events:sub_confirm"

	// Редактирование расписания
	EventsEditStart types.Step = "events:edit_start"
)

func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowEvents, map[types.Step]botengine.StepHandler{
		EventsIntro:         intro,
		EventsList:          list,
		EventsDone:       done, // можно не регать, если сразу выходим из флоу
		EventsSub:           subscribe,
		EventsUnsub:         unsubscribe,
		EventsSubPickDay:    subPickDay,
		EventsSubPickTime:   subPickTime,
		EventsSubAwaitInput: subAwaitTimeText,
		EventsSubConfirm:    subConfirm,
		EventsEditStart:     editSchedule,
	})

	// Команды
	reg.RegisterCommand("events",        botengine.FlowEntry{Flow: FlowEvents, Step: EventsList})
	reg.RegisterCommand("events_time",   botengine.FlowEntry{Flow: FlowEvents, Step: EventsEditStart})
	reg.RegisterCommand("unsubscribe_events", botengine.FlowEntry{Flow: FlowEvents, Step: EventsUnsub})

	// Точные префиксы коллбэков
	reg.RegisterCallbackPrefix("events:subscribe",   botengine.FlowEntry{Flow: FlowEvents, Step: EventsSub})
	reg.RegisterCallbackPrefix("events:edit",        botengine.FlowEntry{Flow: FlowEvents, Step: EventsEditStart})
	reg.RegisterCallbackPrefix("events:unsubscribe", botengine.FlowEntry{Flow: FlowEvents, Step: EventsUnsub})

	// ✅ КЛЮЧЕВОЕ: направляем выбор дня/времени сразу в соответствующие шаги
	reg.RegisterCallbackPrefix("events:sub:day:",    botengine.FlowEntry{Flow: FlowEvents, Step: EventsSubPickDay})
	reg.RegisterCallbackPrefix("events:sub:time:",   botengine.FlowEntry{Flow: FlowEvents, Step: EventsSubPickTime})
}

