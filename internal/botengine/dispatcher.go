package botengine

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/text/message"

	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

type Dispatcher struct {
	bot          *tgbotapi.BotAPI
	cfg          *config.Config
	svcs         types.Services
	registry     *Registry
	printerFunc  func(lang string) *message.Printer
}

func NewDispatcher(bot *tgbotapi.BotAPI, cfg *config.Config, services types.Services, reg *Registry) *Dispatcher {
	return &Dispatcher{
		bot:      bot,
		cfg:      cfg,
		svcs:     services,
		registry: reg,
	}
}

func (d *Dispatcher) AttachPrinter(printer func(lang string) *message.Printer) { d.printerFunc = printer }

func (d *Dispatcher) Run(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{"message", "callback_query", "my_chat_member"}

	updates := d.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil && update.MyChatMember == nil {
			continue
		}

		chatID := ResolveChatID(update)
		if chatID == 0 {
			continue
		}

		sess := d.registry.Store.Get(chatID)

		ev := Classify(update)

		if ev.FromUserID != 0 {
			sess.UserID = ev.FromUserID
		}
		if sess.Lang == "" {
			sess.Lang = "ru"
		}

		deps := Deps{
			Bot:        d.bot,
			Cfg:        d.cfg,
			Svcs:       d.svcs,
			Printer:    d.printerFunc,
			LastUpdate: update,
			State:      d.registry.Store,
		}

		if err := RunFSM(ctx, ev, d.registry, deps, sess); err != nil {
			log.Printf("[dispatcher] fsm error: %v", err)
		}

		d.registry.Store.Set(sess.ChatID, sess)
	}
}
