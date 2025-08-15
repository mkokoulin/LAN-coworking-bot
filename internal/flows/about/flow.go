package flow

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowAbout  types.Flow = "about"
	AboutSend  types.Step = "about:send"
	AboutDone  types.Step = "about:done"
)

func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowAbout, map[types.Step]botengine.StepHandler{
		AboutSend: send,
		AboutDone: done,
	})

	reg.RegisterCommand("about", botengine.FlowEntry{Flow: FlowAbout, Step: AboutSend})
}
