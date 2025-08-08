package commands

import (
	"context"
	"errors"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/locales"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

func AboutCommand(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, services types.Services, state *types.ChatStorage) error {
	p := locales.Printer(state.Language)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, p.Sprintf("about_text"))
	msg.ParseMode = "HTML"

	photoBytes, err := os.ReadFile("internal/assets/Letters_and_Numbers_map.jpg")
	if err != nil {
		return errors.New("failed to read image: " + err.Error())
	}

	photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{
		Name:  "scheme.jpg",
		Bytes: photoBytes,
	})

	// Сначала пробуем отправить картинку
	if _, err := bot.Send(photo); err != nil {
		// Если не получилось — хотя бы текст
		_, err = bot.Send(msg)
		return err
	}

	// Затем отправляем текст
	_, err = bot.Send(msg)
	return err
}
