package flows

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowLanguage   types.Flow = "language"
	LangPrompt     types.Step = "language:prompt"
	LangWaitChoice types.Step = "language:wait_choice"
	LangDone       types.Step = "language:done"
)

func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowLanguage, map[types.Step]botengine.StepHandler{
		LangPrompt:     prompt,
		LangWaitChoice: waitChoice,
		LangDone:       done,
	})

	// входы в сценарий языка
	reg.RegisterCommand("language", botengine.FlowEntry{Flow: FlowLanguage, Step: LangPrompt})
	// можно добавить алиас:
	// reg.RegisterCommand("lang", botengine.FlowEntry{Flow: FlowLanguage, Step: LangPrompt})

	// все callback-и lang:* обрабатываем этим флоу
	reg.RegisterCallbackPrefix("lang:", botengine.FlowEntry{Flow: FlowLanguage, Step: LangWaitChoice})
}
