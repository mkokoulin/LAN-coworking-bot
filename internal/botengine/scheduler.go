package botengine

import (
	"context"
	"log"
	"time"

	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/state"
)

func RunWeeklyEvents(ctx context.Context, d *Dispatcher, reg *Registry, states state.Manager, cfg *config.Config) {
    targetWeekday := time.Monday
    targetHour := 0
    targetMinute := 58

    loc, _ := time.LoadLocation("Asia/Yerevan")
    if loc == nil { loc = time.Local }

    ticker := time.NewTicker(time.Minute)
    go func() {
        defer ticker.Stop()
        for {
            select {
            case <-ctx.Done():
                return
            case now := <-ticker.C:
                now = now.In(loc)
                if now.Weekday() != targetWeekday || now.Hour() != targetHour || now.Minute() != targetMinute {
                    continue
                }

                ids, err := states.ListSubscribedChatIDs() // 👈 читаем из Mongo
                if err != nil || len(ids) == 0 {
                    if err != nil { log.Printf("[weekly-events] list error: %v", err) }
                    continue
                }

                for _, chatID := range ids {
                    s := states.Get(chatID)
                    if s.ChatID == 0 { s.ChatID = chatID }

                    ev := Event{Kind: EventCommand, ChatID: chatID, Command: "events"}
                    deps := Deps{Bot: d.bot, Cfg: d.cfg, Svcs: d.svcs, Printer: d.printerFunc}

                    func() {
                        defer states.Set(chatID, s)
                        if err := RunFSM(ctx, ev, reg, deps, s); err != nil {
                            log.Printf("[weekly-events] fsm error chat=%d: %v", chatID, err)
                        }
                    }()
                }
            }
        }
    }()
}
