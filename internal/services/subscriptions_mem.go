// internal/services/subscriptions_mem.go
package services

import (
    "context"
    "sync"
)

type MemSubscriptions struct {
    mu   sync.RWMutex
    subs map[int64]bool
}

func NewMemSubscriptions() *MemSubscriptions {
    return &MemSubscriptions{subs: map[int64]bool{}}
}

func (m *MemSubscriptions) SetWeeklyEvents(ctx context.Context, chatID int64, enabled bool) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    if enabled {
        m.subs[chatID] = true
    } else {
        delete(m.subs, chatID)
    }
    return nil
}

func (m *MemSubscriptions) ListWeeklyEventsSubscribers(ctx context.Context) ([]int64, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    out := make([]int64, 0, len(m.subs))
    for id := range m.subs {
        out = append(out, id)
    }
    return out, nil
}
