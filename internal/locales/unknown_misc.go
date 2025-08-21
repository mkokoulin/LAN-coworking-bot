package locales

func registerUnknownAndMisc() {
	// Unknown command
	set(LangEN, "unknown_simple", "I don’t know this command 😔 Use /start.")
	set(LangRU, "unknown_simple", "Я не знаю этой команды 😔 Воспользуйтесь /start.")
	set(LangEN, "unknown_detailed", "❓ Unknown command: %s\nAvailable commands: %s")
	set(LangRU, "unknown_detailed", "❓ Неизвестная команда: %s\nДоступные команды: %s")
	set(LangEN, "unknown_command", "❓ Unknown command: %s\nAvailable commands: %s")
	set(LangRU, "unknown_command", "❓ Неизвестная команда: %s\nДоступные команды: %s")

	// Misc / menu
	set(LangEN, "menu_unavailable", "The menu is temporarily unavailable. Please try again later.")
	set(LangRU, "menu_unavailable", "Меню временно недоступно. Пожалуйста, попробуйте позже.")
}
