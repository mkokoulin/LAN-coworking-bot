package flows

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowBar      types.Flow = "bar"
	stepBar      types.Step = "bar:show"
	stepHandle   types.Step = "bar:handle"
	stepAskName  types.Step = "bar:askname"
	stepConfirm  types.Step = "bar:confirm"
	stepDone     types.Step = "bar:done"

	stepBarHome  types.Step = "bar:home"
	stepAskServe types.Step = "bar:ask_serve"
	stepAskZone  types.Step = "bar:ask_zone"

	keyCart     = "order.cart"     // map[string]int
	keyBuyer    = "order.buyer"    // string
	keyLastMsg  = "order.lastmsg"  // int (message id для обновления, по желанию)
	keyCurrency = "order.currency" // string
)

type Item struct {
	ID       string
	Title    string
	PriceAMD int
	PhotoURL string
}

func getMenu() []Item {
	return []Item{
		{ID: "latte", Title: "Латте 300мл", PriceAMD: 1800, PhotoURL: "https://imgs.coffeemania.ru/822087/c848883d619a1ef/1200x800.jpg"},
		{ID: "americano", Title: "Американо 250мл", PriceAMD: 1200, PhotoURL: "https://imgs.coffeemania.ru/822087/c848883d619a1ef/1200x800.jpg"},
		{ID: "cookie", Title: "Печенье с шоколадом", PriceAMD: 900, PhotoURL: "https://imgs.coffeemania.ru/822087/c848883d619a1ef/1200x800.jpg"},
		{ID: "tea", Title: "Чай жасминовый", PriceAMD: 1000, PhotoURL: "https://imgs.coffeemania.ru/822087/c848883d619a1ef/1200x800.jpg"},
	}
}

func Register(reg *botengine.Registry) {
	reg.RegisterFlow(FlowBar, map[types.Step]botengine.StepHandler{
		stepBar:     prompt,
		stepHandle:  handle,
		stepAskName: askName,
		stepAskServe: askServe,
		stepAskZone: askZone,
		stepConfirm: confirm,
		stepDone:    done,
	})

	// Команда входа
	reg.RegisterCommand("bar", botengine.FlowEntry{Flow: FlowBar, Step: stepBar})

	// ВАЖНО: зарегистрировать все реально используемые префиксы
	reg.RegisterCallbackPrefix("bar:add:",      botengine.FlowEntry{Flow: FlowBar, Step: stepHandle})
	reg.RegisterCallbackPrefix("bar:rem:",      botengine.FlowEntry{Flow: FlowBar, Step: stepHandle})
	reg.RegisterCallbackPrefix("bar:peek:",     botengine.FlowEntry{Flow: FlowBar, Step: stepHandle})
	reg.RegisterCallbackPrefix("bar:noop",      botengine.FlowEntry{Flow: FlowBar, Step: stepHandle})

	reg.RegisterCallbackPrefix("bar:cart",      botengine.FlowEntry{Flow: FlowBar, Step: stepHandle})
	reg.RegisterCallbackPrefix("bar:rmitem:",   botengine.FlowEntry{Flow: FlowBar, Step: stepHandle})
	reg.RegisterCallbackPrefix("bar:clear",     botengine.FlowEntry{Flow: FlowBar, Step: stepHandle})
	reg.RegisterCallbackPrefix("bar:checkout",  botengine.FlowEntry{Flow: FlowBar, Step: stepHandle})

	reg.RegisterCallbackPrefix("bar:serve:",    botengine.FlowEntry{Flow: FlowBar, Step: stepHandle})
	reg.RegisterCallbackPrefix("bar:zone:",     botengine.FlowEntry{Flow: FlowBar, Step: stepHandle})

	reg.RegisterCallbackPrefix("bar:confirm",   botengine.FlowEntry{Flow: FlowBar, Step: stepConfirm})
	reg.RegisterCallbackPrefix("bar:notes",     botengine.FlowEntry{Flow: FlowBar, Step: stepConfirm})    // покрывает :clear и :cancel
	reg.RegisterCallbackPrefix("bar:cancel",    botengine.FlowEntry{Flow: FlowBar, Step: stepConfirm})

	// Кнопка «готово» из админского чата
	reg.RegisterCallbackPrefix("bar:done:",     botengine.FlowEntry{Flow: FlowBar, Step: stepHandle})
}
