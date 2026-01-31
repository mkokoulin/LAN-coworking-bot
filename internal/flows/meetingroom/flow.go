package flows

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowMeeting       types.Flow = "meetingroom"
	MeetPrompt        types.Step = "meeting:prompt"
	MeetWaitInterval  types.Step = "meeting:wait_interval"
	MeetNotify        types.Step = "meeting:notify"
	MeetDone          types.Step = "meeting:done"
)

func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowMeeting, map[types.Step]botengine.StepHandler{
		MeetPrompt:       prompt,
		MeetWaitInterval: waitInterval,
		MeetDone:         done,
	})

	reg.RegisterCommand("meetingroom", botengine.FlowEntry{Flow: FlowMeeting, Step: MeetPrompt})
}
