package start

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowStart types.Flow = "start"
	StepShow  types.Step = "start:show"
	StepDone  types.Step = "start:done"
)

func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowStart, map[types.Step]botengine.StepHandler{
		StepShow: show,
		StepDone: done,
	})

	// вход в стартовый сценарий по /start
	reg.RegisterCommand("start", botengine.FlowEntry{Flow: FlowStart, Step: StepShow})
}
