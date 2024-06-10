package commands

import (
	"context"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)

func Menu(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	msg.ParseMode = "html"

	var pdf string

	if args.Storage.Language == Languages[0].Lang {
		pdf = "menu_eng"
	} else if args.Storage.Language == Languages[1].Lang {
		pdf = "menu_rus"
	}
	
	pdfFile, err := os.Open(fmt.Sprintf("internal/assets/%s.pdf", pdf))
	if err != nil {
		panic(err)
	}

	reader := tgbotapi.FileReader{ Name: "menu.pdf", Reader: pdfFile }

	file := tgbotapi.NewDocument(update.Message.Chat.ID, reader)

	_, err = bot.Send(file)
	if err != nil {
		fmt.Println(err)
	}
		
	return err
}
