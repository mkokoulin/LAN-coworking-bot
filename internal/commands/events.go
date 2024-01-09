package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)

type Event struct {
	Capacity string
	Date string
	Description string
	ExternalLink string
	Id string
	Img string
	Link string
	Name string
	ShowForm bool
	Type string
}

func stripHtmlRegex(s string, regex string) string {
    r := regexp.MustCompile(regex)
    return r.ReplaceAllString(s, "")
}

func Events(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	msg.ParseMode = "html"

	client := http.Client{} 
	res, err := client.Get("https://shark-app-wrcei.ondigitalocean.app/api/events") 
	if err != nil { 
		return err
	} 

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil { 
		return err
	} 

	events := []Event{}

	err = json.Unmarshal(body, &events)
	if err != nil { 
		return err
	} 

	msg.ParseMode = "html"

	const regex = `<.*?>`

	eventsList := ""
	
	for _, e := range events {
		eventsList += fmt.Sprintf("%s %s <a href='https://lettersandnumbers.am/events/%s'>регистрация</a> \n\n", e.Date, stripHtmlRegex(e.Description, regex), e.Id)
	}
	
	if args.Storage.Language == Languages[0].Lang {
		// msg.Text = "We have a large number of different events, we publish announcements of events on our social networks: <a href='https://www.instagram.com/lan_yerevan/'>Instagram</a> and <a href='https://t.me/lan_yerevan'>Telegram</a>. Subscribe to keep up to date with cool events. An up-to-date list of events and reservations is maintained via <a href='https://lettersandnumbers.am/events'>site</a>"
		msg.Text = "We have a large number of different events, we publish announcements of events on our social networks: <a href='https://www.instagram.com/lan_yerevan/'>Instagram</a> and <a href='https://t.me/lan_yerevan'>Telegram</a>. The list of events is available below ⬇️"
	} else if args.Storage.Language == Languages[1].Lang {
		// msg.Text = "У нас проходит большое количество разнообразных мероприятий, анонсы событий мы публикуем в наших социальных сетях: <a href='https://www.instagram.com/lan_yerevan/'>Instagram</a> и <a href='https://t.me/lan_yerevan'>Telegram</a>. Подписывайтесь, чтобы быть в курсе классных событий 🎉. Актуальный список мероприятий и бронирование ведется через <a href='https://lettersandnumbers.am/events'>сайт</a>"
		msg.Text = "У нас проходит большое количество разнообразных мероприятий, анонсы событий мы публикуем в наших социальных сетях: <a href='https://www.instagram.com/lan_yerevan/'>Instagram</a> и <a href='https://t.me/lan_yerevan'>Telegram</a>. Подписывайтесь, чтобы быть в курсе классных событий 🎉. Ниже можно ознакомиться со списком мероприятий ⬇️"
	}

	args.Storage.CurrentCommand = ""

	msg.ParseMode = "html"

	_, err = bot.Send(msg)

	msg.Text = eventsList

	_, err = bot.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
		
	return err
}