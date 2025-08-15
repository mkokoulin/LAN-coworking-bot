package state

import (
	"sync"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

type inMemory struct {
	mu   sync.RWMutex
	data map[int64]*types.Session
}

// ListSubscribedChatIDs implements Manager.
func (m *inMemory) ListSubscribedChatIDs() ([]int64, error) {
	panic("unimplemented")
}

func (m *inMemory) Delete(chatID int64) {
	panic("unimplemented")
}

func NewInMemory() Manager {
	return &inMemory{
		data: make(map[int64]*types.Session),
	}
}

func (m *inMemory) Get(chatID int64) *types.Session {
	m.mu.RLock()
	s, ok := m.data[chatID]
	m.mu.RUnlock()
	if ok && s != nil {
		return s
	}

	ns := &types.Session{ChatID: chatID}
	m.mu.Lock()
	m.data[chatID] = ns
	m.mu.Unlock()
	return ns
}

func (m *inMemory) Set(chatID int64, s *types.Session) {
	if s == nil {
		return
	}
	m.mu.Lock()
	m.data[chatID] = s
	m.mu.Unlock()
}
