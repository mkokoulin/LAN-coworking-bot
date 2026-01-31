package flows

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowCoworking  types.Flow = "coworking"
	CoworkingHome  types.Step = "coworking:home"
	dram                      = "֏"
)

// -------- Регистрация --------

func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowCoworking, map[types.Step]botengine.StepHandler{
		CoworkingHome: coworkingHome,
	})

	reg.RegisterCommand("coworking", botengine.FlowEntry{Flow: FlowCoworking, Step: CoworkingHome})
}
