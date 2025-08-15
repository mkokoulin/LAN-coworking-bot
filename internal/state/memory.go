package state

import (
	"sync"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

// type Manager interface {
// 	Get(chatID int64) *types.Session
// 	Set(chatID int64, state *types.Session)
// 	Delete(chatID int64)
// 	// Reset()
// }

type memoryManager struct {
	store map[int64]*types.Session
	mu    sync.RWMutex
}

// ListSubscribedChatIDs implements Manager.
func (m *memoryManager) ListSubscribedChatIDs() ([]int64, error) {
	panic("unimplemented")
}

func New() Manager {
	return &memoryManager{
		store: make(map[int64]*types.Session),
	}
}

func (m *memoryManager) Get(chatID int64) *types.Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if state, ok := m.store[chatID]; ok {
		return state
	}

	m.mu.RUnlock()
	m.mu.Lock()
	defer m.mu.Unlock()
	state := &types.Session{}
	m.store[chatID] = state
	return state
}

func (m *memoryManager) Set(chatID int64, state *types.Session) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[chatID] = state
}

func (m *memoryManager) Delete(chatID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.store, chatID)
}

func (m *memoryManager) Reset() {
	// TODO
}
