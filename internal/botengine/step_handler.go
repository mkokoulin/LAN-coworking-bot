// internal/botengine/step_handler.go
package botengine

import (
	"context"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

// StepHandler — сигнатура обработчика шага FSM.
type StepHandler func(ctx context.Context, ev Event, d Deps, s *types.Session) (types.Step, error)
