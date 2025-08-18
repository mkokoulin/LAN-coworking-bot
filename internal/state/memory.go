package state

import (
	"sync"
	"time"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

type memoryManager struct {
	mu    sync.RWMutex
	store map[int64]*types.Session
}

func NewMemoryManager() Manager {
	return &memoryManager{store: make(map[int64]*types.Session)}
}

func (m *memoryManager) Get(chatID int64) *types.Session {
	m.mu.RLock()
	sess, ok := m.store[chatID]
	m.mu.RUnlock()

	if ok && sess != nil {
		// гарантируем не-nil Data
		if sess.Data == nil {
			sess.Data = map[string]interface{}{}
		}
		return sess
	}

	// создаём сессию по умолчанию
	ns := &types.Session{
		ChatID: chatID,
		Data:   map[string]interface{}{},
		// поля подписки оставляем нулевыми значениями —
		// IsSubscribed=false, NextDigestAt zero => не шлём
	}
	m.Set(chatID, ns)
	return ns
}

func (m *memoryManager) Set(chatID int64, state *types.Session) {
	if state == nil {
		return
	}
	if state.Data == nil {
		state.Data = map[string]interface{}{}
	}
	m.mu.Lock()
	m.store[chatID] = state
	m.mu.Unlock()
}

func (m *memoryManager) Delete(chatID int64) {
	m.mu.Lock()
	delete(m.store, chatID)
	m.mu.Unlock()
}

// ListSubscribedChatIDs — все подписанные (для массовых операций/миграций)
func (m *memoryManager) ListSubscribedChatIDs() ([]int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var ids []int64
	for id, s := range m.store {
		if s != nil && s.IsSubscribed {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// ListDue — те, кому пора слать (NextDigestAt задан и <= now).
// Если NextDigestAt = zero или пользователь не подписан — пропускаем.
func (m *memoryManager) ListDue(now time.Time) ([]int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var due []int64
	for id, s := range m.store {
		if s == nil || !s.IsSubscribed {
			continue
		}
		if s.NextDigestAt.IsZero() {
			continue
		}
		if !s.NextDigestAt.After(now.UTC()) {
			due = append(due, id)
		}
	}
	return due, nil
}

func (m *memoryManager) SetNextDigestAt(chatID int64, next time.Time) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    if s, ok := m.store[chatID]; ok && s != nil {
        s.NextDigestAt = next.UTC()
    }
    return nil
}
