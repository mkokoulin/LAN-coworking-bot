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
	EventsSub   types.Step = "events:subscribe"
	EventsUnsub   types.Step = "events:unsubscribe"
)

func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowEvents, map[types.Step]botengine.StepHandler{
		EventsIntro: intro,
		EventsList:  list,
		EventsDone:  done,
		EventsSub:   subscribe,
		EventsUnsub: unsubscribe,
	})

	reg.RegisterCommand("events", botengine.FlowEntry{Flow: FlowEvents, Step: EventsList})
	reg.RegisterCallbackPrefix("events:", botengine.FlowEntry{Flow: FlowEvents, Step: EventsSub})
	reg.RegisterCommand("unsubscribe_events", botengine.FlowEntry{Flow: FlowEvents, Step: EventsUnsub})
}
