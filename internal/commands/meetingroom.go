package commands

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/locales"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

func MeetingroomCommand(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, services types.Services, state *types.ChatStorage) error {
	p := locales.Printer(state.Language)
	chatID := update.Message.Chat.ID
	text := strings.TrimSpace(update.Message.Text)

	// Первый шаг: просим ввести дату и время
	if !state.IsBookingProcess {
		state.IsBookingProcess = true
		msg := tgbotapi.NewMessage(chatID, p.Sprintf("meeting_prompt"))
		_, err := bot.Send(msg)
		return err
	}

	// Второй шаг: проверяем сообщение
	if text == "" {
		msg := tgbotapi.NewMessage(chatID, p.Sprintf("meeting_empty"))
		_, _ = bot.Send(msg)
		return nil
	}

	// Уведомляем администратора
	adminMsg := fmt.Sprintf("Пользователь @%s просит забронировать переговорку - %s", update.Message.Chat.UserName, text)
	_, _ = bot.Send(tgbotapi.NewMessage(cfg.AdminChatId, adminMsg))

	// Подтверждение пользователю
	confirm := tgbotapi.NewMessage(chatID, p.Sprintf("meeting_confirm"))
	_, _ = bot.Send(confirm)

	state.IsBookingProcess = false
	state.CurrentCommand = ""

	return nil
}
