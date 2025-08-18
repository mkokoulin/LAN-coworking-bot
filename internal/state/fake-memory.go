package state

import (
	"sync"
	"time"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

type inMemory struct {
	mu   sync.RWMutex
	data map[int64]*types.Session
}

// ListDue implements Manager.
func (m *inMemory) ListDue(now time.Time) ([]int64, error) {
	panic("unimplemented")
}

// SetNextDigestAt implements Manager.
func (m *inMemory) SetNextDigestAt(chatID int64, next time.Time) error {
	panic("unimplemented")
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
