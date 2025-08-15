package botengine

import "sync"
import "github.com/mkokoulin/LAN-coworking-bot/internal/types"

type SessionStore struct {
    mu   sync.RWMutex
    byID map[int64]*types.Session
}

func NewSessionStore() *SessionStore {
    return &SessionStore{byID: make(map[int64]*types.Session)}
}

func (s *SessionStore) Get(chatID int64) *types.Session {
    s.mu.RLock()
    sess, ok := s.byID[chatID]
    s.mu.RUnlock()
    if ok && sess != nil {
        return sess
    }
    s.mu.Lock()
    defer s.mu.Unlock()
    // двойная проверка
    if sess, ok := s.byID[chatID]; ok && sess != nil {
        return sess
    }
    sess = &types.Session{ChatID: chatID, Data: map[string]interface{}{}}
    s.byID[chatID] = sess
    return sess
}

func (s *SessionStore) Save(sess *types.Session) {
    if sess == nil { return }
    s.mu.Lock()
    s.byID[sess.ChatID] = sess
    s.mu.Unlock()
}
