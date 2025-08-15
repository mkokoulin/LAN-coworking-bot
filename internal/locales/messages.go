package locales

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	LangEN = language.English
	LangRU = language.Russian
)

func Init() {
	// 🌐 Язык
	message.SetString(LangEN, "language_prompt", "Choose the interface language 🌎")
	message.SetString(LangRU, "language_prompt", "Выберите язык интерфейса 🌎")

	message.SetString(LangEN, "language_selected", "Language set to %s ✅")
	message.SetString(LangRU, "language_selected", "Язык изменён на %s ✅")

	// 🚀 Start
	message.SetString(LangEN, "start_message", `
		<b>Letters & Numbers — what’s inside</b>
		
		• 💻 <b>Coworking</b>
		• ☕ <b>Coffee bar</b>
		• ✨ <b>Event venue</b>

		<i>Tip:</i> check <a>/about</a> for locations and house rules.

		<b>Choose a command</b>:
		• <a>/start</a> — restart
		• <a>/booking</a> — book your first visit 🎁✨
		• <a>/wifi</a> — get the Wi-Fi password
		• <a>/meetingroom</a> — book a meeting room
		• <a>/printout</a> — send docs to print
		• <a>/events</a> — events info
		• <a>/menu</a> — bar menu 🍷
		• <a>/about</a> — about & map
		• <a>/language</a> — change language
		• <a>/bar</a> — bar menu
	`)

	message.SetString(LangRU, "start_message", `
		<b>Letters & Numbers — что внутри</b>
		
		• 💻 <b>Coworking</b>
		• ☕ <b>Coffee bar</b>
		• ✨ <b>Event venue</b>

		<i>Подсказка:</i> загляните в <a>/about</a> — там адреса и правила.

		<b>Выберите команду</b>:
		• <a>/start</a> — перезапуск
		• <a>/booking</a> — забронировать первый визит 🎁✨
		• <a>/wifi</a> — получить пароль Wi-Fi
		• <a>/meetingroom</a> — переговорка
		• <a>/printout</a> — отправить на печать
		• <a>/events</a> — события
		• <a>/menu</a> — барное меню 🍷
		• <a>/about</a> — о пространстве и схема
		• <a>/language</a> — сменить язык
		• <a>/bar</a> — барное меню
	`)

	// Wi-Fi
	message.SetString(LangEN, "select_network", "Select the network options below: guest / coworking")
	message.SetString(LangRU, "select_network", "Выберите ниже варианты сети: гостевой / коворкинг")

	message.SetString(LangEN, "wifi_guest", "L&N_guest network password %s")
	message.SetString(LangRU, "wifi_guest", "сеть L&N_guest пароль %s")

	message.SetString(LangEN, "wifi_guest_name", "Guest")
	message.SetString(LangRU, "wifi_guest_name", "Гостевой")

	message.SetString(LangEN, "wifi_coworking_name", "Coworking")
	message.SetString(LangRU, "wifi_coworking_name", "Коворкинг")

	message.SetString(LangEN, "wifi_coworking", "L&N network password %s")
	message.SetString(LangRU, "wifi_coworking", "сеть L&N пароль %s")

	message.SetString(LangEN, "ask_confirmation", "Enter the number you received from the administrator")
	message.SetString(LangRU, "ask_confirmation", "Введите номер, полученный от администратора")

	message.SetString(LangEN, "wrong_secret", "The password is incorrect, check with the administrator")
	message.SetString(LangRU, "wrong_secret", "Пароль неверный, уточните у администратора")

	// Booking
	message.SetString(LangEN, "booking_text", `
		You can book your first visit to the coworking space for 2 hours for free or a coffee of your choice. Just contact us: <a href='https://t.me/lan_yerevan'>telegram</a>, phone — +37494601303.
		`)

	message.SetString(LangRU, "booking_text", `
		Вы можете забронировать своё первое посещение коворкинга на 2 часа бесплатно или получить кофе на ваш выбор. Просто свяжитесь с нами: <a href='https://t.me/lan_yerevan'>telegram</a>, телефон — +37494601303.
		`)

	// Meeting Room
	message.SetString(LangEN, "meeting_prompt", "Write the date and time interval for which you want to book a meeting room in the format yyyy-mm-dd hh:mm - hh:mm")
	message.SetString(LangRU, "meeting_prompt", "Напишите дату и интервал времени, на который вы хотите забронировать комнату для переговоров в формате yyyy-mm-dd hh:mm - hh:mm")

	message.SetString(LangEN, "meeting_empty", "The message cannot be empty")
	message.SetString(LangRU, "meeting_empty", "Сообщение не может быть пустым")

	message.SetString(LangEN, "meeting_confirm", "Our administrator will contact you soon 🧑‍💼")
	message.SetString(LangRU, "meeting_confirm", "В ближайшее время с вами свяжется наш администратор 🧑‍💼")

	message.SetString(LangEN, "meeting_invalid_date_format", "❌ Please enter the interval in format YYYY-MM-DD HH:MM - HH:MM")
	message.SetString(LangRU, "meeting_invalid_date_format", "❌ Пожалуйста, введите интервал в формате ГГГГ-ММ-ДД ЧЧ:ММ - ЧЧ:ММ")

	message.SetString(LangRU, "meeting_pick_date", "📅 Выберите дату (на неделю вперёд):")
	message.SetString(LangRU, "meeting_pick_start_time", "⏱ Выберите время начала на %s:")
	message.SetString(LangRU, "meeting_pick_end_time", "⏱ Выберите время окончания для %s, начало %s:")
	message.SetString(LangRU, "meeting_select_date_first", "Сначала выберите дату 📅")
	message.SetString(LangRU, "meeting_flow_broken", "Хмм, что-то пошло не так. Давайте начнём заново 👇")
	message.SetString(LangRU, "meeting_invalid_interval", "❌ Неверный интервал: %s")
	message.SetString(LangRU, "meeting_confirm_interval", "✅ Бронь отправлена на подтверждение: %s")

	// Printout
	message.SetString(LangEN, "printout_info", "Send the documents for printing to the account @lan_yerevan (administrator) and check with them the cost of the service")
	message.SetString(LangRU, "printout_info", "Отправьте документы для распечатки в аккаунт @lan_yerevan (администратору) и уточните у него стоимость услуги")

	// Events
	message.SetString(LangEN, "events_intro", `
		We host a variety of events. Follow us on <a href='https://www.instagram.com/lan_yerevan/'>Instagram</a> and <a href='https://t.me/lan_yerevan'>Telegram</a> to stay updated. The list of upcoming events is below ⬇️
		`)

	message.SetString(LangRU, "events_intro", `
		У нас проходит множество мероприятий. Подписывайтесь на <a href='https://www.instagram.com/lan_yerevan/'>Instagram</a> и <a href='https://t.me/lan_yerevan'>Telegram</a>, чтобы быть в курсе событий 🎉. Список ближайших мероприятий ниже ⬇️
		`)

	message.SetString(LangEN, "event_item", "%s %s <a href='https://lettersandnumbers.am/events/%s'>registration</a>\n\n")
	message.SetString(LangRU, "event_item", "%s %s <a href='https://lettersandnumbers.am/events/%s'>регистрация</a>\n\n")

	// About
	message.SetString(LangEN, "about_text", `
		We are directing the layout of the site so that it is easier for you to navigate. The Letters and Numbers space houses: a coworking space, a coffee shop and an event space. Our locations and the rules of behavior in them are marked here.

		🐈 Address: Yerevan<a href='https://yandex.ru/maps/-/CDecr088'>, 35 Tumanyan str.</a>

		— To use the premises and coworking services, you must select and pay the appropriate tariff, you can get acquainted with the tariffs on <a href='https://lettersandnumbers.am/'>the site.</a>

		— We offer the coffee shop hall and the outdoor part of the site to the visitors of the coffee shop.

		💻 There is a quiet and noisy area in the coworking.

		The main coworking hall and part of the outdoor terrace by the window are a quiet area from 10:00 to 19:00. At this time, conversations are not appropriate and it is necessary to use headphones to watch videos. If one of the coworking visitors breaks the silence, then contact the administrator. After 19:00 in the main coworking area, you can call and talk, while maintaining the working atmosphere of the space. You can take coffee, tea, and cookies with you to the coworking room.

		The coffee shop hall and the courtyard are noisy areas (except for the tables by the window on terrace No. 1). Meetings, calls, and meals can be held here. You can bring food with you and store it in the refrigerator (through a barista), order delivery and, of course, purchase it in our cafe. The priority locations of coworkers are marked on the diagram.

		🕜 Coworking hours: weekdays 10-22, weekends 10-18. The playground is open every day from 10 to 22.
		`)

	message.SetString(LangRU, "about_text", `
		🗺️ Направляем схему площадки, чтобы вам было легче сориентироваться. В пространстве Letters and Numbers размещаются: коворкинг, кофейня и площадка для мероприятий. Здесь отмечены наши локации и правила поведения в них.

		🐈 Адрес: г. Ереван<a href='https://yandex.ru/maps/-/CDecr088'>, ул. Туманяна 35Г.</a>

		— Чтобы воспользоваться помещениями и услугами коворкинга, необходимо выбрать и оплатить соответствующий тариф, ознакомиться с тарифами можно на <a href='https://lettersandnumbers.am/'>сайте.</a>

		— Посетителям кофейни мы предлагаем зал кофейни и уличную часть площадки.

		💻 В коворкинге есть тихая и шумная зона.

		🤫 Основной зал коворкинга и часть уличной террасы у окна являются тихой зоной с 10:00 и до 19:00. В это время неуместны разговоры и обязательно использование наушников для просмотра видео. Если кто-то из посетителей коворкинга нарушает тишину, то обратитесь к администратору. После 19:00 в основной зоне коворкинга можно созваниваться и разговаривать, сохраняя рабочую атмосферу пространства. В зал коворкинга можно брать с собой кофе, чай, печенье.

		☕ Зал кофейни и двор являются шумными зонами (кроме столиков у окна на террасе №1). Здесь можно проводить встречи, звонки, принимать пищу. Еду можно принести с собой и оставить на хранение в холодильнике (через бариста), заказать доставку и, конечно, приобрести в нашем кафе. Приоритетные места размещения коворкеров отмечены на схеме.

		🕜 Время работы коворкинга: будни 10-22, выходные 10-18. Площадка открыта каждый день с 10 до 22.
		`)

	// Unknown command
	message.SetString(LangEN, "unknown_command", "I do not know this command 😔 Use the /start command.")
	message.SetString(LangRU, "unknown_command", "Я не знаю этой команды 😔 Воспользуйтесь командой /start.")

	// locales
	message.SetString(LangEN, "menu_unavailable", "The menu is temporarily unavailable. Please try again later.")
	message.SetString(LangRU, "menu_unavailable", "Меню временно недоступно. Пожалуйста, попробуйте позже.")

	// locales
	message.SetString(LangEN, "meeting_request_admin", "Meeting room request: %s")
	message.SetString(LangRU, "meeting_request_admin", "Заявка на переговорку: %s")

	// bar

	// RU
	message.SetString(LangRU, "bar_welcome", "🍽 Выберите позиции из меню. Жмите +/− на карточках товара. Когда будете готовы — откройте корзину ниже.")
	message.SetString(LangRU, "bar_cart_hint", "Когда выберете — откройте корзину: 🧺")
	message.SetString(LangRU, "bar_added", "Добавлено в корзину")
	message.SetString(LangRU, "bar_removed", "Обновили количество")
	message.SetString(LangRU, "bar_cart_title", "Корзина")
	message.SetString(LangRU, "bar_cart_empty", "Пока пусто. Добавьте что-нибудь из меню 🙂")
	message.SetString(LangRU, "bar_cart_cleared", "Корзина очищена")
	message.SetString(LangRU, "bar_ask_name", "На чьё имя оформить заказ? Напишите имя сообщением.")
	message.SetString(LangRU, "bar_ask_name_hint", "Напишите имя одним сообщением, например: «Миша».")
	message.SetString(LangRU, "bar_buyer_is", "Заказ оформляется на: <b>%s</b>")
	message.SetString(LangRU, "bar_order_sent", "✅ Заказ отправлен! Мы сверим детали и подтвердим в чате.")
	message.SetString(LangRU, "bar_order_cancelled", "Отменили. Если передумаете — корзина рядом 😉")

	// EN (если нужно)
	message.SetString(LangEN, "bar_welcome", "🍽 Pick items from the menu. Use +/− on product cards. When ready — open the cart below.")
	message.SetString(LangEN, "bar_cart_hint", "When ready — open your cart: 🧺")
	message.SetString(LangEN, "bar_added", "Added to cart")
	message.SetString(LangEN, "bar_removed", "Updated quantity")
	message.SetString(LangEN, "bar_cart_title", "Cart")
	message.SetString(LangEN, "bar_cart_empty", "The cart is empty yet.")
	message.SetString(LangEN, "bar_cart_cleared", "Cart cleared")
	message.SetString(LangEN, "bar_ask_name", "Whose name should we put the order under? Send a message with the name.")
	message.SetString(LangEN, "bar_ask_name_hint", "Please send a single message with the name, e.g. “Alex”.")
	message.SetString(LangEN, "bar_buyer_is", "Order for: <b>%s</b>")
	message.SetString(LangEN, "bar_order_sent", "✅ Order sent! We’ll confirm details here.")
	message.SetString(LangEN, "bar_order_cancelled", "Cancelled. Open the cart anytime 😉")

	message.SetString(LangRU, "unknown_command", "❓ Неизвестная команда: %s\nДоступные команды: %s")
	message.SetString(LangEN, "unknown_command", "❓ Unknown command: %s\nAvailable commands: %s")


	message.SetString(LangRU, "kotolog_intro", `
		<b>КОТОЛОГ 🐱</b>
		Здесь живут котики, которым нужен дом.
		Наши инициативы: книжный своп, лекции, книжная полочка и отложенные напитки — всё в пользу котиков.
		Выберите раздел ниже:`)
	message.SetString(LangRU, "kotolog_btn_view", "🐾 Посмотреть котиков")
	message.SetString(LangRU, "kotolog_btn_help", "🙌 Как помочь котикам")
	message.SetString(LangRU, "kotolog_btn_back", "⬅️ Назад")
	message.SetString(LangRU, "kotolog_btn_more", "Дальше →")
	message.SetString(LangRU, "kotolog_list_title", "Котики, которые ищут дом")
	message.SetString(LangRU, "kotolog_no_more", "Больше котиков пока нет — загляните позже 😸")
	message.SetString(LangRU, "kotolog_not_found", "Кот не найден — возможно уже дома. Ура! 🐾")
	message.SetString(LangRU, "kotolog_link_article", "📖 Статья")
	message.SetString(LangRU, "kotolog_link_photo", "Фото")
	message.SetString(LangRU, "kotolog_city", "Город")
	message.SetString(LangRU, "kotolog_contacts", "Контакты волонтёров")
	message.SetString(LangRU, "kotolog_help_text", `
		<b>Как помочь котикам</b>
		1) Благотворительный книжный своп — приносите книги, донаты идут котикам.
		2) Лекции в поддержку котиков — вход по донату.
		3) Книжная полочка — берите книги за пожертвование.
		4) Отложенные напитки — оплачиваете напиток заранее, и поддерживаете хвостиков.`)


	message.SetString(LangEN, "kotolog_intro", `
		<b>KOTOLOG 🐱</b>
		Cats looking for a loving home.
		Our initiatives: book swap, talks, bookshelf and suspended drinks — all for cats.
		Pick a section below:`)
	message.SetString(LangEN, "kotolog_btn_view", "🐾 View cats")
	message.SetString(LangEN, "kotolog_btn_help", "🙌 How to help cats")
	message.SetString(LangEN, "kotolog_btn_back", "⬅️ Back")
	message.SetString(LangEN, "kotolog_btn_more", "Next →")
	message.SetString(LangEN, "kotolog_list_title", "Cats looking for a home")
	message.SetString(LangEN, "kotolog_no_more", "No more cats for now — check back soon 😸")
	message.SetString(LangEN, "kotolog_not_found", "Cat not found — maybe already at home! 🐾")
	message.SetString(LangEN, "kotolog_link_article", "📖 Article")
	message.SetString(LangEN, "kotolog_link_photo", "Photo")
	message.SetString(LangEN, "kotolog_city", "City")
	message.SetString(LangEN, "kotolog_contacts", "Volunteer contacts")
	message.SetString(LangEN, "kotolog_help_text", `
		<b>How to help</b>
		1) Charity book swap — bring books, donations help cats.
		2) Talks — pay what you wish, proceeds go to cats.
		3) Bookshelf — take a book for a donation.
		4) Suspended drinks — prepay a drink, support cats.`)


			// ===== Russian =====
	message.SetString(LangRU, "bar_welcome",           "Привет! Это бар. Сначала выберите позиции из меню 👇")
	message.SetString(LangRU, "bar_cart_empty",        "Корзина пуста.")
	message.SetString(LangRU, "bar_cart_title",        "Корзина")
	message.SetString(LangRU, "bar_added",             "Добавлено в корзину")
	message.SetString(LangRU, "bar_removed",           "Убрано из корзины")
	message.SetString(LangRU, "bar_cart_cleared",      "Корзина очищена.")

	message.SetString(LangRU, "bar_ask_name",          "Как вас зовут? Напишите в одном сообщении.")
	message.SetString(LangRU, "bar_ask_name_hint",     "Пожалуйста, напишите имя текстом.")
	message.SetString(LangRU, "bar_ask_serve",         "Как подать заказ?")
	message.SetString(LangRU, "bar_ask_serve_hint",    "Нажмите кнопку ниже, чтобы выбрать способ подачи.")
	message.SetString(LangRU, "bar_ask_zone",          "Выберите зону для подачи:")

	message.SetString(LangRU, "bar_order_cancelled",   "Оформление заказа отменено.")
	message.SetString(LangRU, "bar_buyer_is",          "👤 Заказчик: <b>%s</b>")
	message.SetString(LangRU, "bar_order_sent",        "Заказ отправлен баристе. Мы напишем, когда будет готов!")

	// Промпты/тосты/тексты для комментария
	message.SetString(LangRU, "bar_notes_toast_prompt","Напишите комментарий одним сообщением")
	message.SetString(LangRU, "bar_notes_enter",       "Напишите комментарий для баристы (макс. 300 символов).")
	message.SetString(LangRU, "bar_notes_saved",       "📝 Комментарий сохранён.")
	message.SetString(LangRU, "bar_notes_deleted",     "Комментарий удалён")
	message.SetString(LangRU, "bar_notes_unchanged",   "Без изменений")

	// Подписи в подтверждении
	message.SetString(LangRU, "bar_contact_hint",      "☎️ Если что — пишите: %s")
	message.SetString(LangRU, "bar_comment_label",     "📝 Комментарий:")

	message.SetString(LangEN, "bar_welcome",           "Hi! This is the bar. First, pick items from the menu 👇")
	message.SetString(LangEN, "bar_cart_empty",        "Your cart is empty.")
	message.SetString(LangEN, "bar_cart_title",        "Cart")
	message.SetString(LangEN, "bar_added",             "Added to cart")
	message.SetString(LangEN, "bar_removed",           "Removed from cart")
	message.SetString(LangEN, "bar_cart_cleared",      "Cart cleared.")

	message.SetString(LangEN, "bar_ask_name",          "What’s your name? Please send it in one message.")
	message.SetString(LangEN, "bar_ask_name_hint",     "Please provide your name as text.")
	message.SetString(LangEN, "bar_ask_serve",         "How should we serve your order?")
	message.SetString(LangEN, "bar_ask_serve_hint",    "Use the buttons below to choose how to serve.")
	message.SetString(LangEN, "bar_ask_zone",          "Choose a zone for delivery:")

	message.SetString(LangEN, "bar_order_cancelled",   "Order checkout cancelled.")
	message.SetString(LangEN, "bar_buyer_is",          "👤 Customer: <b>%s</b>")
	message.SetString(LangEN, "bar_order_sent",        "Order sent to the barista. We’ll ping you when it’s ready!")

	// Prompts/toasts/notes texts
	message.SetString(LangEN, "bar_notes_toast_prompt","Type your comment in a single message")
	message.SetString(LangEN, "bar_notes_enter",       "Type a note for the barista (max 300 characters).")
	message.SetString(LangEN, "bar_notes_saved",       "📝 Note saved.")
	message.SetString(LangEN, "bar_notes_deleted",     "Note removed")
	message.SetString(LangEN, "bar_notes_unchanged",   "No changes")

	// Labels in confirmation
	message.SetString(LangEN, "bar_contact_hint",      "☎️ If needed — text: %s")
	message.SetString(LangEN, "bar_comment_label",     "📝 Comment:")
}
