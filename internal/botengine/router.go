package botengine

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

func StartFlow(s *types.Session, flow, step string) {
	if s == nil {
		return
	}
	if s.Data == nil {
		s.Data = map[string]interface{}{}
	}
	s.Flow = types.Flow(flow)
	s.Step = types.Step(step)
}

// startFlowByCommand — маппинг /команд на стартовые Flow/Step.
// Держу как пример, без ссылок на несуществующие поля.
// func StartFlowByCommand(s *types.Session, cmd string) {
// 	if s == nil {
// 		return
// 	}
// 	if s.Data == nil {
// 		s.Data = map[string]interface{}{}
// 	}
//
// 	switch cmd {
// 	case "start":
// 		StartFlow(s, "start", "start:show")
// 	case "language":
// 		StartFlow(s, "language", "language:prompt")
// 	case "wifi":
// 		StartFlow(s, "wifi", "wifi:start")
// 	case "booking":
// 		StartFlow(s, "booking", "booking:info")
// 	case "meetingroom":
// 		StartFlow(s, "meetingroom", "wait_interval")
// 	case "printout":
// 		StartFlow(s, "printout", "printout:show")
// 	case "events":
// 		StartFlow(s, "events", "events:intro")
// 	case "about":
// 		StartFlow(s, "about", "about:send")
// 	case "menu":
// 		StartFlow(s, "menu", "menu:send")
// 	default:
// 		// неизвестная команда — мягкий сброс
// 		StartFlow(s, "", "")
// 	}
// }
