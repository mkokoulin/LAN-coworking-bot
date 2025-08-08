package commands

import (
	"context"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/locales"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

func WifiCommand(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, services types.Services, state *types.ChatStorage) error {
	p := locales.Printer(state.Language)
	text := strings.ToLower(update.Message.Text)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	send := func(text string, markup any) error {
		msg.Text = text
		msg.ReplyMarkup = markup
		_, err := bot.Send(msg)
		return err
	}

	removeKeyboard := tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}

	// Шаг 1: выбор типа сети
	if !state.IsWifiProcess {
		state.IsWifiProcess = true
		return send(
			p.Sprintf("select_network"),
			tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("guest"),
					tgbotapi.NewKeyboardButton("coworking"),
				),
			),
		)
	}

	// Шаг 2: гостевая сеть
	if !state.IsAwaitingConfirmation && (text == "guest" || text == "гостевой") {
		state.IsWifiProcess = false
		state.CurrentCommand = ""
		return send(p.Sprintf("wifi_guest", cfg.GuestWifiPassword), removeKeyboard)
	}

	// Шаг 3: коворкинг — авторизован
	if !state.IsAwaitingConfirmation && (text == "coworking" || text == "коворкинг") {
		if state.IsAuthorized {
			state.IsWifiProcess = false
			state.CurrentCommand = ""
			return send(p.Sprintf("wifi_coworking", cfg.CoworkingWifiPassword), removeKeyboard)
		}

		state.IsAwaitingConfirmation = true
		return send(p.Sprintf("ask_confirmation"), nil)
	}

	// Шаг 4: проверка кода от администратора
	if state.IsAwaitingConfirmation {
		secrets, err := services.CoworkersSheets.GetUnusedSecrets(ctx)
		if err != nil {
			log.Fatalf("fatal error %v", err)
		}

		for _, s := range secrets {
			if text == s {
				state.IsWifiProcess = false
				state.IsAwaitingConfirmation = false
				state.IsAuthorized = true
				state.CurrentCommand = ""
				return send(p.Sprintf("wifi_coworking", cfg.CoworkingWifiPassword), removeKeyboard)
			}
		}

		// Неверный пароль
		return send(p.Sprintf("wrong_secret"), nil)
	}

	return nil
}
