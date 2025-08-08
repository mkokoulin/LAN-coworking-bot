package state

import (
	"context"
	"sync"
	"time"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoManager struct {
	collection *mongo.Collection
	mu         sync.RWMutex
}

func NewMongo(uri, dbName, collectionName string) (Manager, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &mongoManager{collection: collection}, nil
}

func (m *mongoManager) Get(chatID int64) *types.ChatStorage {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var result types.ChatStorage
	err := m.collection.FindOne(ctx, bson.M{"_id": chatID}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		result = types.ChatStorage{}
		m.Set(chatID, &result)
	} else if err != nil {
		return &types.ChatStorage{}
	}

	return &result
}

func (m *mongoManager) Set(chatID int64, state *types.ChatStorage) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	doc := bson.M{
		"_id":                    chatID,
		"language":              state.Language,
		"is_authorized":          state.IsAuthorized,
		"is_booking_process":      state.IsBookingProcess,
		"is_wifi_process":         state.IsWifiProcess,
		"is_awaiting_confirmation": state.IsAwaitingConfirmation,
		"current_command":        state.CurrentCommand,
	}
	_, _ = m.collection.UpdateByID(ctx, chatID, bson.M{"$set": doc}, options.Update().SetUpsert(true))
}

func (m *mongoManager) Delete(chatID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, _ = m.collection.DeleteOne(ctx, bson.M{"_id": chatID})
}
