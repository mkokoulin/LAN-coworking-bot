package locales

import (
	"golang.org/x/text/message"
)

func Printer(lang string) *message.Printer {
	switch lang {
	case "ru":
		return message.NewPrinter(LangRU)
	default:
		return message.NewPrinter(LangEN)
	}
}
