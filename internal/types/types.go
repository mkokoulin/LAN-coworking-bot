package types

import "time"


///////////
type Session struct {
    ChatID    int64             `bson:"_id"`
    Lang      string            `bson:"language"`
    Flow      Flow              `bson:"flow"`
    Step      Step              `bson:"step"`
    Data      map[string]interface{} `bson:"data,omitempty"`
    Attempts  int               `bson:"attempts,omitempty"`
    ExpiresAt time.Time         `bson:"expires_at,omitempty"`
    IsAuthorized bool           `bson:"is_authorized"`
    PendingCmd string           `bson:"pending_cmd,omitempty"`
    IsSubscribed bool           `bson:"is_subscribed,omitempty"`
}
