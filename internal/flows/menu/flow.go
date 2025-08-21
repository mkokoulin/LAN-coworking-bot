package flows

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowMenu types.Flow = "menu"
	MenuSend types.Step = "menu:send"
	MenuDone types.Step = "menu:done"
)

func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowMenu, map[types.Step]botengine.StepHandler{
		MenuSend: send,
		MenuDone: done,
	})

	// вход по /menu
	reg.RegisterCommand("menu", botengine.FlowEntry{Flow: FlowMenu, Step: MenuSend})
}
