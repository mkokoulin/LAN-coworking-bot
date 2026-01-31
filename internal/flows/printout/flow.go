package flows

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowPrintout types.Flow = "printout"

	PrintoutShow types.Step = "printout:show"
	PrintoutDone types.Step = "printout:done"
)

func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowPrintout, map[types.Step]botengine.StepHandler{
		PrintoutShow: showInfo,
		PrintoutDone: done,
	})

	reg.RegisterCommand("printout", botengine.FlowEntry{Flow: FlowPrintout, Step: PrintoutShow})
}
