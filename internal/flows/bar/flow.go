package flow

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	stepBar       types.Step = "bar:show"
	stepHandle    types.Step = "bar:handle"
	stepAskName   types.Step = "bar:askname"
	stepConfirm   types.Step = "bar:confirm"
	stepDone      types.Step = "bar:done"

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
    // Шаги flow
    reg.RegisterFlow("bar", map[types.Step]botengine.StepHandler{
        stepBar:      prompt,
        stepHandle:   handle,
        stepAskName:  askName,
        stepAskServe: askServe,
        stepAskZone:  askZone,
        stepConfirm:  confirm,
        stepDone:     done,
    })

    // Вход в flow по команде
    reg.RegisterCommand("bar", botengine.FlowEntry{Flow: "bar", Step: stepBar})

    // ===== Callback'и, которые должен обрабатывать stepHandle =====
    for _, p := range []string{
        "bar:add:",
        "bar:rem:",
        "bar:peek:",
        "bar:cart",
        "bar:rmitem:",
        "bar:clear",
        "bar:checkout",
        "bar:serve:",
        "bar:zone:",
        "bar:noop",

        // ВАЖНО: кнопка из админского чата
        "bar:done:",
    } {
        reg.RegisterCallbackPrefix(p, botengine.FlowEntry{Flow: "bar", Step: stepHandle})
    }

    // ===== Callback'и, когда пользователь на экране подтверждения (stepConfirm) =====
    for _, p := range []string{
        "bar:confirm",
        "bar:cancel",
        "bar:notes",
        "bar:notes:clear",
        "bar:notes:cancel",
    } {
        reg.RegisterCallbackPrefix(p, botengine.FlowEntry{Flow: "bar", Step: stepConfirm})
    }
}
