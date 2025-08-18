package types

import "time"


///////////
type Session struct {
	Data   map[string]interface{} `bson:"data,omitempty" json:"data,omitempty"`

	ChatID int64  `bson:"_id,omitempty"  json:"chatID,omitempty"` // можно оставить "chatID", если хочешь, но с _id декодинг надёжнее
	UserID int64  `bson:"userID,omitempty" json:"userID,omitempty"`
	Lang   string `bson:"lang,omitempty"   json:"lang,omitempty"`
	Flow   Flow   `bson:"flow,omitempty"   json:"flow,omitempty"`
	Step   Step   `bson:"step,omitempty"   json:"step,omitempty"`

	// поля подписки (ВАЖНО: snake_case)
	IsSubscribed    bool      `bson:"is_subscribed,omitempty"   json:"is_subscribed,omitempty"`
	EventsSubDOW    int       `bson:"events_sub_dow,omitempty"  json:"events_sub_dow,omitempty"`   // 0..6
	EventsSubHour   int       `bson:"events_sub_hour,omitempty" json:"events_sub_hour,omitempty"`  // 0..23
	EventsSubMinute int       `bson:"events_sub_min,omitempty"  json:"events_sub_min,omitempty"`   // 0..59
	NextDigestAt    time.Time `bson:"next_digest_at,omitempty"  json:"next_digest_at,omitempty"`   // UTC
}
