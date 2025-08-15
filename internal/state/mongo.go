package state

import (
	"context"
	"time"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoManager struct {
	collection *mongo.Collection
}

func NewMongo(uri, dbName, collectionName string) (Manager, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &mongoManager{collection: client.Database(dbName).Collection(collectionName)}, nil
}

func (m *mongoManager) ListSubscribedChatIDs() ([]int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
        if err := cur.Decode(&row); err != nil {
            return nil, err
        }
        ids = append(ids, row.ID)
    }
    return ids, cur.Err()
}

func (m *mongoManager) Get(chatID int64) *types.Session {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var s types.Session
	err := m.collection.FindOne(ctx, bson.M{"_id": chatID}).Decode(&s)
	if err == mongo.ErrNoDocuments {
		// Новый чистый сеанс
		s = types.Session{
			ChatID:        chatID,
			Flow:      "",
			Step:      "",
			Lang:      "", // пускай /language спросит
			Data:      map[string]interface{}{},
			Attempts:  0,
			ExpiresAt: time.Now().Add(30 * time.Minute),
		}
		m.Set(chatID, &s)
	} else if err != nil {
		// В случае ошибки вернём пустой, чтобы не падать
		s = types.Session{
			ChatID:   chatID,
			Data: map[string]interface{}{},
		}
	}

	// safety: если нет map — создадим
	if s.Data == nil {
		s.Data = map[string]interface{}{}
	}

	return &s
}

func (m *mongoManager) Set(chatID int64, s *types.Session) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	doc := bson.M{
		"_id":           chatID,
		"flow":          s.Flow,
		"step":          s.Step,
		"language":      s.Lang,
		"data":          s.Data,
		"attempts":      s.Attempts,
		"expires_at":    s.ExpiresAt,
		"is_authorized": s.IsAuthorized,
		"is_subscribed": s.IsSubscribed,
	}

	_, _ = m.collection.UpdateByID(
		ctx,
		chatID,
		bson.M{"$set": doc},
		options.Update().SetUpsert(true),
	)
}

func (m *mongoManager) Delete(chatID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, _ = m.collection.DeleteOne(ctx, bson.M{"_id": chatID})
}
