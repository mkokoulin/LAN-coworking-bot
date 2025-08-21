// Path: internal/botengine/fsm.go  (файл где RunFSM)
package botengine

import (
	"context"
	"fmt"
	"sort"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

// InternalContinue — служебный шаг, сигнализирует FSM "продолжить немедленно".
const InternalContinue types.Step = "__internal:continue"

func RunFSM(ctx context.Context, ev Event, reg *Registry, d Deps, s *types.Session) error {
	GuardMaybeReset(ev, reg, s)

	const maxInternalIters = 5

	for it := 0; it < maxInternalIters; it++ {

		// 0) Callback со слеш-командой → считаем навигацией по команде
		if ev.Kind == EventCallback && strings.HasPrefix(ev.CallbackData, "/") {
			// ACK, чтобы не висели "часики" на кнопке
			if ev.CallbackQueryID != "" {
				cb := tgbotapi.NewCallback(ev.CallbackQueryID, "")
				_, _ = d.Bot.Request(cb)
			}
			cmd := ev.Command
			if cmd == "" { // на всякий случай
				cmd = normalizeCommand(ev.CallbackData)
			}
			if entry, ok := reg.commands[cmd]; ok {
				s.Flow, s.Step = entry.Flow, entry.Step
			}
			// Переходим к отрисовке нового шага (тот же ev)
			// не return — рендерим ниже как обычный шаг
		} else
		// ✅ 0.5) Форс-вход по командам из сообщений
		if ev.Kind == EventCommand {
			if entry, ok := reg.ResolveEntry(ev); ok {
				s.Flow, s.Step = entry.Flow, entry.Step
			} else {
				// Неизвестная команда — мягко сообщаем и не ломаем текущий сценарий
				cmds := make([]string, 0, len(reg.commands))
				for c := range reg.commands {
					cmds = append(cmds, "/"+c)
				}
				sort.Strings(cmds)
				var hint string
				if len(cmds) > 0 {
					hint = "\nДоступные команды: " + strings.Join(cmds, ", ")
				}
				text := fmt.Sprintf("❓ Неизвестная команда: /%s.%s", ev.Command, hint)
				_, _ = d.Bot.Send(tgbotapi.NewMessage(s.ChatID, text))
				return nil
			}
		} else
		// 1) Если текущий шаг пустой/битый — пробуем авто-вход по событию
		if s.Flow == "" || s.Step == "" || !reg.Has(s.Flow, s.Step) {
			if entry, ok := reg.ResolveEntry(ev); ok {
				s.Flow, s.Step = entry.Flow, entry.Step
			} else {
				return nil
			}
		}

		// 2) Берём хендлер текущего шага
		steps, ok := reg.flows[s.Flow]
		if !ok {
			return nil
		}
		h, ok := steps[s.Step]
		if !ok {
			return nil
		}

		prevFlow, prevStep := s.Flow, s.Step

		// 3) Выполняем шаг
		next, err := h(ctx, ev, d, s)
		if err != nil {
			return err
		}

		// 4) InternalContinue
		if next == InternalContinue {
			if s.Flow == prevFlow && s.Step == prevStep {
				return nil
			}
			continue
		}

		// 5) Обычный переход
		s.Step = next
		return nil
	}

	return nil
}
