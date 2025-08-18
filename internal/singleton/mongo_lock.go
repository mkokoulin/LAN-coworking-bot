package singleton

import (
	"context"
	"encoding/hex"
	"errors"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoLocker struct {
	coll *mongo.Collection
}

func NewMongoLocker(coll *mongo.Collection) *mongoLocker {
	return &mongoLocker{coll: coll}
}

var ErrLostLock = errors.New("lock lost")

type lockDoc struct {
	ID       string    `bson:"_id"`
	Owner    string    `bson:"owner"`
	ExpireAt time.Time `bson:"expireAt"`
	Created  time.Time `bson:"createdAt,omitempty"`
}

// Вынесите ensureTTLIndex в init вашего приложения (один раз на старте).
func ensureTTLIndex(ctx context.Context, coll *mongo.Collection) error {
	idx := mongo.IndexModel{
		Keys:    bson.D{{Key: "expireAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	_, err := coll.Indexes().CreateOne(ctx, idx)
	return err
}

func newOwnerID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// небольшой джиттер, чтобы не билось в одну секунду
func withJitter(d time.Duration, frac float64) time.Duration {
	if frac <= 0 {
		return d
	}
	delta := time.Duration(frac * float64(d))
	return d - delta + time.Duration(rand.Int63n(int64(2*delta)+1))
}

// WaitAcquire — ждём, пока сможем захватить лок. Возвращает release().
func (l *mongoLocker) WaitAcquire(ctx context.Context, key string, ttl time.Duration) (func() error, error) {
	// !!! В идеале вызвать ensureTTLIndex один раз при старте приложения.
	// Оставлено здесь на случай, если вы ещё не вынесли это.
	if err := ensureTTLIndex(ctx, l.coll); err != nil {
		return nil, err
	}

	// РЕКОМЕНДАЦИИ: ttl ~ 2-5s. Чем меньше — тем быстрее восстановление после падения.
	// Защита от совсем микроскопического ttl
	if ttl < 600*time.Millisecond {
		ttl = 600 * time.Millisecond
	}

	owner := newOwnerID()

	for {
		ok, wait, err := l.tryLock(ctx, key, ttl, owner)
		if err == nil && ok {
			stop := make(chan struct{})

			// Heartbeat чаще: ttl/3 (+джиттер), чтобы надёжно удерживать короткие TTL.
			renewEvery := max(ttl / 3, 200 * time.Millisecond)

			go func() {
				timer := time.NewTimer(withJitter(renewEvery, 0.2))
				defer timer.Stop()
				for {
					select {
					case <-timer.C:
						if l.renew(context.Background(), key, owner, ttl) != nil {
							// Потеряли лок — просто перестаём продлевать.
							return
						}
						timer.Reset(withJitter(renewEvery, 0.2))
					case <-stop:
						return
					}
				}
			}()

			release := func() error {
				close(stop)
				_, _ = l.coll.DeleteOne(context.Background(), bson.M{"_id": key, "owner": owner})
				return nil
			}
			return release, nil
		}

		// Транзиентная ошибка — короткий бэкофф
		if err != nil && wait <= 0 {
			wait = 100 * time.Millisecond // было 500ms
		}

		// Немного джиттера, чтобы конкуренты не просыпались синхронно
		wait = withJitter(wait, 0.2)

		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (l *mongoLocker) tryLock(ctx context.Context, key string, ttl time.Duration, owner string) (ok bool, wait time.Duration, err error) {
	now := time.Now()
	upd := bson.M{
		"$setOnInsert": bson.M{"createdAt": now},
		"$set":         bson.M{"owner": owner, "expireAt": now.Add(ttl)},
	}
	filter := bson.M{
		"_id": key,
		"$or": []bson.M{
			{"expireAt": bson.M{"$lte": now}}, // истёкший лок — можно "забрать"
			{"owner": owner},                  // наш — продлеваем
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var out lockDoc
	err = l.coll.FindOneAndUpdate(ctx, filter, upd, opts).Decode(&out)
	if err == nil {
		return true, 0, nil
	}

	// Не взяли — узнаём оставшийся TTL
	var cur lockDoc
	if getErr := l.coll.FindOne(ctx, bson.M{"_id": key}).Decode(&cur); getErr != nil {
		// ключа нет — быстро ретраим
		return false, 50 * time.Millisecond, nil // было 300ms
	}

	left := max(
// Минимальный интервал ожидания чуть ниже — реактивнее при коротких TTL
time.Until(cur.ExpireAt), 
// было 200ms
50 * time.Millisecond)
	return false, left, nil
}

func (l *mongoLocker) renew(ctx context.Context, key, owner string, ttl time.Duration) error {
	now := time.Now()
	res := l.coll.FindOneAndUpdate(
		ctx,
		bson.M{"_id": key, "owner": owner},
		bson.M{"$set": bson.M{"expireAt": now.Add(ttl)}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var d lockDoc
	if err := res.Decode(&d); err != nil {
		return ErrLostLock
	}
	return nil
}
