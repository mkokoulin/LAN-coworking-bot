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

	// Команды
	reg.RegisterCommand("kotolog", botengine.FlowEntry{Flow: FlowKotolog, Step: KotologHome})
	reg.RegisterCommand("cats",    botengine.FlowEntry{Flow: FlowKotolog, Step: KotologHome})

	// ⬇️ ВАЖНО: конкретные префиксы вместо общего "kotolog:"
	reg.RegisterCallbackPrefix("kotolog:home",  botengine.FlowEntry{Flow: FlowKotolog, Step: KotologHome})
	reg.RegisterCallbackPrefix("kotolog:list:", botengine.FlowEntry{Flow: FlowKotolog, Step: KotologList})
	reg.RegisterCallbackPrefix("kotolog:cat:",  botengine.FlowEntry{Flow: FlowKotolog, Step: KotologCat})
	reg.RegisterCallbackPrefix("kotolog:help",  botengine.FlowEntry{Flow: FlowKotolog, Step: KotologHelp})
}
