package flows

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
	"golang.org/x/text/message"
)

// -------------------- Данные --------------------

type Cat struct {
	ID         string
	Name       string
	Age        string // "~8 месяцев", "2 года"
	Sex        string // "мальчик" | "девочка"
	Sterilized bool
	Vaccinated bool
	Character  string
	City       string
	PhotoURL   string // публичная ссылка на фото (опционально)
	Contacts   string // @username волонтёра или телефон
	ArticleURL string // ссылка на Telegra.ph (Instant View)
}

const kotologViewURL = "https://www.canva.com/design/DAGRxCeCUy0/_kGPED4IzghEra57q2IRZw/view"

// Seed-список
var kotologCats = []Cat{
	{ID: "simba", Name: "Симба", Age: "1,5 месяца", Sex: "мальчик", Sterilized: false, Vaccinated: false, Character: "Очень деловой малыш; на передержке; обработан от паразитов; любит ласку, но не сидит на руках; очень подвижный.", City: "Ереван", Contacts: "@lan_yerevan", PhotoURL: "internal/assets/simba.png"},
	{ID: "mikki", Name: "Микки (Микеланджела)", Age: "6 месяцев", Sex: "девочка", Sterilized: false, Vaccinated: false, Character: "Самая ласковая, мурчащая, доверчивая; любит вкусняшки; обработана от паразитов.", City: "Ереван", Contacts: "+37494601303", PhotoURL: "internal/assets/mikki.png"},
	{ID: "raffi", Name: "Раффи (Рафаэлло)", Age: "6 месяцев", Sex: "девочка", Sterilized: false, Vaccinated: false, Character: "Активная, любознательная, игривая; ориентирована на человека; обработана от паразитов.", City: "Ереван", Contacts: "+37494601303", PhotoURL: "internal/assets/raffi.png"},
	{ID: "roni", Name: "Рони", Age: "2 года", Sex: "девочка", Sterilized: true, Vaccinated: true, Character: "Социальная и разговорчивая миниатюрная кошечка-тигрица; стерилизована и привита.", City: "Ереван", Contacts: "@lan_yerevan", PhotoURL: "internal/assets/roni.png"},
	{ID: "ronald", Name: "Рональд", Age: "2–3 года", Sex: "мальчик", Sterilized: true, Vaccinated: false, Character: "Кастрирован; задирает других кошек; очень боится людей; лучше единственным котом в доме.", City: "Ереван", Contacts: "@lan_yerevan", PhotoURL: "internal/assets/ronald.png"},
	{ID: "oreshka", Name: "Орешка", Age: "1 год", Sex: "девочка", Sterilized: true, Vaccinated: false, Character: "Пугливая, но любопытная; готова к поглаживаниям; на передержке; обработана от паразитов; покоряет любые высоты.", City: "Ереван", Contacts: "@lan_yerevan", PhotoURL: "internal/assets/oreshka.png"},
	{ID: "pestrushka", Name: "Пеструшка", Age: "3 года", Sex: "девочка", Sterilized: true, Vaccinated: false, Character: "Осторожная, близко не подходит; пережила потерю котят; мама Орешки; лучше единственной кошкой; обработана от паразитов.", City: "Ереван", Contacts: "@lan_yerevan", PhotoURL: "internal/assets/pestrushka.png"},
	{ID: "musya-korovkina", Name: "Муся Коровкина", Age: "не меньше 2 лет", Sex: "девочка", Sterilized: true, Vaccinated: false, Character: "Нежная и пугливая; любит наблюдать за людьми; обработана; не идёт на руки, но дома раскроется; кошка-компаньон.", City: "Ереван", Contacts: "@lan_yerevan", PhotoURL: "internal/assets/musya-korovkina.png"},
	{ID: "zaika", Name: "Зайка", Age: "1–2 года", Sex: "девочка", Sterilized: true, Vaccinated: false, Character: "Ласковая; похоже, была домашней; не боится людей; не дерётся с другими кошками; красивая шубка.", City: "Ереван", Contacts: "@lan_yerevan", PhotoURL: "internal/assets/zaika.png"},
	{ID: "shrek", Name: "Шрек", Age: "около 1 года", Sex: "мальчик", Sterilized: false, Vaccinated: false, Character: "Компанейский и общительный рыжик; лучше без других животных; не кастрирован (в очереди); любит поглаживания; тихий.", City: "Ереван", Contacts: "@lan_yerevan", PhotoURL: "internal/assets/shrek.png"},
	{ID: "masyanya", Name: "Масяня", Age: "3 года", Sex: "девочка", Sterilized: true, Vaccinated: false, Character: "Строгая кошка-единоличница; с другими не уживается; сопровождает человека, но любит наблюдать со стороны.", City: "Ереван", Contacts: "@lan_yerevan", PhotoURL: "internal/assets/masyanya.png"},
	{ID: "arbuzer", Name: "Арбузер", Age: "не меньше 2 лет", Sex: "мальчик", Sterilized: false, Vaccinated: false, Character: "Кот-мем: вид грустный, характер непредсказуем; пытается дружить, но пока не готов; дерётся с кошками; нужен на кастрацию.", City: "Ереван", Contacts: "@lan_yerevan", PhotoURL: "internal/assets/arbuzer.png"},
	{ID: "krikoslava", Name: "Крикослава", Age: "2 года", Sex: "девочка", Sterilized: true, Vaccinated: false, Character: "Пугливая, но нежная; стерилизована и обработана; любит компанию людей; самодостаточная охотница.", City: "Ереван", Contacts: "@lan_yerevan", PhotoURL: "internal/assets/krikoslava.png"},
	{ID: "krikoslav", Name: "Крикослав", Age: "2 года", Sex: "мальчик", Sterilized: true, Vaccinated: false, Character: "Громкий и наглый, но нежный; кастрирован и обработан; любит компанию людей; любит игрушки и мягкие кресла.", City: "Ереван", Contacts: "@lan_yerevan", PhotoURL: "internal/assets/krikoslav.png"},
	{ID: "sherkhan", Name: "Шерхан", Age: "2–3 года", Sex: "мальчик", Sterilized: true, Vaccinated: false, Character: "Пушистый царь гаражей; недоверчивый и осторожный; кастрирован; ладит с кошками; опекает Пеструшку и Орешку.", City: "Ереван", Contacts: "@lan_yerevan", PhotoURL: "internal/assets/sherkhan.png"},
}

// -------------------- Handlers --------------------


func home(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	if ev.Kind == botengine.EventCallback {
		ackCallback(d, ev)
		switch {
		// ✅ удаляем обработку "kotolog:list:*" — больше не нужна
		case ev.CallbackData == "kotolog:help":
			s.Step = KotologHelp
			return botengine.InternalContinue, nil
		case ev.CallbackData == "kotolog:home":
			// Покажем интро ниже без смены шага
		default:
			s.Step = KotologHome
		}
	}

	kb := ui.Inline(
		// ✅ URL-кнопка вместо callback
		ui.Row(tgbotapi.NewInlineKeyboardButtonURL(p.Sprintf("kotolog_btn_view"), kotologViewURL)),
		ui.Row(ui.Cb(p.Sprintf("kotolog_btn_help"), "kotolog:help")),
	)

	sendOrEditHTML(d, s, ev, p.Sprintf("kotolog_intro"), kb)
	return KotologHome, nil
}


func list(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	if ev.Kind == botengine.EventCallback {
		ackCallback(d, ev)
		switch {
		case strings.HasPrefix(ev.CallbackData, "kotolog:cat:"):
			s.Step = KotologCat
			return botengine.InternalContinue, nil
		case ev.CallbackData == "kotolog:home":
			s.Step = KotologHome
			return botengine.InternalContinue, nil
		case ev.CallbackData == "kotolog:help":
			s.Step = KotologHelp
			return botengine.InternalContinue, nil
		}
	}

	page := 1
	if ev.Kind == botengine.EventCallback && strings.HasPrefix(ev.CallbackData, "kotolog:list:p") {
		if n, err := strconv.Atoi(strings.TrimPrefix(ev.CallbackData, "kotolog:list:p")); err == nil && n > 0 {
			page = n
		}
	}
	const perPage = 5
	total := len(kotologCats)
	maxPage := (total + perPage - 1) / perPage
	if page > maxPage {
		page = 1
	}
	start := (page - 1) * perPage
	end := start + perPage
	if end > total {
		end = total
	}

	var b strings.Builder
	b.WriteString("<b>" + p.Sprintf("kotolog_list_title") + "</b>\n\n")
	visible := kotologCats[start:end]
	for _, c := range visible {
		b.WriteString(catCardHTML(p, c))
		b.WriteString("\n\n")
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(visible)+1)
	for _, c := range visible {
		rows = append(rows, ui.Row(
			ui.Cb(p.Sprintf("kotolog_btn_more_about", c.Name), "kotolog:cat:"+c.ID),
		))
	}
	nav := []tgbotapi.InlineKeyboardButton{ui.Cb(p.Sprintf("kotolog_btn_back"), "kotolog:home")}
	if page < maxPage {
		nav = append(nav, ui.Cb(p.Sprintf("kotolog_btn_more"), fmt.Sprintf("kotolog:list:p%d", page+1)))
	}
	rows = append(rows, ui.Row(nav...))
	kb := ui.Inline(rows...)

	sendOrEditHTML(d, s, ev, b.String(), kb)
	return KotologList, nil
}

func fileExists(p string) bool {
    info, err := os.Stat(p)
    return err == nil && !info.IsDir()
}

func cat(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
    if ev.Kind != botengine.EventCallback || !strings.HasPrefix(ev.CallbackData, "kotolog:cat:") {
        s.Step = KotologHome
        return botengine.InternalContinue, nil
    }
    ackCallback(d, ev)

    id := strings.TrimPrefix(ev.CallbackData, "kotolog:cat:")
    c, ok := findCat(id)
    p := d.Printer(s.Lang)
    if !ok {
        kb := ui.Inline(ui.Row(
            ui.Cb(p.Sprintf("kotolog_btn_back_to_list"), "kotolog:list:p1"),
            ui.Cb(p.Sprintf("kotolog_btn_home"), "kotolog:home"),
        ))
        sendOrEditHTML(d, s, ev, "⚠️ "+p.Sprintf("kotolog_not_found"), kb)
        return KotologList, nil
    }

    text := catFullHTML(p, c)
    kb := ui.Inline(ui.Row(
        ui.Cb(p.Sprintf("kotolog_btn_back_to_list"), "kotolog:list:p1"),
        ui.Cb(p.Sprintf("kotolog_btn_home"), "kotolog:home"),
    ))
    sendOrEditHTML(d, s, ev, text, kb)

    caption := fmt.Sprintf("<b>%s</b> — %s\n<b>City:</b> Ереван\n<b>Volunteer contacts:</b> @lan_yerevan", c.Name, c.Age)

    var photo tgbotapi.PhotoConfig
    switch {
    case isHTTP(c.PhotoURL):
        // публичный URL — пусть Telegram сам скачивает
        photo = tgbotapi.NewPhoto(s.ChatID, tgbotapi.FileURL(c.PhotoURL))
    case fileExists(c.PhotoURL):
        // локальный файл — загружаем с сервера
        photo = tgbotapi.NewPhoto(s.ChatID, tgbotapi.FilePath(c.PhotoURL))
    default:
        log.Printf("[kotolog] photo not found or not http(s): %q", c.PhotoURL)
        return KotologCat, nil
    }

    photo.ParseMode = "HTML"
    photo.Caption = caption

    if _, err := d.Bot.Send(photo); err != nil {
        log.Printf("[kotolog] send photo fail: %v (source=%q)", err, c.PhotoURL)
    }

    return KotologCat, nil
}

const supportCardNumber = "0000 0000 0000 0000" // TODO: замените на реальный номер

func help(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	if ev.Kind == botengine.EventCallback && ev.CallbackData == "kotolog:copy_card" {
		ackCallback(d, ev)
		msg := tgbotapi.NewMessage(
			s.ChatID,
			p.Sprintf("kotolog_copy_msg", supportCardNumber), // "<code>%s</code>" внутри локали
		)
		msg.ParseMode = "HTML"
		msg.DisableWebPagePreview = true
		_, _ = d.Bot.Send(msg)
		// остаёмся на экране помощи
		return KotologHelp, nil
	}

	if ev.Kind == botengine.EventCallback {
		ackCallback(d, ev)
	}

	helpText := p.Sprintf("kotolog_help_text")
	donateNote := p.Sprintf("kotolog_donate_note", supportCardNumber)
	text := helpText + "\n\n" + donateNote

	kb := ui.Inline(
		ui.Row(tgbotapi.NewInlineKeyboardButtonURL(p.Sprintf("kotolog_btn_view"), kotologViewURL)),
		ui.Row(ui.Cb(p.Sprintf("kotolog_btn_back"), "kotolog:home")),
	)
	sendOrEditHTML(d, s, ev, text, kb)
	return KotologHelp, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return KotologHome, nil
}

func ackCallback(d botengine.Deps, ev botengine.Event) {
	if ev.Kind == botengine.EventCallback && ev.CallbackQueryID != "" {
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "") // пустой текст = просто закрыть спиннер
	}
}

func sendOrEditHTML(d botengine.Deps, s *types.Session, ev botengine.Event, html string, kb tgbotapi.InlineKeyboardMarkup) {
	if ev.Kind == botengine.EventCallback && ev.MessageID != 0 {
		edit := tgbotapi.NewEditMessageText(s.ChatID, ev.MessageID, html)
		edit.ParseMode = "HTML"
		edit.DisableWebPagePreview = true
		edit.ReplyMarkup = &kb
		_, _ = d.Bot.Request(edit)
		return
	}
	msg := tgbotapi.NewMessage(s.ChatID, html)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	msg.ReplyMarkup = kb
	_, _ = d.Bot.Send(msg)
}

// -------------------- Рендер --------------------

func catCardHTML(p *message.Printer, c Cat) string {
	flags := make([]string, 0, 2)
	if c.Sterilized {
		flags = append(flags, p.Sprintf("kotolog_flag_sterilized"))
	}
	if c.Vaccinated {
		flags = append(flags, p.Sprintf("kotolog_flag_vaccinated"))
	}
	meta := strings.Join(flags, ", ")
	if meta != "" {
		meta = " • " + meta
	}

	link := ""
	if c.ArticleURL != "" {
		link = fmt.Sprintf(" | <a href=\"%s\">%s</a>", c.ArticleURL, p.Sprintf("kotolog_link_article"))
	} else if isHTTP(c.PhotoURL) {
		link = fmt.Sprintf(" | <a href=\"%s\">%s</a>", c.PhotoURL, p.Sprintf("kotolog_link_photo"))
	}

	return fmt.Sprintf(
		`<b>%s</b> — %s, %s%s
%s%s`,
		c.Name, c.Sex, c.Age, meta, c.Character, link,
	)
}

func catFullHTML(p *message.Printer, c Cat) string {
	flags := make([]string, 0, 2)
	if c.Sterilized {
		flags = append(flags, p.Sprintf("kotolog_flag_sterilized"))
	}
	if c.Vaccinated {
		flags = append(flags, p.Sprintf("kotolog_flag_vaccinated"))
	}
	meta := strings.Join(flags, ", ")
	if meta != "" {
		meta = " • " + meta
	}

	article := ""
	if c.ArticleURL != "" {
		article = fmt.Sprintf("\n%s: %s", p.Sprintf("kotolog_link_article"), c.ArticleURL)
	} else if isHTTP(c.PhotoURL) {
		article = fmt.Sprintf("\n%s: %s", p.Sprintf("kotolog_link_photo"), c.PhotoURL)
	}

	return fmt.Sprintf(
		`<b>%s</b> — %s, %s%s
%s
%s: <i>%s</i>%s
%s: %s`,
		c.Name, c.Sex, c.Age, meta, c.Character,
		p.Sprintf("kotolog_city"), c.City, article,
		p.Sprintf("kotolog_contacts"), c.Contacts,
	)
}

func isHTTP(u string) bool {
	return strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")
}

func findCat(id string) (Cat, bool) {
	for _, c := range kotologCats {
		if c.ID == id {
			return c, true
		}
	}
	return Cat{}, false
}
