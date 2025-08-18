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
		• /bar — coffee bar (menu & orders). <i>Full menu:</i> <b>/menu</b> 🍷

		<b>Info</b>
		• /events — events info
		• /about — about & map
		• /language — change language
		• /kotolog — 🐱 kotolog
		• /start — restart

		<b>Support us</b>
		• /donation — donate to the project
		`,
	)

	message.SetString(LangEN, "kotolog_btn_copy_card", "📋 Copy card number")
	message.SetString(LangEN, "kotolog_copy_msg",
		"Here is the number — long-press to copy:\n<code>%s</code>")

	message.SetString(LangEN, "kotolog_donate_note",
		"💛 <b>How to support</b>\n" +
		"You can send a donation to the card <code>%s</code>.\n" +
		"Please include <code>lan cats</code> in the payment title so we know it’s for the cats. Thank you 🐾")

	message.SetString(LangRU, "start_message", `
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
		• /bar — бар (заказы и меню). <i>Полное меню:</i> <b>/menu</b> 🍷

		<b>Инфо</b>
		• /events — события
		• /about — инфо и карта
		• /language — язык
		• /kotolog — 🐱 котолог
		• /start — перезапуск

		<b>Поддержите нас</b>
		• /donation — поддержать проект
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

	message.SetString(LangRU, "kotolog_btn_copy_card", "📋 Скопировать номер")
	message.SetString(LangRU, "kotolog_copy_msg",
		"Вот номер — удерживайте, чтобы скопировать:\n<code>%s</code>")

	message.SetString(LangRU, "kotolog_donate_note",
		"💛 <b>Как поддержать</b>\n" +
		"Можно оставить донат на карту <code>%s</code>.\n" +
		"Пожалуйста, укажите в названии платежа <code>lan cats</code> — так мы поймём, что это на котиков. Спасибо 🐾")

	message.SetString(LangRU, "kotolog_btn_back", "← Назад")
	message.SetString(LangRU, "kotolog_btn_home", "🏠 Домой") // или "На главную"


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
	message.SetString(LangRU, "kotolog_help_text", `<b>Как помочь котикам 🐾😺</b>

1) 😺📚 <b>Книжный своп</b> — приносите книги, которые уже не нужны. Пожертвования за обмен идут на корм и лечение котиков. 🐈
2) 😺🎤 <b>Лекции и беседы</b> — формат «сколько не жалко». Все средства направляем на поддержку котиков. 🐱
3) 😺📚 <b>Книжная полка</b> — берите книгу за донат любой суммы. Так мы поддерживаем полку и хвостатых. 🐈‍⬛
4) 😺☕️ <b>Отложенные напитки</b> — предоплатите чай/кофе для кого-то; деньги идут котикам. 🐾

Больше о наших проектах — @lan_yerevan. Будем рады пообщаться! 😺`)
	message.SetString(LangRU, "kotolog_btn_more_about", "Подробнее о %s")
	message.SetString(LangRU, "kotolog_btn_back_to_list", "← Назад к списку")

// Kotolog flags

message.SetString(LangRU, "kotolog_flag_sterilized", "стерилизован(а)")
message.SetString(LangRU, "kotolog_flag_vaccinated", "вакцинирован(а)")
	// внутри func Init()
	message.SetString(LangEN, "kotolog_btn_back", "← Back")
	message.SetString(LangEN, "kotolog_btn_home", "🏠 Home")
	
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
	message.SetString(LangEN, "kotolog_help_text", `<b>How to help 🐾😺</b>

1) 😺📚 <b>Charity book swap</b> — bring books you no longer need. Donations for every swap help pay for food and treatment for cats. 🐈
2) 😺🎤 <b>Talks</b> — pay what you wish. All proceeds go to support cats and their care. 🐱
3) 😺📚 <b>Bookshelf</b> — take a book for a donation of any size. Your support keeps the shelf alive. 🐈‍⬛
4) 😺☕️ <b>Suspended drinks</b> — prepay a drink for someone; the money goes to cats. 🐾

More about our projects — @lan_yerevan. We’ll be happy to chat! 😺`)
	message.SetString(LangEN, "kotolog_btn_more_about", "More about %s")

	message.SetString(LangEN, "kotolog_btn_back_to_list", "← Back to list")

	message.SetString(LangEN, "kotolog_flag_vaccinated", "vaccinated")
	message.SetString(LangEN, "kotolog_flag_sterilized", "sterilized")

	// ===== Bar — extra keys (RU) =====
	message.SetString(LangRU, "bar_menu_title", "🍹 <b>Меню</b>")
	message.SetString(LangRU, "bar_menu_hint", "Жмите +/− рядом с позицией. Корзина — кнопкой ниже.\n👁 Фото позиции — отдельным превью (само удалится через 8 сек).")
	message.SetString(LangRU, "bar_cart_total", "Итого: <b>%d AMD</b>")
	message.SetString(LangRU, "bar_item_price", "%s — %d AMD")
	message.SetString(LangRU, "bar_price_qty", "%d AMD • ×%d")
	message.SetString(LangRU, "bar_btn_photo", "👁 Фото")
	message.SetString(LangRU, "bar_btn_cart", "🧺 Корзина")
	message.SetString(LangRU, "bar_btn_clear", "🧹 Очистить")
	message.SetString(LangRU, "bar_btn_checkout", "✅ Оформить")
	message.SetString(LangRU, "bar_btn_back", "↩️ Назад")
	message.SetString(LangRU, "bar_btn_cancel", "↩️ Отменить")
	message.SetString(LangRU, "bar_btn_confirm", "✅ Подтвердить заказ")
	message.SetString(LangRU, "bar_btn_edit_note", "✏️ Изменить комментарий")
	message.SetString(LangRU, "bar_btn_delete_note", "🧽 Удалить комментарий")
	message.SetString(LangRU, "bar_btn_add_note", "📝 Комментарий для баристы")
	message.SetString(LangRU, "bar_in_cart_label", "В корзине: %d")

	message.SetString(LangRU, "bar_zone_coworking_name", "Коворкинг")
	message.SetString(LangRU, "bar_zone_cafe_name", "Кафе")
	message.SetString(LangRU, "bar_zone_street_name", "Улица")

	message.SetString(LangRU, "bar_serve_pickup_btn", "🧑‍🍳 Самовывоз с бара")
	message.SetString(LangRU, "bar_serve_tozone_btn", "🛎 Принести в зону")

	message.SetString(LangRU, "bar_serve_summary_label", "📍 Подача: <b>%s</b>")
	message.SetString(LangRU, "bar_serve_pickup_label", "Самовывоз с бара")
	message.SetString(LangRU, "bar_serve_tozone_label", "Принести в зону")
	message.SetString(LangRU, "bar_serve_tozone_with_label", "Принести в зону — %s")
	message.SetString(LangRU, "bar_not_specified", "не указано")

	// уведомления гостю при готовности
	message.SetString(LangRU, "bar_ready_pickup", "✅ Ваш заказ готов — можно забрать у бара.")
	message.SetString(LangRU, "bar_ready_tozone_generic", "✅ Ваш заказ готов — скоро принесём к вам.")
	message.SetString(LangRU, "bar_ready_tozone_zone", "✅ Ваш заказ готов — скоро принесём в %s.")
	message.SetString(LangRU, "bar_ready_generic", "✅ Ваш заказ готов.")

	// админская часть
	message.SetString(LangRU, "bar_admin_ready_btn", "✅ Готово к выдаче")
	message.SetString(LangRU, "bar_admin_issued_label", "✅ Выдано")
	message.SetString(LangRU, "bar_admin_ack_ok", "Принято")
	message.SetString(LangRU, "bar_admin_bad_button", "Некорректные данные кнопки")
	message.SetString(LangRU, "bar_admin_user_notified", "✅ Гостю отправлено уведомление о готовности.")
	message.SetString(LangRU, "bar_admin_notify_fail", "⚠️ Не удалось уведомить гостя (chat_id=%d): %v")

	message.SetString(LangRU, "bar_admin_new_order_title", "🧾 <b>Новый заказ</b>")
	message.SetString(LangRU, "bar_admin_order_no", "🔖 Номер: %s")
	message.SetString(LangRU, "bar_admin_name", "👤 Имя: %s")
	message.SetString(LangRU, "bar_line_item", "• %s × %d = %d AMD")
	message.SetString(LangRU, "bar_admin_serve_line", "📍 Подача: %s")
	message.SetString(LangRU, "bar_admin_questions_title", "❓ Уточнения:")
	message.SetString(LangRU, "bar_admin_q_delivery", "• Самовывоз или поднос до стола? Если стол — номер/описание?")
	message.SetString(LangRU, "bar_admin_q_disposables", "• Нужна ли одноразовая посуда/мешалка/сахар?")
	message.SetString(LangRU, "bar_admin_q_time", "• Время подачи (ASAP / ко времени)?")
	message.SetString(LangRU, "bar_admin_q_payment", "• Способ оплаты (нал/безнал)?")
	message.SetString(LangRU, "bar_admin_contact_line", "• Бариста для связи: %s")
	message.SetString(LangRU, "bar_admin_contact_meta", "• contact @%s, chat_id=%d")

	// подтверждение для гостя
	message.SetString(LangRU, "bar_order_number_label", "🔖 Номер заказа: <b>%s</b>")
	message.SetString(LangRU, "bar_order_customer_label", "👤 Заказчик: <b>%s</b>")
	message.SetString(LangRU, "bar_chat_label", "💬 Чат: %s")
	message.SetString(LangRU, "bar_open_chat", "открыть чат")

	// Промпты/тосты/тексты для комментария
	message.SetString(LangRU, "bar_notes_toast_prompt", "Напишите комментарий одним сообщением")
	message.SetString(LangRU, "bar_notes_enter", "Напишите комментарий для баристы (макс. 300 символов).")
	message.SetString(LangRU, "bar_notes_saved", "📝 Комментарий сохранён.")
	message.SetString(LangRU, "bar_notes_deleted", "Комментарий удалён")
	message.SetString(LangRU, "bar_notes_unchanged", "Без изменений")

	// Подписи в подтверждении
	message.SetString(LangRU, "bar_contact_hint", "☎️ Если что — пишите: %s")
	message.SetString(LangRU, "bar_comment_label", "📝 Комментарий:")

	message.SetString(LangRU, "bar_ask_zone", "В какую зону принести заказ?")
	message.SetString(LangRU, "bar_ask_serve", "Как подать заказ? Выберите способ обслуживания:")

	message.SetString(LangRU, "bar_overview", `
		<b>LAN Bar</b>

		• Кофе, чай, десерты, сезонные предложения  
		• Заказывайте у стойки или в этом чате

		<i>Полное меню:</i> <b>/menu</b> 🍷
		`)

		// ---------- RU ----------
	// Home
	message.SetString(language.Russian,
		"Letters & Numbers is an independent project. We exist thanks to your support ❤️",
		"Letters & Numbers — независимый проект. Мы живём за счёт вашей поддержки ❤️")
	message.SetString(language.Russian, "How you can support:", "Как можно поддержать:")
	message.SetString(language.Russian, "Attend our 🎟 events", "Приходить на наши 🎟 мероприятия")
	message.SetString(language.Russian, "Grab a coffee and desserts at the ☕ bar", "Заглядывать в ☕ бар за кофе и десертами")
	message.SetString(language.Russian, "Work from our 💻 coworking", "Работать у нас в 💻 коворкинге")
	message.SetString(language.Russian, "Or send a 💳 card donation (add note “lan cats”)",
		"Сделать 💳 донат на карту (в названии платежа указать «lan cats»)")
	message.SetString(language.Russian, "Choose an option:", "Выберите вариант:")

	// Buttons
	message.SetString(language.Russian, "💳 Card donation", "💳 Донат на карту")
	message.SetString(language.Russian, "🎟 Events", "🎟 Мероприятия")
	message.SetString(language.Russian, "☕ Bar", "☕ Бар")
	message.SetString(language.Russian, "💻 Coworking", "💻 Коворкинг")
	message.SetString(language.Russian, "⬅️ Home", "⬅️ На главную")
	message.SetString(language.Russian, "📋 Copy number", "📋 Скопировать номер")
	message.SetString(language.Russian, "⬅️ Back", "⬅️ Назад")
	message.SetString(language.Russian, "✅ Done", "✅ Готово")

	// Card screen
	message.SetString(language.Russian, "Card number:", "Номер карты:")
	message.SetString(language.Russian,
		"Important: add **lan cats** in payment note — this helps us understand the purpose.",
		"Важно: в названии платежа укажите **lan cats** — так мы быстрее поймём назначение.")
	message.SetString(language.Russian, "Thank you for your support! 🐱", "Спасибо за поддержку! 🐱")
	message.SetString(language.Russian, "Copy the card number from the message above:", "Скопируйте номер карты из сообщения выше:")
	message.SetString(language.Russian, "Thank you! /donation is always available.", "Спасибо! Раздел /donation всегда под рукой.")

	// ===== Bar — extra keys (EN) =====
	message.SetString(LangEN, "bar_menu_title", "🍹 <b>Menu</b>")
	message.SetString(LangEN, "bar_menu_hint", "Tap +/− near an item. Open the cart with the button below.\n👁 Item photo — preview (auto-deletes in 8s).")
	message.SetString(LangEN, "bar_cart_total", "Total: <b>%d AMD</b>")
	message.SetString(LangEN, "bar_item_price", "%s — %d AMD")
	message.SetString(LangEN, "bar_price_qty", "%d AMD • ×%d")
	message.SetString(LangEN, "bar_btn_photo", "👁 Photo")
	message.SetString(LangEN, "bar_btn_cart", "🧺 Cart")
	message.SetString(LangEN, "bar_btn_clear", "🧹 Clear")
	message.SetString(LangEN, "bar_btn_checkout", "✅ Checkout")
	message.SetString(LangEN, "bar_btn_back", "↩️ Back")
	message.SetString(LangEN, "bar_btn_cancel", "↩️ Cancel")
	message.SetString(LangEN, "bar_btn_confirm", "✅ Confirm order")
	message.SetString(LangEN, "bar_btn_edit_note", "✏️ Edit note")
	message.SetString(LangEN, "bar_btn_delete_note", "🧽 Delete note")
	message.SetString(LangEN, "bar_btn_add_note", "📝 Note for barista")
	message.SetString(LangEN, "bar_in_cart_label", "In cart: %d")

	message.SetString(LangEN, "bar_zone_coworking_name", "Coworking")
	message.SetString(LangEN, "bar_zone_cafe_name", "Cafe")
	message.SetString(LangEN, "bar_zone_street_name", "Street")

	message.SetString(LangEN, "bar_serve_pickup_btn", "🧑‍🍳 Pick up at bar")
	message.SetString(LangEN, "bar_serve_tozone_btn", "🛎 Bring to zone")

	message.SetString(LangEN, "bar_serve_summary_label", "📍 Serving: <b>%s</b>")
	message.SetString(LangEN, "bar_serve_pickup_label", "Pick up at the bar")
	message.SetString(LangEN, "bar_serve_tozone_label", "Bring to zone")
	message.SetString(LangEN, "bar_serve_tozone_with_label", "Bring to zone — %s")
	message.SetString(LangEN, "bar_not_specified", "not specified")

	message.SetString(LangEN, "bar_ready_pickup", "✅ Your order is ready — pick it up at the bar.")
	message.SetString(LangEN, "bar_ready_tozone_generic", "✅ Your order is ready — we’ll bring it to you shortly.")
	message.SetString(LangEN, "bar_ready_tozone_zone", "✅ Your order is ready — we’ll bring it to %s shortly.")
	message.SetString(LangEN, "bar_ready_generic", "✅ Your order is ready.")

	message.SetString(LangEN, "bar_admin_ready_btn", "✅ Ready to serve")
	message.SetString(LangEN, "bar_admin_issued_label", "✅ Served")
	message.SetString(LangEN, "bar_admin_ack_ok", "Accepted")
	message.SetString(LangEN, "bar_admin_bad_button", "Invalid button payload")
	message.SetString(LangEN, "bar_admin_user_notified", "✅ Guest has been notified.")
	message.SetString(LangEN, "bar_admin_notify_fail", "⚠️ Failed to notify guest (chat_id=%d): %v")

	message.SetString(LangEN, "bar_admin_new_order_title", "🧾 <b>New order</b>")
	message.SetString(LangEN, "bar_admin_order_no", "🔖 No: %s")
	message.SetString(LangEN, "bar_admin_name", "👤 Name: %s")
	message.SetString(LangEN, "bar_line_item", "• %s × %d = %d AMD")
	message.SetString(LangEN, "bar_admin_serve_line", "📍 Serving: %s")
	message.SetString(LangEN, "bar_admin_questions_title", "❓ Clarify:")
	message.SetString(LangEN, "bar_admin_q_delivery", "• Pickup or to table? If to table — number/description?")
	message.SetString(LangEN, "bar_admin_q_disposables", "• Disposable cup/stirrer/sugar?")
	message.SetString(LangEN, "bar_admin_q_time", "• Serving time (ASAP / specific time)?")
	message.SetString(LangEN, "bar_admin_q_payment", "• Payment (cash/card)?")
	message.SetString(LangEN, "bar_admin_contact_line", "• Barista contact: %s")
	message.SetString(LangEN, "bar_admin_contact_meta", "• contact @%s, chat_id=%d")

	message.SetString(LangEN, "bar_order_number_label", "🔖 Order number: <b>%s</b>")
	message.SetString(LangEN, "bar_order_customer_label", "👤 Customer: <b>%s</b>")
	message.SetString(LangEN, "bar_chat_label", "💬 Chat: %s")
	message.SetString(LangEN, "bar_open_chat", "open chat")

	// Prompts/toasts/notes texts
	message.SetString(LangEN, "bar_notes_toast_prompt", "Type your comment in a single message")
	message.SetString(LangEN, "bar_notes_enter", "Type a note for the barista (max 300 characters).")
	message.SetString(LangEN, "bar_notes_saved", "📝 Note saved.")
	message.SetString(LangEN, "bar_notes_deleted", "Note removed")
	message.SetString(LangEN, "bar_notes_unchanged", "No changes")

	// Labels in confirmation
	message.SetString(LangEN, "bar_contact_hint", "☎️ If needed — text: %s")
	message.SetString(LangEN, "bar_comment_label", "📝 Comment:")

	// locales/init.go (фрагмент)
	message.SetString(LangEN, "bar_ask_serve", "How should we serve your order? Choose a service type:")
	message.SetString(LangEN, "bar_ask_zone", "Which zone should we bring it to?")

	message.SetString(LangEN, "bar_overview", `
		<b>LAN Bar</b>

		• Coffee, tea, desserts, seasonal specials  
		• Order at the counter or via this chat

		<i>Explore the full menu:</i> <b>/menu</b> 🍷
		`)


// ---------- EN ----------
	// Home
	message.SetString(language.English,
		"Letters & Numbers is an independent project. We exist thanks to your support ❤️",
		"Letters & Numbers is an independent project. We exist thanks to your support ❤️")
	message.SetString(language.English, "How you can support:", "How you can support:")
	message.SetString(language.English, "Attend our 🎟 events", "Attend our 🎟 events")
	message.SetString(language.English, "Grab a coffee and desserts at the ☕ bar", "Grab a coffee and desserts at the ☕ bar")
	message.SetString(language.English, "Work from our 💻 coworking", "Work from our 💻 coworking")
	message.SetString(language.English, "Or send a 💳 card donation (add note “lan cats”)",
		"Or send a 💳 card donation (add note “lan cats”)")
	message.SetString(language.English, "Choose an option:", "Choose an option:")

	// Buttons
	message.SetString(language.English, "💳 Card donation", "💳 Card donation")
	message.SetString(language.English, "🎟 Events", "🎟 Events")
	message.SetString(language.English, "☕ Bar", "☕ Bar")
	message.SetString(language.English, "💻 Coworking", "💻 Coworking")
	message.SetString(language.English, "⬅️ Home", "⬅️ Home")
	message.SetString(language.English, "📋 Copy number", "📋 Copy number")
	message.SetString(language.English, "⬅️ Back", "⬅️ Back")
	message.SetString(language.English, "✅ Done", "✅ Done")

	// Card screen
	message.SetString(language.English, "Card number:", "Card number:")
	message.SetString(language.English,
		"Important: add **lan cats** in payment note — this helps us understand the purpose.",
		"Important: add **lan cats** in payment note — this helps us understand the purpose.")
	message.SetString(language.English, "Thank you for your support! 🐱", "Thank you for your support! 🐱")
	message.SetString(language.English, "Copy the card number from the message above:", "Copy the card number from the message above:")
	message.SetString(language.English, "Thank you! /donation is always available.", "Thank you! /donation is always available.")

// RU
	message.SetString(language.Russian, "coworking_intro",
		"💼 Letters & Numbers — коворкинг, бар и мероприятия в центре Еревана.\n"+
			"Ниже — актуальные цены и опции. Если хочется попробовать, начните с /booking — там есть приятный первый визит 😉")

	message.SetString(language.Russian, "coworking_prices",
		"💳 Тарифы коворкинга:\n"+
			"• 1 час — 1 300֏\n"+
			"• 4 часа — 3 000֏\n"+
			"• 1 день — 5 000֏\n"+
			"• 7 дней — 25 000֏\n"+
			"• 30 дней — 75 000֏\n"+
			"• LAN+ (60 дней) — 120 000֏")

	message.SetString(language.Russian, "coworking_meeting",
		"🧑‍💼 Переговорная (до 6 человек), цена за 1 час:\n"+
			"• «1+1» — 3 500֏\n"+
			"• «до 6 человек» — 5 000֏\n"+
			"• Для резидентов (доп. часы) — 2 000֏\n"+
			"*Скидка при бронировании от 5 часов. Бронирование кратно 30 минут.*")

	message.SetString(language.Russian, "coworking_options",
		"✨ Опции коворкинга:\n"+
			"• Часы работы: будни 10:00–22:00, выходные 10:00–16:00\n"+
			"• Безлимитный фильтр-кофе или чай\n"+
			"• Быстрый интернет ~200 Mbit/s\n"+
			"• Кондиционирование, естественное и искусственное освещение\n"+
			"• Удобное расположение в центре, собственный двор-сад\n"+
			"• Зоны для обеда; specialty coffee, café & bar\n"+
			"• Переговорная; сообщество и события\n"+
			"• Уличные террасы для работы\n"+
			"• Хранение багажа\n"+
			"• Любые способы оплаты (нал/карта/счёт)\n"+
			"• Спецусловия для команд и корпоративных клиентов")

	message.SetString(language.Russian, "coworking_btn_booking", "🎁 Первый визит")
	message.SetString(language.Russian, "coworking_btn_meetingroom", "📅 Переговорная")
	message.SetString(language.Russian, "coworking_btn_events", "🎟 События")
	message.SetString(language.Russian, "coworking_btn_bar", "☕ Бар")
	message.SetString(language.Russian, "coworking_btn_about", "ℹ️ О нас / карта")
	message.SetString(language.Russian, "coworking_btn_language", "🌐 Язык")

	// EN
	message.SetString(language.English, "coworking_intro",
		"💼 Letters & Numbers — coworking, bar and events in the heart of Yerevan.\n"+
			"Below are current prices and options. Want to try? Start with /booking — first visit comes with a nice bonus 😉")

	message.SetString(language.English, "coworking_prices",
		"💳 Coworking prices:\n"+
			"• 1 hour — 1,300֏\n"+
			"• 4 hours — 3,000֏\n"+
			"• 1 day — 5,000֏\n"+
			"• 7 days — 25,000֏\n"+
			"• 30 days — 75,000֏\n"+
			"• LAN+ (60 days) — 120,000֏")

	message.SetString(language.English, "coworking_meeting",
		"🧑‍💼 Meeting room (up to 6 people), price per hour:\n"+
			"• “1+1” — 3,500֏\n"+
			"• “up to 6 people” — 5,000֏\n"+
			"• Residents (additional hours) — 2,000֏\n"+
			"*Discount for bookings of 5 hours or more. Booking in 30-minute increments.*")

	message.SetString(language.English, "coworking_options",
		"✨ Coworking options:\n"+
			"• Working hours: weekdays 10:00–22:00, weekends 10:00–16:00\n"+
			"• Unlimited filter coffee or tea\n"+
			"• High-speed internet ~200 Mbit/s\n"+
			"• Air conditioning, natural & artificial lighting\n"+
			"• Convenient central location, private courtyard with garden\n"+
			"• Lunch zone; specialty coffee, café & bar\n"+
			"• Meeting room; community & events\n"+
			"• Outdoor terraces for work\n"+
			"• Luggage storage\n"+
			"• All payment types (cash/card/bank account)\n"+
			"• Special offers for teams and corporate clients")

	message.SetString(language.English, "coworking_btn_booking", "🎁 First visit")
	message.SetString(language.English, "coworking_btn_meetingroom", "📅 Meeting room")
	message.SetString(language.English, "coworking_btn_events", "🎟 Events")
	message.SetString(language.English, "coworking_btn_bar", "☕ Bar")
	message.SetString(language.English, "coworking_btn_about", "ℹ️ About & map")
	message.SetString(language.English, "coworking_btn_language", "🌐 Language")
}
