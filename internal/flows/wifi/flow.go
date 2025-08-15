package flows

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowWifi types.Flow = "wifi"

	WifiStart      types.Step = "wifi:start"
	WifiWaitChoice types.Step = "wifi:wait_choice"
	WifiWaitCode   types.Step = "wifi:wait_code"
	WifiDone       types.Step = "wifi:done"
)


func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowWifi, map[types.Step]botengine.StepHandler{
		WifiStart:      start,
		WifiWaitChoice: waitChoice,
		WifiWaitCode:   waitCode,
		WifiDone:       done,
	})

	// входы в сценарий
	reg.RegisterCommand("wifi", botengine.FlowEntry{Flow: FlowWifi, Step: WifiStart})
	reg.RegisterCallbackPrefix("wifi:", botengine.FlowEntry{Flow: FlowWifi, Step: WifiWaitChoice})
}
