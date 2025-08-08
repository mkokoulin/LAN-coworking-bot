package commands

import (
	"context"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

func MenuCommand(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, services types.Services, state *types.ChatStorage) error {
	var fileName string

	switch state.Language {
	case "en":
		fileName = "menu_eng"
	case "ru":
		fileName = "menu_rus"
	default:
		fileName = "menu_eng"
	}

	path := fmt.Sprintf("internal/assets/%s.pdf", fileName)
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open menu file: %w", err)
	}
	defer file.Close()

	doc := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FileReader{
		Name:   "menu.pdf",
		Reader: file,
	})

	_, err = bot.Send(doc)
	return err
}
