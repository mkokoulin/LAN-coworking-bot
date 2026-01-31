package locales

import (
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Public language tags
var (
	LangEN = language.English
	LangRU = language.Russian
)

// s trims multiline literals nicely
func s(v string) string { return strings.TrimSpace(v) }

// set registers a key in snake_case and also a dot-alias (and vice versa) for backward compatibility
func set(lang language.Tag, key, val string) {
	v := strings.TrimSpace(val)
	message.SetString(lang, key, v)
	if strings.Contains(key, "_") {
		message.SetString(lang, strings.ReplaceAll(key, "_", "."), v)
	} else if strings.Contains(key, ".") {
		message.SetString(lang, strings.ReplaceAll(key, ".", "_"), v)
	}
}

// Init registers all localized strings
func Init() {
	registerLanguage()
	registerStart()
	registerWiFi()
	registerBooking()
	registerMeeting()
	registerPrintout()
	registerEvents()
	registerAbout()
	registerUnknownAndMisc()
	registerKotolog()
	registerDonation()
	registerBarGuest()
	registerBarAdmin()
	registerCoworking()
	registerMenu()
}
