package botengine

import (
	"context"
	"strings"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

// isCommand извлекает команду без аргументов: "/bar arg" -> "bar".
func isCommand(text string) (string, bool) {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "/") {
		return "", false
	}
	cmd := text[1:]
	// отрезаем @botname и аргументы
	if i := strings.IndexAny(cmd, " @"); i >= 0 {
		cmd = cmd[:i]
	}
	if i := strings.IndexByte(cmd, ' '); i >= 0 {
		cmd = cmd[:i]
	}
	cmd = strings.ToLower(strings.TrimSpace(cmd))
	if cmd == "" {
		return "", false
	}
	return cmd, true
}

// // clearNamespace удаляет все KV-ключи с префиксом.
// func clearNamespace(s *types.Session, prefix string) {
// 	if s == nil || s.KV == nil {
// 		return
// 	}
// 	for k := range s.KV {
// 		if strings.HasPrefix(k, prefix) {
// 			delete(s.KV, k)
// 		}
// 	}
// }

// HardResetToStart — глобальный сброс: очищаем флоу/шаг/временные ключи и переводим в старт.
func HardResetToStart(ctx context.Context, s *types.Session) {
	if s == nil {
		return
	}
	s.Flow = ""
	s.Step = ""
	// // Чистим всё «временное» по известным неймспейсам (добавляй свои)
	// clearNamespace(s, "bar:")
	// clearNamespace(s, "booking:")
	// clearNamespace(s, "lang:")
	// // можно добавить флаг "fresh_start" в KV — если используешь где-то в /start
	// if s.KV != nil {
	// 	s.KV["fresh_start"] = "1"
	// }
}

// GuardOnNewCommand применяет политику:
// 1) Повторная /bar при незавершённом шаге (/bar) → чистим bar:* и перезапускаем bar.
// 2) Любая новая команда при незавершённом шаге → жёсткий сброс и переходим в /start.
// Возвращает действие: "" | "restart_bar" | "force_start"
func GuardOnNewCommand(text string, s *types.Session) string {
	text = strings.TrimSpace(text)
	if text == "" || text[0] != '/' {
		return ""
	}

	// вычистим атрибуты типа /cmd@bot и хвост после пробела
	cmd := text[1:]
	if i := strings.IndexByte(cmd, ' '); i >= 0 {
		cmd = cmd[:i]
	}
	if j := strings.IndexByte(cmd, '@'); j >= 0 {
		cmd = cmd[:j]
	}

	// Явное правило: новая /bar всегда перезапускает бар
	if cmd == "bar" {
		s.ResetFlow() // <= только Flow/Step
		return "restart_bar"
	}

	// Если мы в процессе какого-то флоу — любая другая команда => /start
	if s.Flow != "" || s.Step != "" {
		s.ResetFlow()
		return "force_start"
	}
	return ""
}

func GuardMaybeReset(ev Event, reg *Registry, s *types.Session) {
	// если мы и так ни в каком флоу — нечего сбрасывать
	if s.Flow == "" && s.Step == "" {
		return
	}
	// проверяем, распознаётся ли событие как вход в флоу (команда/коллбэк)
	if _, ok := reg.ResolveEntry(ev); !ok {
		return
	}
	// сбрасываем маршрут; дальше авто-вход отработает по этому же ev
	s.ResetFlow()

	// Если хочешь вместо "войти в присланную команду" всегда уводить в /start,
	// раскомментируй 3 строки ниже и подставь свои константы флоу/шага старта:
	// s.Flow = flows.FlowStart
	// s.Step = flows.StartPrompt
	// (и не забудь импортировать пакет flows)
}
