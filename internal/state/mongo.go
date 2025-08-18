package state

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoManager struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoManager(ctx context.Context, uri, dbName, collectionName string) (*mongoManager, error) {
	cl, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	col := cl.Database(dbName).Collection(collectionName)

	// Индексы под рассылку
	_, _ = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "is_subscribed", Value: 1}}},
		{Keys: bson.D{{Key: "next_digest_at", Value: 1}}},
	})

	return &mongoManager{client: cl, collection: col}, nil
}

func (m *mongoManager) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *mongoManager) Get(chatID int64) *types.Session {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var s types.Session
	err := m.collection.FindOne(ctx, bson.M{"_id": chatID}).Decode(&s)
	if err == mongo.ErrNoDocuments {
		return &types.Session{
			ChatID: chatID,
			Data:   map[string]interface{}{},
		}
	}
	if err != nil {
		log.Printf("[state.mongo] Get(%d) error: %v", chatID, err)
		return &types.Session{
			ChatID: chatID,
			Data:   map[string]interface{}{},
		}
	}
	if s.Data == nil {
		s.Data = map[string]interface{}{}
	}
	// Подстрахуем ChatID, если теги не совпадают
	if s.ChatID == 0 {
		s.ChatID = chatID
	}
	return &s
}

func (m *mongoManager) Set(chatID int64, s *types.Session) {
	if s == nil {
		return
	}
	if s.Data == nil {
		s.Data = map[string]interface{}{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	doc := bson.M{
		"_id":  chatID,
		"lang": s.Lang, // чаще всего поле тегируется как `bson:"lang"`
		"flow": s.Flow,
		"step": s.Step,
		"data": s.Data,

		// Подписка на анонсы
		"is_subscribed":   s.IsSubscribed,
		"events_sub_dow":  s.EventsSubDOW,    // 0..6
		"events_sub_hour": s.EventsSubHour,   // 0..23
		"events_sub_min":  s.EventsSubMinute, // 0..59
		"next_digest_at":  s.NextDigestAt,    // UTC
	}

	_, err := m.collection.UpdateByID(ctx, chatID, bson.M{"$set": doc}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("[state.mongo] Set(%d) error: %v", chatID, err)
	}
}

func (m *mongoManager) Delete(chatID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, _ = m.collection.DeleteOne(ctx, bson.M{"_id": chatID})
}

func (m *mongoManager) ListSubscribedChatIDs() ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cur, err := m.collection.Find(ctx, bson.M{"is_subscribed": true}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var ids []int64
	for cur.Next(ctx) {
		var row struct {
			ID int64 `bson:"_id"`
		}
		if err := cur.Decode(&row); err == nil {
			ids = append(ids, row.ID)
		}
	}
	return ids, cur.Err()
}

// ===== Методы для планировщика рассылки =====

// ListDue — кому пора слать (next_digest_at <= now, подписка активна)
func (m *mongoManager) ListDue(now time.Time) ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"is_subscribed":  true,
		"next_digest_at": bson.M{"$lte": now.UTC()},
	}
	cur, err := m.collection.Find(ctx, filter, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var ids []int64
	for cur.Next(ctx) {
		var row struct {
			ID int64 `bson:"_id"`
		}
		if err := cur.Decode(&row); err == nil {
			ids = append(ids, row.ID)
		}
	}
	return ids, cur.Err()
}

// SetNextDigestAt — атомарно переносит следующий запуск
func (m *mongoManager) SetNextDigestAt(chatID int64, next time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.collection.UpdateByID(ctx, chatID, bson.M{
		"$set": bson.M{"next_digest_at": next.UTC()},
	})
	return err
}

// NewFromEnv — выбрать менеджер из env или память
// MONGO_URI=... (обязательно для включения Mongo)
// MONGO_DB=lan-bot (по умолчанию)
// MONGO_COLLECTION=sessions (по умолчанию)
func NewFromEnv(ctx context.Context) (Manager, func(context.Context) error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return NewMemoryManager(), func(context.Context) error { return nil }
	}
	db := os.Getenv("MONGO_DB")
	if db == "" {
		db = "lan-bot"
	}
	col := os.Getenv("MONGO_COLLECTION")
	if col == "" {
		col = "sessions"
	}

	m, err := NewMongoManager(ctx, uri, db, col)
	if err != nil {
		log.Printf("[state] cannot init mongo, fallback to memory: %v", err)
		return NewMemoryManager(), func(context.Context) error { return nil }
	}
	log.Printf("[state] using mongo %s/%s", db, col)
	return m, m.Close
}
