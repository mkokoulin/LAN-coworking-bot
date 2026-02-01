package locales

func registerAbout() {
	set(LangEN, "about_text", s(`
		We are aligning the venue map to help you navigate. Letters & Numbers hosts:
		a coworking space, a coffee shop and an event area. Locations and rules are marked here.

		🐈 Address: Yerevan<a href='https://yandex.ru/maps/-/CDecr088'>, 35 Tumanyan St.</a>

		— To use the coworking and premises you need to choose and pay the appropriate plan.
		See plans on <a href='https://lettersandnumbers.am/'>the website</a>.

		— Café guests can use the café hall and the outdoor area.

		💻 Coworking has quiet and non-quiet zones.

		The main hall and part of the window-side terrace are a quiet zone 10:00–19:00:
		please avoid conversations and use headphones for video. Contact the admin if someone breaks silence.
		After 19:00 calls and talks are ok while keeping the working vibe. Drinks and cookies are welcome.

		The café hall and the yard are non-quiet zones (except window tables on terrace #1):
		meetings, calls and meals are fine. You can bring food (store via barista), order delivery, or buy in our café.

		🕜 Coworking hours: weekdays 10–22, weekends 10–16. The venue is open daily 10–22.
	`))

	set(LangRU, "about_text", s(`
		🗺️ Схема площадки, чтобы легче сориентироваться. В L&N есть коворкинг, кофейня и площадка для мероприятий.
		Здесь отмечены локации и правила поведения.

		🐈 Адрес: г. Ереван<a href='https://yandex.ru/maps/-/CDecr088'>, ул. Туманяна 35Г.</a>

		— Чтобы пользоваться помещениями и услугами коворкинга, выберите и оплатите соответствующий тариф —
		см. <a href='https://lettersandnumbers.am/'>сайт</a>.

		— Посетителям кофейни доступны зал кофейни и уличная часть площадки.

		💻 В коворкинге есть тихая и шумная зоны.

		🤫 Основной зал и часть террасы у окна — тихая зона 10:00–19:00: разговоры неуместны, видео только в наушниках.
		Если нарушают тишину — обратитесь к администратору. После 19:00 можно созваниваться, сохраняя рабочую атмосферу.
		Кофе/чай/печенье можно брать в зал.

		☕ Зал кофейни и двор — шумные зоны (кроме столиков у окна на террасе №1): встречи, звонки, приём пищи разрешены.
		Еду можно принести (хранение через бариста), заказать доставку или купить у нас.

		🕜 Коворкинг: будни 10–22, выходные 10–16. Площадка ежедневно 10–22.
	`))
}
