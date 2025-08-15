// flows/kotolog_register.go
package flows

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
	"golang.org/x/text/message"
)

// --- Модель данных ---
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

// ⚠️ Это seed-данные. PhotoURL/ArticleURL можно заполнить позже (после загрузки фото и публикации статей).
var kotologCats = []Cat{
	{ID: "simba", Name: "Симба", Age: "1,5 месяца", Sex: "мальчик", Sterilized: false, Vaccinated: false,
		Character: "Очень деловой малыш; на передержке; обработан от паразитов; любит ласку, но не сидит на руках; очень подвижный.", City: "Ереван", Contacts: "@lan_yerevan"},
	{ID: "mikki", Name: "Микки (Микеланджела)", Age: "6 месяцев", Sex: "девочка", Sterilized: false, Vaccinated: false,
		Character: "Самая ласковая, мурчащая, доверчивая; любит вкусняшки; обработана от паразитов.", City: "Ереван", Contacts: "+37494601303"},
	{ID: "raffi", Name: "Раффи (Рафаэлло)", Age: "6 месяцев", Sex: "девочка", Sterilized: false, Vaccinated: false,
		Character: "Активная, любознательная, игривая; ориентирована на человека; обработана от паразитов.", City: "Ереван", Contacts: "+37494601303"},
	{ID: "roni", Name: "Рони", Age: "2 года", Sex: "девочка", Sterilized: true, Vaccinated: true,
		Character: "Социальная и разговорчивая миниатюрная кошечка-тигрица; стерилизована и привита.", City: "Ереван", Contacts: "@lan_yerevan"},
	{ID: "ronald", Name: "Рональд", Age: "2–3 года", Sex: "мальчик", Sterilized: true, Vaccinated: false,
		Character: "Кастрирован; задирает других кошек; очень боится людей; лучше единственным котом в доме.", City: "Ереван", Contacts: "@lan_yerevan"},
	{ID: "oreshka", Name: "Орешка", Age: "1 год", Sex: "девочка", Sterilized: true, Vaccinated: false,
		Character: "Пугливая, но любопытная; готова к поглаживаниям; на передержке; обработана от паразитов; покоряет любые высоты.", City: "Ереван", Contacts: "@lan_yerevan"},
	{ID: "pestrushka", Name: "Пеструшка", Age: "3 года", Sex: "девочка", Sterilized: true, Vaccinated: false,
		Character: "Осторожная, близко не подходит; пережила потерю котят; мама Орешки; лучше единственной кошкой; обработана от паразитов.", City: "Ереван", Contacts: "@lan_yerevan"},
	{ID: "musya-korovkina", Name: "Муся Коровкина", Age: "не меньше 2 лет", Sex: "девочка", Sterilized: true, Vaccinated: false,
		Character: "Нежная и пугливая; любит наблюдать за людьми; обработана; не идёт на руки, но дома раскроется; кошка-компаньон.", City: "Ереван", Contacts: "@lan_yerevan"},
	{ID: "zaika", Name: "Зайка", Age: "1–2 года", Sex: "девочка", Sterilized: true, Vaccinated: false,
		Character: "Ласковая; похоже, была домашней; не боится людей; не дерётся с другими кошками; красивая шубка.", City: "Ереван", Contacts: "@lan_yerevan"},
	{ID: "shrek", Name: "Шрек", Age: "около 1 года", Sex: "мальчик", Sterilized: false, Vaccinated: false,
		Character: "Компанейский и общительный рыжик; лучше без других животных; не кастрирован (в очереди); любит поглаживания; тихий.", City: "Ереван", Contacts: "@lan_yerevan"},
	{ID: "masyanya", Name: "Масяня", Age: "3 года", Sex: "девочка", Sterilized: true, Vaccinated: false,
		Character: "Строгая кошка-единоличница; с другими не уживается; сопровождает человека, но любит наблюдать со стороны.", City: "Ереван", Contacts: "@lan_yerevan"},
	{ID: "arbuzer", Name: "Арбузер", Age: "не меньше 2 лет", Sex: "мальчик", Sterilized: false, Vaccinated: false,
		Character: "Кот-мем: вид грустный, характер непредсказуем; пытается дружить, но пока не готов; дерётся с кошками; нужен на кастрацию.", City: "Ереван", Contacts: "@lan_yerevan"},
	{ID: "krikoslava", Name: "Крикослава", Age: "2 года", Sex: "девочка", Sterilized: true, Vaccinated: false,
		Character: "Пугливая, но нежная; стерилизована и обработана; любит компанию людей; самодостаточная охотница.", City: "Ереван", Contacts: "@lan_yerevan"},
	{ID: "krikoslav", Name: "Крикослав", Age: "2 года", Sex: "мальчик", Sterilized: true, Vaccinated: false,
		Character: "Громкий и наглый, но нежный; кастрирован и обработан; любит компанию людей; любит игрушки и мягкие кресла.", City: "Ереван", Contacts: "@lan_yerevan"},
	{ID: "sherkhan", Name: "Шерхан", Age: "2–3 года", Sex: "мальчик", Sterilized: true, Vaccinated: false,
		Character: "Пушистый царь гаражей; недоверчивый и осторожный; кастрирован; ладит с кошками; опекает Пеструшку и Орешку.", City: "Ереван", Contacts: "@lan_yerevan"},
}

func home(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)
	kb := ui.Inline(
		ui.Row(ui.Cb(p.Sprintf("kotolog_btn_view"), "kotolog:list:p1")),
		ui.Row(ui.Cb(p.Sprintf("kotolog_btn_help"), "kotolog:help")),
	)
	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("kotolog_intro"), kb)
	return KotologHome, nil
}

func list(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)
	page := 1
	if ev.Kind == botengine.EventCallback && strings.HasPrefix(ev.CallbackData, "kotolog:list:p") {
		if n, err := strconv.Atoi(strings.TrimPrefix(ev.CallbackData, "kotolog:list:p")); err == nil && n > 0 {
			page = n
		}
	}
	perPage := 5
	start := (page - 1) * perPage
	if start >= len(kotologCats) {
		// _ = ui.AnswerCallback(d.Bot, ev, p.Sprintf("kotolog_no_more"))
		return KotologList, nil
	}
	end := start + perPage
	if end > len(kotologCats) {
		end = len(kotologCats)
	}

var b strings.Builder
b.WriteString("<b>" + p.Sprintf("kotolog_list_title") + "</b>\n\n")
// for _, c := range kotologCats[start:end] {
// 	// b.WriteString(catCardHTML(p, c))
// 	b.WriteString("\n\n")
// }

	if end < len(kotologCats) {
		kb := ui.Inline(
			ui.Row(ui.Cb(p.Sprintf("kotolog_btn_back"), "kotolog:home")),
			ui.Row(ui.Cb(p.Sprintf("kotolog_btn_more"), fmt.Sprintf("kotolog:list:p%d", page+1))),
		)
		_ = ui.SendHTML(d.Bot, s.ChatID, b.String(), kb)
	} else {
		kb := ui.Inline(
			ui.Row(ui.Cb(p.Sprintf("kotolog_btn_back"), "kotolog:home")),
		)
		_ = ui.SendHTML(d.Bot, s.ChatID, b.String(), kb)
	}
	return KotologList, nil
}

func cat(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	// Детальная карточка открывается из списка по ссылке «Подробнее»
	id := ""
	if ev.Kind == botengine.EventCallback && strings.HasPrefix(ev.CallbackData, "kotolog:cat:") {
		id = strings.TrimPrefix(ev.CallbackData, "kotolog:cat:")
	}
	if id == "" {
		return KotologHome, nil
	}
	// var c Cat
	// found := false
	// for _, x := range kotologCats {
	// 	if x.ID == id {
	// 		c = x; found = true; break
	// 	}
	// }
	// p := d.Printer(s.Lang)
	// if !found {
	// 	_ = ui.AnswerCallback(d.Bot, ev, p.Sprintf("kotolog_not_found"))
	// 	return KotologList, nil
	// }
	// text := catFullHTML(p, c)
	// kb := ui.Inline(
	// 	ui.Row(ui.Cb(p.Sprintf("kotolog_btn_back"), "kotolog:list:p1")),
	// )
	// _ = ui.SendHTML(d.Bot, s.ChatID, text, kb)
	return KotologCat, nil
}

func help(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)
	kb := ui.Inline(
		ui.Row(ui.Cb("🐾 "+p.Sprintf("kotolog_btn_view"), "kotolog:list:p1")),
		ui.Row(ui.Cb(p.Sprintf("kotolog_btn_back"), "kotolog:home")),
	)
	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("kotolog_help_text"), kb)
	return KotologHelp, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return KotologHome, nil
}

// --- Вспомогательные функции рендера ---
func catCardHTML(p message.Printer, c Cat) string {
	flags := make([]string, 0, 2)
	if c.Sterilized { flags = append(flags, p.Sprintf("kotolog_flag_sterilized")) }
	if c.Vaccinated { flags = append(flags, p.Sprintf("kotolog_flag_vaccinated")) }
	meta := strings.Join(flags, ", ")
	if meta != "" { meta = " • " + meta }

	link := ""
	if c.ArticleURL != "" {
		link = fmt.Sprintf(" | <a href=\"%s\">%s</a>", c.ArticleURL, p.Sprintf("kotolog_link_article"))
	} else if isHTTP(c.PhotoURL) {
		link = fmt.Sprintf(" | <a href=\"%s\">%s</a>", c.PhotoURL, p.Sprintf("kotolog_link_photo"))
	}

	return fmt.Sprintf(`<b>%s</b> — %s, %s%s
%s%s`, c.Name, c.Sex, c.Age, meta, c.Character, link)
}

func catFullHTML(p message.Printer, c Cat) string {
	flags := make([]string, 0, 2)
	if c.Sterilized { flags = append(flags, p.Sprintf("kotolog_flag_sterilized")) }
	if c.Vaccinated { flags = append(flags, p.Sprintf("kotolog_flag_vaccinated")) }
	meta := strings.Join(flags, ", ")
	if meta != "" { meta = " • " + meta }

article := ""
if c.ArticleURL != "" {
	article = fmt.Sprintf("\n%s: %s", p.Sprintf("kotolog_link_article"), c.ArticleURL)
} else if isHTTP(c.PhotoURL) {
	article = fmt.Sprintf("\n%s: %s", p.Sprintf("kotolog_link_photo"), c.PhotoURL)
}

	return fmt.Sprintf(`
<b>%s</b> — %s, %s%s
%s
%s: <i>%s</i>%s
%s: %s
`, c.Name, c.Sex, c.Age, meta, c.Character,
		p.Sprintf("kotolog_city"), c.City, article,
		p.Sprintf("kotolog_contacts"), c.Contacts)
}

func isHTTP(u string) bool { return strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") }

// locales/kotolog.go — добавьте эти строки в ваш пакет locales (в Init или отдельной функции)
// package locales
// import (
// 	"golang.org/x/text/language"
// 	"golang.org/x/text/message"
// )
// func RegisterKotologStrings() {
// 	// RU
// 	message.SetString(LangRU, "kotolog_intro", `
// <b>КОТОЛОГ 🐱</b>
// Здесь живут котики, которым нужен дом.
// Наши инициативы: книжный своп, лекции, книжная полочка и отложенные напитки — всё в пользу котиков.
// Выберите раздел ниже:`)
// 	message.SetString(LangRU, "kotolog_btn_view", "🐾 Посмотреть котиков")
// 	message.SetString(LangRU, "kotolog_btn_help", "🙌 Как помочь котикам")
// 	message.SetString(LangRU, "kotolog_btn_back", "⬅️ Назад")
// 	message.SetString(LangRU, "kotolog_btn_more", "Дальше →")
// 	message.SetString(LangRU, "kotolog_list_title", "Котики, которые ищут дом")
// 	message.SetString(LangRU, "kotolog_no_more", "Больше котиков пока нет — загляните позже 😸")
// 	message.SetString(LangRU, "kotolog_not_found", "Кот не найден — возможно уже дома. Ура! 🐾")
// 	message.SetString(LangRU, "kotolog_link_article", "📖 Статья")
// 	message.SetString(LangRU, "kotolog_link_photo", "Фото")
// 	message.SetString(LangRU, "kotolog_city", "Город")
// 	message.SetString(LangRU, "kotolog_contacts", "Контакты волонтёров")
// 	message.SetString(LangRU, "kotolog_help_text", `
// <b>Как помочь котикам</b>
// 1) Благотворительный книжный своп — приносите книги, донаты идут котикам.
// 2) Лекции в поддержку котиков — вход по донату.
// 3) Книжная полочка — берите книги за пожертвование.
// 4) Отложенные напитки — оплачиваете напиток заранее, и поддерживаете хвостиков.`)
// 
// 	// EN
// 	message.SetString(LangEN, "kotolog_intro", `
// <b>KOTOLOG 🐱</b>
// Cats looking for a loving home.
// Our initiatives: book swap, talks, bookshelf and suspended drinks — all for cats.
// Pick a section below:`)
// 	message.SetString(LangEN, "kotolog_btn_view", "🐾 View cats")
// 	message.SetString(LangEN, "kotolog_btn_help", "🙌 How to help cats")
// 	message.SetString(LangEN, "kotolog_btn_back", "⬅️ Back")
// 	message.SetString(LangEN, "kotolog_btn_more", "Next →")
// 	message.SetString(LangEN, "kotolog_list_title", "Cats looking for a home")
// 	message.SetString(LangEN, "kotolog_no_more", "No more cats for now — check back soon 😸")
// 	message.SetString(LangEN, "kotolog_not_found", "Cat not found — maybe already at home! 🐾")
// 	message.SetString(LangEN, "kotolog_link_article", "📖 Article")
// 	message.SetString(LangEN, "kotolog_link_photo", "Photo")
// 	message.SetString(LangEN, "kotolog_city", "City")
// 	message.SetString(LangEN, "kotolog_contacts", "Volunteer contacts")
// 	message.SetString(LangEN, "kotolog_help_text", `
// <b>How to help</b>
// 1) Charity book swap — bring books, donations help cats.
// 2) Talks — pay what you wish, proceeds go to cats.
// 3) Bookshelf — take a book for a donation.
// 4) Suspended drinks — prepay a drink, support cats.`)
// }

// Вызовите locales.RegisterKotologStrings() из вашего locales.Init().
