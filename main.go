package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/text/message"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/flows"
	"github.com/mkokoulin/LAN-coworking-bot/internal/locales"
	"github.com/mkokoulin/LAN-coworking-bot/internal/services"
	"github.com/mkokoulin/LAN-coworking-bot/internal/state"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

func main() {
	go func() {
		_ = http.ListenAndServe(":"+os.Getenv("PORT"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("ok"))
		}))
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	locales.Init()

	// 1) Конфиг
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("[boot] load config: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("[boot] telegram: %v", err)
	}
	log.Printf("Bot started as @%s (debug=%v)", bot.Self.UserName, bot.Debug)

	lockDisabled := os.Getenv("LOCK_DISABLE") == "1"
	lockID := "telegram_updates_lock:" + bot.Self.UserName
	var release func() error = func() error { return nil }

	if !lockDisabled {
		mongoDB := "coworking_bot"
		mongoColl := "locks"

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
		if err != nil {
			log.Fatalf("[singleton] mongo connect: %v", err)
		}
		coll := client.Database(mongoDB).Collection(mongoColl)
		if err := ensureTTLIndex(ctx, coll); err != nil {
			log.Fatalf("[singleton] ensure TTL index: %v", err)
		}

		// Аварийный сброс лока (если прошлый владелец умер без TTL/GC)
		if os.Getenv("LOCK_FORCE") == "1" {
			if _, err := coll.DeleteOne(ctx, bson.M{"_id": lockID}); err != nil {
				log.Fatalf("[singleton] force release failed: %v", err)
			}
			log.Printf("[singleton] force-released %s", lockID)
		}

		// Логируем текущего владельца (если есть)
		var cur lockDoc
		if err := coll.FindOne(ctx, bson.M{"_id": lockID}).Decode(&cur); err == nil {
			log.Printf("[singleton] current lock owner: %s (expires %s)", cur.Owner, cur.ExpireAt.Format(time.RFC3339))
		}

		// Ждём лок (TTL=3m, heartbeat каждые 90s)
		release, err = mongoWaitAcquire(ctx, coll, lockID, 10*time.Second)
		if err != nil {
			log.Fatalf("[singleton] cannot acquire lock: %v", err)
		}
		defer func() { _ = release() }()
		log.Println("[singleton] lock acquired — starting bot…")
	} else {
		log.Println("[singleton] LOCK_DISABLE=1 — запускаемся БЕЗ лока (не запускай второй экземпляр!)")
	}

	dropPending := os.Getenv("DROP_PENDING_UPDATES") == "1"
	if _, err := bot.Request(tgbotapi.DeleteWebhookConfig{DropPendingUpdates: dropPending}); err != nil {
		log.Printf("[boot] deleteWebhook warn: %v", err)
	}

	svcs, err := initServices(ctx, cfg)
	if err != nil {
		log.Fatalf("[boot] services: %v", err)
	}

	stateMgr, err := state.NewMongoManager(ctx, cfg.MongoURI, "coworking_bot", "user_states")
	if err != nil {
		log.Fatalf("[boot] state: %v", err)
	}

	reg := botengine.NewRegistry(stateMgr)
	flows.RegisterAll(reg)

	dispatcher := botengine.NewDispatcher(bot, cfg, svcs, reg)
	dispatcher.AttachPrinter(func(lang string) *message.Printer { return locales.Printer(lang) })

	if cfg.OrdersChatId != 0 {
		if chat, err := bot.GetChat(tgbotapi.ChatInfoConfig{
			tgbotapi.ChatConfig{ChatID: cfg.OrdersChatId},
		}); err != nil {
			log.Printf("[boot] OrdersChatId GETCHAT FAIL: %v", err)
		} else {
			log.Printf("[boot] OrdersChatId ok: type=%s title=%q id=%d", chat.Type, chat.Title, chat.ID)
		}

		if member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
			tgbotapi.ChatConfigWithUser{ChatID: cfg.OrdersChatId, UserID: bot.Self.ID},
		}); err != nil {
			log.Printf("[boot] OrdersChatId GETCHATMEMBER FAIL: %v", err)
		} else {
			log.Printf("[boot] Bot membership in OrdersChatId: status=%s", member.Status)
		}

		if os.Getenv("PING_ORDERS_ON_START") == "1" {
			ping := tgbotapi.NewMessage(cfg.OrdersChatId, "🤖 Bot online · orders will appear here")
			if _, err := bot.Send(ping); err != nil {
				log.Printf("[boot] OrdersChatId STARTUP PING FAIL: %v", err)
			} else {
				log.Printf("[boot] OrdersChatId STARTUP PING ok")
			}
		}
	}

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		<-signals
		log.Println("[boot] shutdown signal received")
		cancel()
	}()

	go botengine.RunWeeklyEvents(ctx, dispatcher, reg, stateMgr, cfg)
	dispatcher.Run(ctx)
	log.Println("Bye 👋")
}

func initServices(ctx context.Context, cfg *config.Config) (types.Services, error) {
	googleClient, err := services.NewGoogleClient(ctx, cfg.GoogleCloudConfig, cfg.Scope)
	if err != nil {
		return types.Services{}, err
	}

	coworkersSheets, err := services.NewCoworkersSheets(
		ctx,
		googleClient,
		cfg.CoworkersSpreadsheetId,
		cfg.CoworkersReadRange,
	)
	if err != nil {
		return types.Services{}, err
	}

	guestsSheets, err := services.NewGuestSheets(
		ctx,
		googleClient,
		cfg.CoworkersSpreadsheetId,
		cfg.GuestsReadRange,
	)
	if err != nil {
		return types.Services{}, err
	}

	botLogs, err := services.NewBotLogsSheets(
		ctx,
		googleClient,
		cfg.CoworkersSpreadsheetId,
		cfg.BotLogsReadRange,
	)
	if err != nil {
		return types.Services{}, err
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}

	eventsService := services.NewEventsService(
		httpClient,
		"https://shark-app-wrcei.ondigitalocean.app/api/events",
	)

	subs := services.NewMemSubscriptions()

	// Либо сохранить, либо удалить
	// barService := services.NewHaysellBarService(
	// 	httpClient,
	// 	cfg.HaysellBaseURL,
	// 	cfg.HaysellAPIKey,
	// )

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return types.Services{}, err
	}

	db := mongoClient.Database("coworking_bot")

	registrationsRepo, err := services.NewCoworkingRegistrationMongo(db)
	if err != nil {
		return types.Services{}, err
	}

	return types.Services{
		CoworkersSheets:        coworkersSheets,
		Guests:                 guestsSheets,
		BotLogs:                botLogs,
		Events:                 eventsService,
		Subscriptions:          subs,
		CoworkingRegistrations: registrationsRepo,
	}, nil
}

// ====== Mongo lock helpers (локальный, без отдельного пакета) ======

type lockDoc struct {
	ID       string    `bson:"_id"`
	Owner    string    `bson:"owner"`
	ExpireAt time.Time `bson:"expireAt"`
	Created  time.Time `bson:"createdAt,omitempty"`
}

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

// mongoWaitAcquire ждёт захвата лока; возвращает release().
func mongoWaitAcquire(ctx context.Context, coll *mongo.Collection, key string, ttl time.Duration) (func() error, error) {
	owner := newOwnerID()
	for {
		ok, wait, err := mongoTryLock(ctx, coll, key, ttl, owner)
		if err == nil && ok {
			// Heartbeat — продлеваем TTL каждые ttl/2
			stop := make(chan struct{})
			go func() {
				t := time.NewTicker(ttl / 2)
				defer t.Stop()
				for {
					select {
					case <-t.C:
						_ = mongoRenew(context.Background(), coll, key, owner, ttl)
					case <-stop:
						return
					}
				}
			}()
			release := func() error {
				close(stop)
				_, _ = coll.DeleteOne(context.Background(), bson.M{"_id": key, "owner": owner})
				return nil
			}
			return release, nil
		}
		// бэкофф при ошибке или ожидании чужого TTL
		if err != nil && wait <= 0 {
			wait = 500 * time.Millisecond
		}
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// mongoTryLock — атомарная попытка забрать/продлить лок.
func mongoTryLock(ctx context.Context, coll *mongo.Collection, key string, ttl time.Duration, owner string) (ok bool, wait time.Duration, err error) {
	now := time.Now()
	upd := bson.M{
		"$setOnInsert": bson.M{"createdAt": now},
		"$set":         bson.M{"owner": owner, "expireAt": now.Add(ttl)},
	}
	filter := bson.M{
		"_id": key,
		"$or": []bson.M{
			{"expireAt": bson.M{"$lte": now}}, // истёкший лок
			{"owner": owner},                  // наш — продлеваем
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var out lockDoc
	err = coll.FindOneAndUpdate(ctx, filter, upd, opts).Decode(&out)
	if err == nil {
		return true, 0, nil
	}

	// Не взяли — оценим оставшийся TTL
	var cur lockDoc
	if getErr := coll.FindOne(ctx, bson.M{"_id": key}).Decode(&cur); getErr != nil {
		return false, 300 * time.Millisecond, nil
	}
	left := time.Until(cur.ExpireAt)
	if left < 200*time.Millisecond {
		left = 200 * time.Millisecond
	}
	return false, left, nil
}

func mongoRenew(ctx context.Context, coll *mongo.Collection, key, owner string, ttl time.Duration) error {
	now := time.Now()
	res := coll.FindOneAndUpdate(
		ctx,
		bson.M{"_id": key, "owner": owner},
		bson.M{"$set": bson.M{"expireAt": now.Add(ttl)}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var d lockDoc
	return res.Decode(&d)
}
