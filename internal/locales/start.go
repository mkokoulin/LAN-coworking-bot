package locales

func registerStart() {
	set(LangEN, "start_message", s(`
		<b>Letters & Numbers — what’s inside</b>

		• 💻 <b>Coworking</b>
		• ☕️ <b>LAN Bar</b>
		• ✨ <b>Event venue</b>

		<i>Tip:</i> check <b>/about</b> for locations and house rules.

		<b>Choose a command</b>

		<b>Work & bookings</b>
		• /coworking — about coworking
		• /booking — book your first visit 🎁✨
		• /meetingroom — book a meeting room

		<b>Tools</b>
		• /wifi — get the Wi-Fi password
		• /printout — send docs to print

		<b>Bar</b>
		• /menu — full bar menu
		` +
		// `• /bar — coffee bar (menu & orders). <i>Full menu:</i> <b>/menu</b> ` +🍷
		`
		<b>Info</b>
		• /events — events info
		• /about — about & map
		• /language — change language
		• /kotolog — 🐱 kotolog
		• /start — restart

		<b>Support us</b>
		• /donation — donate to the project
	`))

	set(LangRU, "start_message", s(`
		<b>Letters & Numbers — что внутри</b>

		• 💻 <b>Коворкинг</b>
		• ☕️ <b>LAN Bar</b>
		• ✨ <b>Площадка для событий</b>

		<i>Подсказка:</i> в <b>/about</b> — адреса и правила.

		<b>Команды</b>

		<b>Работа и брони</b>
		• /coworking — о коворкинге
		• /booking — первая бронь 🎁✨
		• /meetingroom — переговорка

		<b>Инструменты</b>
		• /wifi — пароль Wi-Fi
		• /printout — печать документов

		<b>Бар</b>
		• /menu — полное меню бара
		` +
		// `• /bar — бар (заказы и меню). <i>Полное меню:</i> <b>/menu</b> ` +🍷
		`
		<b>Инфо</b>
		• /events — события
		• /about — инфо и карта
		• /language — язык
		• /kotolog — 🐱 котолог
		• /start — перезапуск

		<b>Поддержите нас</b>
		• /donation — поддержать проект
	`))
}
