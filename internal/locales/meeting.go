// internal/locales/meeting.go
package locales

func registerMeeting() {
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

	// UPDATED: now takes 2 placeholders: interval + contact
	set(LangEN, "meeting_request_admin", "Meeting room request: %s\nContact: %s")
	set(LangRU, "meeting_request_admin", "Заявка на переговорку: %s\nКонтакт: %s")

	// NEW: ask for contact when username is missing
	set(LangEN, "meeting_need_contact", "✅ Booking drafted: %s\nI can’t see your Telegram username. Please send how to contact you (e.g. @handle or t.me link).")
	set(LangRU, "meeting_need_contact", "✅ Черновик брони: %s\nЯ не вижу ваш username в Telegram. Напишите, как с вами связаться (например @ник или ссылку t.me).")

	set(LangEN, "meeting_contact_too_long", "❌ Contact is too long. Please send a shorter handle/link.")
	set(LangRU, "meeting_contact_too_long", "❌ Слишком длинный контакт. Пришлите покороче (ник/ссылка).")
}
