package flow

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowUnknown types.Flow = "unknown"
	UnknownSend      types.Step = "unknown:send"
	UnknownDone       types.Step = "unknown:done"
)


func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowUnknown, map[types.Step]botengine.StepHandler{
		UnknownSend:      send,
		UnknownDone:       done,
	})

	reg.RegisterCommand("unknown", botengine.FlowEntry{Flow: FlowUnknown, Step: UnknownSend})
}
