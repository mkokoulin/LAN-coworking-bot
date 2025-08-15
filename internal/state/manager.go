package state

import "github.com/mkokoulin/LAN-coworking-bot/internal/types"

type Manager interface {
	Get(chatID int64) *types.Session
	Set(chatID int64, state *types.Session)
	Delete(chatID int64)
	ListSubscribedChatIDs() ([]int64, error)
}
