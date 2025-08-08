package state

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommandLogger struct {
	collection *mongo.Collection
}

func NewCommandLogger(db *mongo.Database, collectionName string) *CommandLogger {
	return &CommandLogger{
		collection: db.Collection(collectionName),
	}
}

func (l *CommandLogger) Log(chatID int64, command string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	doc := bson.M{
		"chat_id":   chatID,
		"command":   command,
		"timestamp": time.Now(),
	}

	_, err := l.collection.InsertOne(ctx, doc)
	return err
}
