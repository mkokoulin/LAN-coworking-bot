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

	// небольшой лимит на внутренние "продолжить", чтобы не зациклиться
	const maxInternalIters = 5

	for it := 0; it < maxInternalIters; it++ {

		// ✅ 0) Форс-вход по командам (перебивает любой текущий флоу)
		if ev.Kind == EventCommand {
			if entry, ok := reg.ResolveEntry(ev); ok {
				s.Flow, s.Step = entry.Flow, entry.Step
			} else {
				// Неизвестная команда — мягко сообщаем и НЕ ломаем текущий сценарий
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
				// не нашли, мягко выходим
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

		// 4) Обработка служебного шага — "продолжить немедленно"
		if next == InternalContinue {
			// если хендлер не поменял шаг/флоу — нечего продолжать
			if s.Flow == prevFlow && s.Step == prevStep {
				return nil
			}
			// иначе крутим следующую итерацию с тем же событием
			continue
		}

		// 5) Обычный переход на следующий шаг
		s.Step = next
		return nil
	}

	// превысили лимит внутренних итераций — на всякий случай выходим
	return nil
}
