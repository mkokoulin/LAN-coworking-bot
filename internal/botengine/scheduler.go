package botengine

import (
    "context"
    "log"
    "time"

    "github.com/mkokoulin/LAN-coworking-bot/internal/state"
    "github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

// локаль по умолчанию — Ереван
func userLoc(_ *types.Session) *time.Location {
    loc, err := time.LoadLocation("Asia/Yerevan")
    if err != nil {
        return time.FixedZone("Asia/Yerevan", 4*3600)
    }
    return loc
}

func computeNextRunUTC(hh, mm int, dow time.Weekday, loc *time.Location) time.Time {
    now := time.Now().In(loc)
    shift := (int(dow) - int(now.Weekday()) + 7) % 7
    cand := time.Date(now.Year(), now.Month(), now.Day(), hh, mm, 0, 0, loc).AddDate(0, 0, shift)
    if !cand.After(now) {
        cand = cand.AddDate(0, 0, 7)
    }
    return cand.UTC()
}

// RunWeeklyEvents теперь не «по понедельникам», а «по персональному расписанию».
// Он раз в 30 секунд опрашивает базу, кому уже пора (next_digest_at <= now),
// рассылает список и переносит next_digest_at на следующий раз.
func RunWeeklyEvents(ctx context.Context, d *Dispatcher, reg *Registry, states state.Manager, _ interface{}) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Println("[scheduler] stop")
            return
        case <-ticker.C:
            now := time.Now()
            ids, err := states.ListDue(now)
            if err != nil {
                log.Printf("[scheduler] ListDue error: %v", err)
                continue
            }
            if len(ids) == 0 {
                continue
            }

            for _, chatID := range ids {
                s := states.Get(chatID)
                if s == nil || !s.IsSubscribed {
                    // на всякий, чтобы не зациклиться
                    _ = states.SetNextDigestAt(chatID, now.Add(24*time.Hour))
                    continue
                }

                // шлём командой /events (через FSM), чтобы использовать уже готовый список
                ev := Event{Kind: EventCommand, ChatID: chatID, Command: "events"}
                deps := Deps{Bot: d.bot, Cfg: d.cfg, Svcs: d.svcs, Printer: d.printerFunc}

                // запустили FSM для этого чата
                if err := RunFSM(ctx, ev, reg, deps, s); err != nil {
                    log.Printf("[scheduler] fsm error chat=%d: %v", chatID, err)
                }

                // переносим следующий запуск согласно персональным настройкам пользователя
                loc := userLoc(s)
                next := computeNextRunUTC(s.EventsSubHour, s.EventsSubMinute, time.Weekday(s.EventsSubDOW), loc)
                if err := states.SetNextDigestAt(chatID, next); err != nil {
                    log.Printf("[scheduler] SetNextDigestAt chat=%d: %v", chatID, err)
                }
            }
        }
    }
}
