package ui

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CbCmd(label, cmd string) tgbotapi.InlineKeyboardButton {
	if !strings.HasPrefix(cmd, "/") {
		cmd = "/" + cmd
	}
	return Cb(label, cmd)
}

