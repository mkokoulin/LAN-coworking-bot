package flow

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowBooking types.Flow = "booking"
	BookInfo    types.Step = "booking:info"
	BookDone    types.Step = "booking:done"
)

func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowBooking, map[types.Step]botengine.StepHandler{
		BookInfo: info,
		BookDone: done,
	})

	// вход по команде
	reg.RegisterCommand("booking", botengine.FlowEntry{Flow: FlowBooking, Step: BookInfo})
	// можно алиас:
	// reg.RegisterCommand("book", botengine.FlowEntry{Flow: FlowBooking, Step: BookInfo})
}
