// internal/botengine/deps.go
package botengine

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/state"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"golang.org/x/text/message"
)

// Deps — зависимости, которые мы прокидываем в шаги флоу.
type Deps struct {
	Bot        *tgbotapi.BotAPI
	Cfg        *config.Config
	Svcs       types.Services
	Printer    func(lang string) *message.Printer
	LastUpdate tgbotapi.Update
	State      state.Manager  
}
