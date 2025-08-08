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

	message.SetString(LangEN, "language_selected", "Selected %s")
	message.SetString(LangRU, "language_selected", "Выбран %s")

	// 🚀 Start
	message.SetString(LangEN, "start_message", `
The Letters and Numbers space contains:
💻 coworking,
☕️ coffee shop and 
✨ event venue.

Be sure to check out the /about section — there you will find information about our locations and the rules of conduct in them.

Select the command to continue the dialog:

<b>Commands:</b>
/start – restart
/booking – book your first visit 🎁✨
/wifi – get a password from wifi
/meetingroom – book a meeting
/printout – send documents for printing
/events – information about events
/menu – bar menu 🍷
/about – information about the site and the scheme
/language – change interface language
`)

	message.SetString(LangRU, "start_message", `
В пространстве Letters and Numbers размещаются:
💻 коворкинг,
☕️ кофейня и
✨ площадка для мероприятий.

Обязательно ознакомьтесь с разделом /about — там вы найдете информацию о наших локациях и правилах поведения в них.

Выберите команду для продолжения диалога:

<b>Команды:</b>
/start – перезапуск
/booking – забронировать первое посещение 🎁✨
/wifi – получить пароль от вайфай
/meetingroom – забронировать переговорку
/printout – отправить документы на печать
/events – информация о мероприятиях
/menu – меню бара 🍷
/about – информация о площадке и схема
/language – смена языка интерфейса
`)

	// Wi-Fi
	message.SetString(LangEN, "select_network", "Select the network options below: guest / coworking")
	message.SetString(LangRU, "select_network", "Выберите ниже варианты сети: гостевой / коворкинг")

	message.SetString(LangEN, "wifi_guest", "L&N_guest network password %s")
	message.SetString(LangRU, "wifi_guest", "сеть L&N_guest пароль %s")

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
}
