package botengine

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

func StartFlow(s *types.Session, flow, step string) {
    if s.Data == nil {
        s.Data = map[string]interface{}{}
    }
    s.Attempts = 0
    s.Flow = types.Flow(flow)
    s.Step = types.Step(step)
}

// startFlowByCommand — маппинг /команд на стартовые Flow/Step.
// Здесь используем строковые константы, чтобы не плодить импортные циклы.
// Они должны соответствовать константам в пакетах flows/*.
// func StartFlowByCommand(s *types.Session, cmd string) {
// 	// сброс/инициализация контекста сценария
// 	if s.Data == nil {
// 		s.Data = map[string]interface{}{}
// 	}
// 	s.Attempts = 0
// 	// опционально: срок жизни сессии, если используешь
// 	// s.ExpiresAt = time.Now().Add(30 * time.Minute)

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
