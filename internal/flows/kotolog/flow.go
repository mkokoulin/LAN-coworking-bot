package flows

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowKotolog types.Flow = "kotolog"

	KotologHome types.Step = "kotolog:home"
	KotologList types.Step = "kotolog:list"
	KotologCat  types.Step = "kotolog:cat"
	KotologHelp types.Step = "kotolog:help"
	KotologDone types.Step = "kotolog:done"
)

func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowKotolog, map[types.Step]botengine.StepHandler{
		KotologHome: home,
		KotologList: list,
		KotologCat:  cat,
		KotologHelp: help,
		KotologDone: done,
	})

	// Команда /kotolog
	reg.RegisterCommand("kotolog", botengine.FlowEntry{Flow: FlowKotolog, Step: KotologHome})
	// Все callback-и с префиксом kotolog:* ведём в этот флоу
	reg.RegisterCallbackPrefix("kotolog:", botengine.FlowEntry{Flow: FlowKotolog, Step: KotologHome})
}
