package botengine

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/state"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

type Dispatcher struct {
	bot           *tgbotapi.BotAPI
	cfg           *config.Config
	services      types.Services
	stateManager  state.Manager
	commandRouter *CommandRouter
	commandLogger  *state.CommandLogger
}

func (d *Dispatcher) AttachLogger(logger *state.CommandLogger) {
	d.commandLogger = logger
}

func NewDispatcher(bot *tgbotapi.BotAPI, cfg *config.Config, services types.Services, stateManager state.Manager) *Dispatcher {
	return &Dispatcher{
		bot:           bot,
		cfg:           cfg,
		services:      services,
		stateManager:  stateManager,
		commandRouter: NewCommandRouter(),
	}
}

func (d *Dispatcher) Run(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := d.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		state := d.stateManager.Get(chatID)

		if update.Message.IsCommand() {
			// state.Reset()
			state.CurrentCommand = update.Message.Command()

			if d.commandLogger != nil {
				_ = d.commandLogger.Log(chatID, state.CurrentCommand)
			}
		}

		err := d.commandRouter.Handle(ctx, update, d.bot, d.cfg, d.servicesToArgs(), state)
		if err != nil {
			log.Printf("[Dispatcher] handler error: %v", err)
		}
	}
}

func (d *Dispatcher) servicesToArgs() types.Services {
	return d.services
}
