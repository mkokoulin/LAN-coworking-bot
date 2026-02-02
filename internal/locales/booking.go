package locales

func registerBooking() {
	set(LangEN, "booking_text", s(`
		You can book your first coworking visit for free for 2 hours or enjoy a coffee of your choice.
		Just get in touch with us:
		<a href='https://t.me/lan_yerevan'>Telegram</a>, phone — +37494601303. 
		Please mention the promo code LAN-BOT.
	`))
	set(LangRU, "booking_text", s(`
		Вы можете забронировать своё первое посещение коворкинга на 2 часа бесплатно или получить кофе на ваш выбор. 
		Просто свяжитесь с нами:
		<a href='https://t.me/lan_yerevan'>telegram</a>, телефон — +37494601303. Сообщите промокод ЛАН-БОТ (LAN-BOT)
	`))
}