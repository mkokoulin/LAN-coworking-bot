package botengine

import (
	"strings"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

func InterceptSlashNav(ev Event, ack func()) (bool, types.Step) {
	if ev.Kind == EventCallback && strings.HasPrefix(ev.CallbackData, "/") {
		ack()
		return true, InternalContinue
	}
	return false, ""
}
