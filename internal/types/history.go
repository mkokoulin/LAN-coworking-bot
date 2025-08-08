package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommandLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ChatID    int64              `bson:"chat_id"`
	Command   string             `bson:"command"`
	Timestamp time.Time          `bson:"timestamp"`
}