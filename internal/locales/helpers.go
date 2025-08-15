package locales

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// func Printer(lang string) *message.Printer {
//     return message.NewPrinter(language.Make(lang))
// }

func Printer(lang string) *message.Printer {
    if lang == "" {
        lang = language.English.String()
    }
    tag, err := language.Parse(lang)
    if err != nil {
        tag = language.English
    }
    return message.NewPrinter(tag)
}
