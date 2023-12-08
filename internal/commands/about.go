package commands

import (
	"context"
	"io/ioutil"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)

func About(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	msg.ParseMode = "html"
	if args.Storage.Language == Languages[0].Lang {
		msg.Text =
			"We are directing the layout of the site so that it is easier for you to navigate. The Letters and Numbers space houses: a coworking space, a coffee shop and an event space. Our locations and the rules of behavior in them are marked here.\n\n"+
			"🐈 Address: Yerevan<a href='https://yandex.ru/maps/-/CDecr088'>, 35 Tumanyan str.</a>\n\n"+
			"— To use the premises and coworking services, you must select and pay the appropriate tariff, you can get acquainted with the tariffs on <a href='https://lettersandnumbers.am /'>the site.</a>\n\n" +
			"— We offer the coffee shop hall and the outdoor part of the site to the visitors of the coffee shop.\n\n" +
			"💻 There is a quiet and noisy area in the coworking.\n\n" +
			"The main coworking hall and part of the outdoor terrace by the window are a quiet area from 10:00 to 19:00. At this time, conversations are not appropriate and it is necessary to use headphones to watch videos. If one of the coworking visitors breaks the silence, then contact the administrator. After 19:00 in the main coworking area, you can call and talk, while maintaining the working atmosphere of the space. You can take coffee, tea, and cookies with you to the coworking room.\n\n" +
			"The coffee shop hall and the courtyard are noisy areas (except for the tables by the window on terrace No. 1). Meetings, calls, and meals can be held here. You can bring food with you and store it in the refrigerator (through a barista), order delivery and, of course, purchase it in our cafe. The priority locations of coworkers are marked on the diagram.\n\n"+
			"🕜 Coworking hours: weekdays 10-22, weekends 10-18. The playground is open every day from 10 to 22."
	} else if args.Storage.Language == Languages[1].Lang {
		msg.Text = 
			"🗺️ Направляем схему площадки, чтобы вам было легче сориентироваться. В пространстве Letters and Numbers размещаются: коворкинг, кофейня и площадка для мероприятий. Здесь отмечены наши локации и правила поведения в них.\n\n" +
			"🐈 Адрес: г. Ереван<a href='https://yandex.ru/maps/-/CDecr088'>, ул. Туманяна 35Г.</a>\n\n" +
			"— Чтобы воспользоваться помещениями и услугами коворкинга, необходимо выбрать и оплатить соответствующий тариф, ознакомиться с тарифами можно на <a href='https://lettersandnumbers.am/'>сайте.</a>\n\n" +
			"— Посетителям кофейни мы предлагаем зал кофейни и уличную часть площадки.\n\n" +
			"💻 В коворкинге есть тихая и шумная зона.\n\n" +
			"🤫 Основной зал коворкинга и часть уличной террасы у окна являются тихой зоной с 10:00 и до 19:00. В это время неуместны разговоры и обязательно использование наушников для просмотра видео. Если кто-то из посетителей коворкинга нарушает тишину, то обратитесь к администратору. После 19:00 в основной зоне коворкинга можно созваниваться и разговаривать, сохраняя рабочую атмосферу пространства. В зал коворкинга можно брать с собой кофе, чай, печенье.\n\n" +
			"☕ Зал кофейни и двор являются шумными зонами (кроме столиков у окна на террасе №1). Здесь можно проводить встречи, звонки, принимать пищу. Еду можно принести с собой и оставить на хранение в холодильнике (через бариста), заказать доставку и, конечно, приобрести в нашем кафе. Приоритетные места размещения коворкеров отмечены на схеме.\n\n" +
			"🕜 Время работы коворкинга: будни 10-22, выходные 10-18. Площадка открыта каждый день с 10 до 22."
	}


	photoBytes, err := ioutil.ReadFile("internal/assets/Letters_and_Numbers_map.jpg")
	if err != nil {
		panic(err)
	}
	photoFileBytes := tgbotapi.FileBytes{
		Name:  "scheme",
		Bytes: photoBytes,
	}

	_, err = bot.Send(tgbotapi.NewPhoto(update.Message.Chat.ID, photoFileBytes))
	if err != nil {
		_, err = bot.Send(msg)
	}

	_, err = bot.Send(msg)
		
	return err
}