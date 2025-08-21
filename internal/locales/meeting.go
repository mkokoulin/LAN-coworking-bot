package locales

func registerMeeting() {
	set(LangEN, "meeting_prompt", "Write the date and time interval to book a meeting room in the format YYYY-MM-DD HH:MM - HH:MM")
	set(LangRU, "meeting_prompt", "Напишите дату и интервал времени для бронирования переговорки в формате ГГГГ-ММ-ДД ЧЧ:ММ - ЧЧ:ММ")

	set(LangEN, "meeting_empty", "The message cannot be empty")
	set(LangRU, "meeting_empty", "Сообщение не может быть пустым")

	set(LangEN, "meeting_confirm", "Our administrator will contact you soon 🧑💼")
	set(LangRU, "meeting_confirm", "В ближайшее время с вами свяжется наш администратор 🧑💼")

	set(LangEN, "meeting_invalid_date_format", "❌ Please enter the interval in format YYYY-MM-DD HH:MM - HH:MM")
	set(LangRU, "meeting_invalid_date_format", "❌ Пожалуйста, введите интервал в формате ГГГГ-ММ-ДД ЧЧ:ММ - ЧЧ:ММ")

	set(LangEN, "meeting_pick_date", "📅 Choose a date (up to one week ahead):")
	set(LangRU, "meeting_pick_date", "📅 Выберите дату (на неделю вперёд):")

	set(LangEN, "meeting_pick_start_time", "⏱ Choose a start time on %s:")
	set(LangRU, "meeting_pick_start_time", "⏱ Выберите время начала на %s:")

	set(LangEN, "meeting_pick_end_time", "⏱ Choose an end time on %s (start: %s):")
	set(LangRU, "meeting_pick_end_time", "⏱ Выберите время окончания для %s, начало %s:")

	set(LangEN, "meeting_select_date_first", "Please choose a date first 📅")
	set(LangRU, "meeting_select_date_first", "Сначала выберите дату 📅")

	set(LangEN, "meeting_flow_broken", "Hmm, something went wrong. Let’s start over 👇")
	set(LangRU, "meeting_flow_broken", "Хмм, что-то пошло не так. Давайте начнём заново 👇")

	set(LangEN, "meeting_invalid_interval", "❌ Invalid time interval: %s")
	set(LangRU, "meeting_invalid_interval", "❌ Неверный интервал: %s")

	set(LangEN, "meeting_confirm_interval", "✅ Booking sent for confirmation: %s")
	set(LangRU, "meeting_confirm_interval", "✅ Бронь отправлена на подтверждение: %s")

	set(LangEN, "meeting_request_admin", "Meeting room request: %s")
	set(LangRU, "meeting_request_admin", "Заявка на переговорку: %s")
}
