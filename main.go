package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/text/message"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/flows"
	"github.com/mkokoulin/LAN-coworking-bot/internal/locales"
	"github.com/mkokoulin/LAN-coworking-bot/internal/services"
	"github.com/mkokoulin/LAN-coworking-bot/internal/singleton"
	"github.com/mkokoulin/LAN-coworking-bot/internal/state"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

// быстрая проверка LP-доступности
func preflightLP(bot *tgbotapi.BotAPI) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1
	u.AllowedUpdates = []string{"message", "callback_query", "my_chat_member", "chat_member"}
	_, err := bot.GetUpdates(u)
	return err
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	locales.Init()

	// 1) Конфиг
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("[boot] load config: %v", err)
	}

	// 2) Telegram Bot
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("[boot] telegram: %v", err)
	}
	log.Printf("Bot started as @%s (debug=%v)", bot.Self.UserName, bot.Debug)

	// 3) LP-only: выключаем webhook (чтобы getUpdates работал)
	if _, err := bot.Request(tgbotapi.DeleteWebhookConfig{DropPendingUpdates: true}); err != nil {
		log.Printf("[boot] deleteWebhook warn: %v", err)
	}

	// 4) Preflight: если кто-то уже поллит токен / включён вебхук — выходим
	if err := preflightLP(bot); err != nil {
		es := strings.ToLower(err.Error())
		if strings.Contains(es, "conflict") ||
			strings.Contains(es, "terminated by other getupdates") ||
			strings.Contains(es, "webhook") {
			log.Fatalf("[boot] LP недоступен: %v", err)
		}
		log.Printf("[boot] preflight getUpdates warn: %v", err)
	}

	// 5) Монолок (персонифицированный под бота)
	lockID := "telegram_updates_lock:" + bot.Self.UserName
	if os.Getenv("LOCK_FORCE") == "1" {
		if err := singleton.ForceRelease(ctx, cfg.MongoURI, "coworking_bot", lockID); err != nil {
			log.Fatalf("[singleton] force release failed: %v", err)
		}
		log.Printf("[singleton] force-released %s", lockID)
	}
	if owner, exp, err := singleton.CurrentOwner(ctx, cfg.MongoURI, "coworking_bot", lockID); err == nil {
		log.Printf("[singleton] current lock owner: %s (expires %s)", owner, exp.Format(time.RFC3339))
	}
	lock := singleton.EnsureSingletonOrExit(ctx, cfg.MongoURI, "coworking_bot", lockID)
	defer lock.Release(context.Background(), lockID)

	// 6) Сервисы
	svcs, err := initServices(ctx, cfg)
	if err != nil {
		log.Fatalf("[boot] services: %v", err)
	}

	// 7) State manager
	stateMgr, err := state.NewMongo(cfg.MongoURI, "coworking_bot", "user_states")
	if err != nil {
		log.Fatalf("[boot] state: %v", err)
	}

	// 8) Registry + flows
	reg := botengine.NewRegistry()
	flows.RegisterAll(reg)

	// 9) Dispatcher
	dispatcher := botengine.NewDispatcher(bot, cfg, svcs, reg)
	dispatcher.AttachPrinter(func(lang string) *message.Printer { return locales.Printer(lang) })

	// 9.1) Проверим OrdersChatId и права бота — ПОЗИЦИОННЫЕ литералы для embedded полей
	if cfg.OrdersChatId != 0 {
		// GetChat: ChatInfoConfig содержит embedded ChatConfig — задаём позиционно
		if chat, err := bot.GetChat(tgbotapi.ChatInfoConfig{
			tgbotapi.ChatConfig{
				ChatID: cfg.OrdersChatId,
				// SuperGroupUsername: "", // альтернатива по username, если нужно
			},
		}); err != nil {
			log.Printf("[boot] OrdersChatId GETCHAT FAIL: %v", err)
		} else {
			log.Printf("[boot] OrdersChatId ok: type=%s title=%q id=%d",
				chat.Type, chat.Title, chat.ID)
		}

		// GetChatMember: GetChatMemberConfig содержит embedded ChatConfigWithUser — тоже позиционно
		if member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
			tgbotapi.ChatConfigWithUser{
				ChatID: cfg.OrdersChatId,
				UserID: bot.Self.ID,
			},
		}); err != nil {
			log.Printf("[boot] OrdersChatId GETCHATMEMBER FAIL: %v", err)
		} else {
			log.Printf("[boot] Bot membership in OrdersChatId: status=%s", member.Status)
		}

		// Опционально — пинг в заказной чат на старте
		if os.Getenv("PING_ORDERS_ON_START") == "1" {
			ping := tgbotapi.NewMessage(cfg.OrdersChatId, "🤖 Bot online · orders will appear here")
			if _, err := bot.Send(ping); err != nil {
				log.Printf("[boot] OrdersChatId STARTUP PING FAIL: %v", err)
			} else {
				log.Printf("[boot] OrdersChatId STARTUP PING ok")
			}
		}
	}

	// 10) Грейсфул-стоп
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		<-signals
		log.Println("[boot] shutdown signal received")
		cancel()
	}()

	// 11) Поехали
	go botengine.RunWeeklyEvents(ctx, dispatcher, reg, stateMgr, cfg)
	dispatcher.Run(ctx)
	log.Println("Bye 👋")
}

func initServices(ctx context.Context, cfg *config.Config) (types.Services, error) {
	googleClient, err := services.NewGoogleClient(ctx, cfg.GoogleCloudConfig, cfg.Scope)
	if err != nil {
		return types.Services{}, err
	}
	coworkersSheets, err := services.NewCoworkersSheets(ctx, googleClient, cfg.CoworkersSpreadsheetId, cfg.CoworkersReadRange)
	if err != nil {
		return types.Services{}, err
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}
	eventsService := services.NewEventsService(httpClient, "https://shark-app-wrcei.ondigitalocean.app/api/events")
	subs := services.NewMemSubscriptions()

	return types.Services{
		CoworkersSheets: coworkersSheets,
		Events:          eventsService,
		Subscriptions:   subs,
	}, nil
}
