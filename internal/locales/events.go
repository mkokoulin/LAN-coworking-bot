package locales

func registerEvents() {
	set(LangEN, "events_intro", s(`
		We host a variety of events. Follow us on
		<a href='https://www.instagram.com/lan_yerevan/'>Instagram</a> and
		<a href='https://t.me/lan_yerevan'>Telegram</a> to stay updated.
		The list of upcoming events is below ⬇️
	`))
	set(LangRU, "events_intro", s(`
		У нас проходит множество мероприятий. Подписывайтесь на
		<a href='https://www.instagram.com/lan_yerevan/'>Instagram</a> и
		<a href='https://t.me/lan_yerevan'>Telegram</a>, чтобы быть в курсе 🎉.
		Список ближайших мероприятий ниже ⬇️
	`))

	set(LangEN, "event_item", "%s %s <a href='https://lettersandnumbers.am/events/%s'>registration</a>\n\n")
	set(LangRU, "event_item", "%s %s <a href='https://lettersandnumbers.am/events/%s'>регистрация</a>\n\n")
}
