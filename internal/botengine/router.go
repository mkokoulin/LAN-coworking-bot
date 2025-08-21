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
