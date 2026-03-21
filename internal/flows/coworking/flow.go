package flows

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowCoworking types.Flow = "coworking"

	CoworkingHome           types.Step = "coworking:home"
	CoworkingNewName        types.Step = "coworking:new_name"
	CoworkingNewPhone       types.Step = "coworking:new_phone"
	CoworkingNewTariff      types.Step = "coworking:new_tariff"
	CoworkingConfirm        types.Step = "coworking:confirm"
	CoworkingPending        types.Step = "coworking:pending"
	CoworkingReturningPhone types.Step = "coworking:returning_phone"
	CoworkingAdminAction    types.Step = "coworking:admin_action"
	CoworkingProfile 		types.Step = "coworking:profile"
	CoworkingTariff  			types.Step = "coworking:tariff"
)

func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowCoworking, map[types.Step]botengine.StepHandler{
		CoworkingHome:           coworkingHome,
		CoworkingNewName:        coworkingNewName,
		CoworkingNewPhone:       coworkingNewPhone,
		CoworkingNewTariff:      coworkingNewTariff,
		CoworkingConfirm:        coworkingConfirm,
		CoworkingPending:        coworkingPending,
		CoworkingReturningPhone: coworkingReturningPhone,
		CoworkingAdminAction:    coworkingAdminAction,
		CoworkingProfile:        coworkingProfile,
		CoworkingTariff:         coworkingTariff,
	})

	reg.RegisterCommand("coworking", botengine.FlowEntry{
		Flow: FlowCoworking,
		Step: CoworkingHome,
	})

	reg.RegisterCallbackPrefix("cw:", botengine.FlowEntry{
		Flow: FlowCoworking,
		Step: CoworkingHome,
	})

	reg.RegisterCallbackPrefix("cwa:", botengine.FlowEntry{
		Flow: FlowCoworking,
		Step: CoworkingAdminAction,
	})

	reg.RegisterCallbackPrefix("cwconfirm:", botengine.FlowEntry{
		Flow: FlowCoworking,
		Step: CoworkingConfirm,
	})
}