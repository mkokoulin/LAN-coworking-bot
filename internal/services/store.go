package services

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

type Store struct {
	*FireDB
}

func NewStore(fireDB *FireDB) *Store {	
	return &Store{
		FireDB: fireDB,
	}
}

func (s *Store) Create(cs *types.ChatStorage) (*types.ChatStorage, error) {
	if err := s.NewRef("chats/" + strconv.Itoa(int(cs.ChatID))).Set(context.Background(), cs); err != nil {
		return nil, err
	}

	return cs, nil
}

func (s *Store) CreateIfNotExists(cs *types.ChatStorage) (*types.ChatStorage, error) {
	c, err := s.GetByID(cs.ChatID)
	if err != nil {
		return nil, err
	}

	if c != nil {
		return c, nil
	}

	return s.Create(cs)
}
   
//    func (s *Store) Delete(b *Bin) error {
// 	return s.NewRef("bins/" + b.BIN).Delete(context.Background())
//    }
   
func (s *Store) GetByID(chatID int64) (*types.ChatStorage, error) {
	chats := types.ChatStorage{}
	
	if err := s.NewRef("chats/" + strconv.Itoa(int(chatID))).Get(context.Background(), &chats); err != nil {
		return nil, err
	}

	if chats.ChatID == 0 {
		return nil, nil
	}
	
	return &chats, nil
}
   
func (s *Store) Update(chatID int64, cs types.ChatStorage) error {
	var chatStateInInterface map[string]interface{}
	
    chatState, _ := json.Marshal(cs)
    json.Unmarshal(chatState, &chatStateInInterface)
	
	return s.NewRef("chats/" + strconv.Itoa(int(chatID))).Update(context.Background(), chatStateInInterface)
}