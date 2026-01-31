package flows

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

// --- ключи состояния ---
const (
	keyServe      = "bar:serve"
	keyZone       = "bar:zone"
	keyOrderID    = "bar:order_id"
	keyNotes      = "bar:notes"       // текст комментария
	keyAwaitNotes = "bar:await_notes" // флаг: ждём текст комментария
)

const (
	baristaContact = "@LAN_Barista" // для текста клиенту
	baristaMention = "@lan_barista" // для пинга в админском чате
)

// --- персист ключи ---
const (
	keyBarStateBlob  = "user:bar_state"
	keyCmdHistory    = "user:command_history"
	keyOrdersHistory = "bar:orders_history"
)

const (
    keyCurrentCategory = "bar:category" // выбранная категория
)

var (
	barCategories      []BarCategory
	barCategoriesByID  = map[string]BarCategory{}
	barCategoriesReady bool

	barItemsByID  map[string]Item    
)

type ProductItem struct {
    ID          int    `json:"id"`
    ImagePath   string `json:"imagePath"`
    ItemName    string `json:"itemName"`
    ShortName   string `json:"shortName"`
    Price       int    `json:"price"`
    Category    string `json:"category"`
    Description string `json:"description"`
    Balance     int    `json:"balance"`
}

type CategoryItem struct { 
	Sort int `json:"sort" bson:"sort"` 
	Deleted bool `json:"deleted" bson:"deleted"`
	Additional string `json:"additional" bson:"additional"` 
	Id int `json:"id" bson:"id"` 
	NameRu string `json:"nameRu" bson:"name_ru"` 
	MetaDescEn string `json:"metaDescEn" bson:"meta_desc_en"` 
	MetaTitleRu string `json:"metaTitleRu" bson:"meta_title_ru"` 
	ExtraCategories string `json:"extraCategories" bson:"extra_categories"`
	Recomended string `json:"recomended" bson:"recomended"` 
	ShortUrlEn string `json:"shortUrlEn" bson:"short_url_en"` 
	MetaKeyEn string `json:"metaKeyEn" bson:"meta_key_en"` 
	MetaDescRu string `json:"metaDescRu" bson:"meta_desc_ru"` 
	Logo string `json:"logo" bson:"logo"` 
	ShortUrlRu string `json:"shortUrlRu" bson:"short_url_ru"` 
	MetaDesc string `json:"metaDesc" bson:"meta_desc"` 
	MetaKeyRu string `json:"metaKeyRu" bson:"meta_key_ru"` 
	CatId int `json:"catId" bson:"cat_id"` 
	Name string `json:"name" bson:"name"`
	NameEn string `json:"nameEn" bson:"name_en"` 
	ParentId int `json:"parentId" bson:"parent_id"` 
	Level int `json:"level" bson:"level"` 
	ShortUrl string `json:"shortUrl" bson:"short_url"` 
	MetaKey string `json:"metaKey" bson:"meta_key"` 
	Image string `json:"image" bson:"image"` 
	MetaTitle string `json:"metaTitle" bson:"meta_title"` 
	MetaTitleEn string `json:"metaTitleEn" bson:"meta_title_en"` 
	Active int `json:"active" bson:"active"` 
	IsTop int `json:"isTop" bson:"is_top"`
}

// ----- модели персиста -----

type BarState struct {
	Cart      map[string]int `json:"cart,omitempty"`
	Buyer     string         `json:"buyer,omitempty"`
	Serve     string         `json:"serve,omitempty"`
	Zone      string         `json:"zone,omitempty"`
	Currency  string         `json:"currency,omitempty"`
	Notes     string         `json:"notes,omitempty"`
	OrderID   string         `json:"order_id,omitempty"`
	UpdatedAt int64          `json:"updated_at,omitempty"`
}

type CommandEntry struct {
	Cmd     string `json:"cmd"`
	AtUnix  int64  `json:"at"`
	Payload string `json:"payload,omitempty"`
}

type BarOrderItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Qty      int    `json:"qty"`
	PriceAMD int    `json:"price_amd"`
	SumAMD   int    `json:"sum_amd"`
}

type BarOrder struct {
	OrderID   string         `json:"order_id"`
	Buyer     string         `json:"buyer"`
	Items     []BarOrderItem `json:"items"`
	TotalAMD  int            `json:"total_amd"`
	Serve     string         `json:"serve"`
	Zone      string         `json:"zone"`
	Notes     string         `json:"notes,omitempty"`
	CreatedAt int64          `json:"created_at"`
}

type BarCategory struct {
	ID   string
	Name string
	Sort int
}

type haysellProduct struct {
    ID          int    `json:"id"`
    ImagePath   string `json:"imagePath"`
    ItemName    string `json:"itemName"`
    ShortName   string `json:"shortName"`
    Price       int    `json:"price"`
    Category    string `json:"category"`
    Description string `json:"description"`
    Balance     int    `json:"balance"`
}

var (
    // barMenuOnce   sync.Once
    barMenuErr    error
    barItemsAll   []Item
    barItemsByCat map[string][]Item
)

func allowedCategoryIDs() map[string]struct{} {
	m := make(map[string]struct{}, len(barCategories))
	for _, c := range barCategories {
		m[c.ID] = struct{}{}
	}
	return m
}

func buildBarMenuFromProducts(products []ProductItem) {
	barItemsAll = make([]Item, 0, len(products))
	barItemsByCat = make(map[string][]Item)
	barItemsByID = make(map[string]Item)

	// множество разрешённых категорий
	allowed := allowedCategoryIDs()
	useFilter := len(allowed) > 0 // если категорий ещё нет — не фильтруем, чтобы не сломать старое поведение

	for _, p := range products {
		// если категории уже загружены — берём только продукты из этих категорий
		if useFilter {
			if _, ok := allowed[p.Category]; !ok {
				continue
			}
		}

		it := Item{
			ID:         strconv.Itoa(p.ID),
			Title:      p.ItemName,
			PriceAMD:   p.Price,
			PhotoURL:   p.ImagePath,
			CategoryID: p.Category,
			ShortName: p.ShortName,
		}

		barItemsAll = append(barItemsAll, it)
		barItemsByCat[it.CategoryID] = append(barItemsByCat[it.CategoryID], it)
		barItemsByID[it.ID] = it
	}

	// сортируем общий список по названию
	sort.Slice(barItemsAll, func(i, j int) bool {
		return barItemsAll[i].ShortName < barItemsAll[j].ShortName
	})

	// сортируем внутри каждой категории
	for catID, items := range barItemsByCat {
		sort.Slice(items, func(i, j int) bool {
			return items[i].ShortName < items[j].ShortName
		})
		barItemsByCat[catID] = items
	}
}



// func ensureBarMenuLoaded(ctx context.Context, d botengine.Deps) {
// 	barMenuOnce.Do(func() {
// 		prods, err := fetchProducts(ctx, d) // ← []ProductItem
// 		if err != nil {
// 			barMenuErr = err
// 			log.Printf("[bar] fetchProducts failed: %v", err)
// 			return
// 		}
// 		buildBarMenuFromProducts(prods) // ← теперь принимает []ProductItem
// 	})
// }

// func ensureBarDataLoaded(ctx context.Context, d botengine.Deps) error {
// 	// Сначала продукты и меню
// 	ensureBarMenuLoaded(ctx, d)
// 	if barMenuErr != nil {
// 		// если не смогли подтянуть продукты — логируем, но даём шанс категориям
// 		log.Printf("[bar] failed to load bar menu: %v", barMenuErr)
// 	}

// 	// Потом категории (они уже могут видеть barItemsByCat)
// 	if err := ensureBarCategoriesLoaded(ctx, d); err != nil {
// 		return err
// 	}

// 	return barMenuErr
// }

func fetchCategoriesByParentId(ctx context.Context, d botengine.Deps) ([]CategoryItem, error) {
	baseURL := strings.TrimRight(d.Cfg.HaysellBaseURL, "/")
	if baseURL == "" {
		return nil, fmt.Errorf("HAYSELL_BASE_URL is empty in config")
	}

	u := baseURL + "/api/categories?parentId=1240"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var cats []CategoryItem
	if err := json.NewDecoder(resp.Body).Decode(&cats); err != nil {
		return nil, fmt.Errorf("decode categories: %w", err)
	}

	sort.Slice(cats, func(i, j int) bool {
		if cats[i].Sort == cats[j].Sort {
			nameI := cats[i].NameRu
			if nameI == "" {
				nameI = cats[i].NameEn
			}
			if nameI == "" {
				nameI = cats[i].Name
			}
			nameJ := cats[j].NameRu
			if nameJ == "" {
				nameJ = cats[j].NameEn
			}
			if nameJ == "" {
				nameJ = cats[j].Name
			}
			return nameI < nameJ
		}
		return cats[i].Sort < cats[j].Sort
	})

	return cats, nil
}


func fetchProducts(ctx context.Context, d botengine.Deps) ([]ProductItem, error) {
    baseURL := strings.TrimRight(d.Cfg.HaysellBaseURL, "/")
    if baseURL == "" {
        return nil, fmt.Errorf("HAYSELL_BASE_URL is empty in config")
    }

    u, err := url.Parse(baseURL + "/api/products")
    if err != nil {
        return nil, fmt.Errorf("parse base url: %w", err)
    }

    q := u.Query()
    // for _, id := range ids {
    //     q.Add("id", strconv.Itoa(id))
    // }
    u.RawQuery = q.Encode()

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
    if err != nil {
        return nil, fmt.Errorf("new request: %w", err)
    }

    client := &http.Client{Timeout: 5 * time.Second}

    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("do request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
    }

	var rawProducts map[string]haysellProduct
    if err := json.NewDecoder(resp.Body).Decode(&rawProducts); err != nil {
        return nil, fmt.Errorf("decode products: %w", err)
    }

    // res1 := make([]ProductItem, 0, len(rawProducts))
    // for _, p := range rawProducts {
    //     res1 = append(res1, ProductItem{
    //         ID:          p.ID,
    //         Category:  p.Category,
    //         ItemName:        p.ItemName,
    //         ShortName:   p.ShortName,
    //         Price:    p.Price,
    //         ImagePath:    p.ImagePath,
    //         Description: p.Description,
    //         Balance:     p.Balance,
    //     })
    // }

    // var raw map[string]ProductItem
    // if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
    //     return nil, fmt.Errorf("decode products: %w", err)
    // }

    res := make([]ProductItem, 0, len(rawProducts))
    for _, p := range rawProducts {
        res = append(res, ProductItem{
            ID:          p.ID,
            Category:  p.Category,
            ItemName:        p.ItemName,
            ShortName:   p.ShortName,
            Price:    p.Price,
            ImagePath:    p.ImagePath,
            Description: p.Description,
            Balance:     p.Balance,
        })
    }

    // sort.Slice(res, func(i, j int) bool {
    //     return res[i].ItemName < res[j].ItemName
    // })

    return res, nil
}


func loadBarCategoriesFromAPI(ctx context.Context, d botengine.Deps) error {
	cats, err := fetchCategoriesByParentId(ctx, d)
	if err != nil {
		return fmt.Errorf("load categories: %w", err)
	}

	var res []BarCategory
	byID := make(map[string]BarCategory)

	for _, c := range cats {
		if c.Deleted {
			continue
		}
		// if c.Additional != "coffochka" {
		// 	continue
		// }
		// if c.ParentId != 0 {
		// 	continue
		// }
		// если захочешь только кофейные:
		// if c.CatId != 2 { continue }

		idStr := strconv.Itoa(c.Id)

		name := c.NameRu
		if name == "" {
			name = c.NameEn
		}
		if name == "" {
			name = c.Name
		}
		if name == "" {
			name = idStr
		}

		bc := BarCategory{
			ID:   idStr,
			Name: name,
			Sort: c.Sort,
		}

		res = append(res, bc)
		byID[idStr] = bc
	}

	// 🔥 ВАЖНО: скрываем категории без товаров
	if barItemsByCat != nil {
		filtered := make([]BarCategory, 0, len(res))
		for _, c := range res {
			if items, ok := barItemsByCat[c.ID]; ok && len(items) > 0 {
				filtered = append(filtered, c)
			}
		}
		res = filtered
	}

	// сортировка по Sort, затем по Name
	sort.Slice(res, func(i, j int) bool {
		if res[i].Sort == res[j].Sort {
			return res[i].Name < res[j].Name
		}
		return res[i].Sort < res[j].Sort
	})

	barCategories = res
	barCategoriesByID = byID
	barCategoriesReady = true
	return nil
}


func allItems() []Item {
	return barItemsAll
}

func itemsByCategory(catID string) []Item {
	return barItemsByCat[catID]
}

func getBarCategories() []BarCategory {
    return barCategories
}


// func ensureBarCategoriesLoaded(ctx context.Context, d botengine.Deps) error {
// 	if barCategoriesReady {
// 		return nil
// 	}
// 	return loadBarCategoriesFromAPI(ctx, d)
// }

func categoryNameByID(id string) string {
    if c, ok := barCategoriesByID[id]; ok {
        return c.Name
    }
    return id
}

func renderCategoriesText(d botengine.Deps, s *types.Session) string {
	p := d.Printer(s.Lang)
	var b strings.Builder

	b.WriteString(p.Sprintf("bar_categories_title"))
	b.WriteString("\n\n")

	cats := getBarCategories()
	if len(cats) == 0 {
		b.WriteString(p.Sprintf("bar_categories_empty"))
		return b.String()
	}

	for _, c := range cats {
		b.WriteString("• ")
		b.WriteString(safeHTML(c.Name))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(p.Sprintf("bar_categories_hint"))
	return b.String()
}

func categoriesKeyboard(d botengine.Deps, s *types.Session) tgbotapi.InlineKeyboardMarkup {
	cats := getBarCategories()
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, c := range cats {
		rows = append(rows, ui.Row(
			ui.Cb(c.Name, "bar:cat:"+c.ID),
		))
	}

	return ui.Inline(rows...)
}

func renderMenuForItems(d botengine.Deps, s *types.Session, items []Item, header string) string {
	p := d.Printer(s.Lang)
	var b strings.Builder

	if header != "" {
		b.WriteString(header)
		b.WriteString("\n")
	} else {
		b.WriteString(p.Sprintf("bar_menu_title"))
		b.WriteString("\n")
	}

	if len(items) == 0 {
		b.WriteString(p.Sprintf("bar_category_empty"))
		return b.String()
	}

	for _, it := range items {
		qty := qtyInCart(s, it.ID)
		if qty > 0 {
			b.WriteString(fmt.Sprintf("• %s — %d AMD (×%d)\n", it.ShortName, it.PriceAMD, qty))
		} else {
			b.WriteString(fmt.Sprintf("• %s — %d AMD\n", it.ShortName, it.PriceAMD))
		}
	}

	b.WriteString(p.Sprintf("bar_menu_hint"))
	return b.String()
}

func menuKeyboardForItems(d botengine.Deps, s *types.Session, items []Item, withBackToCategories bool) tgbotapi.InlineKeyboardMarkup {
	p := d.Printer(s.Lang)
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, it := range items {
		qty := qtyInCart(s, it.ID)
		rows = append(rows, ui.Row(
			ui.Cb(it.ShortName, "bar:noop"),
		))

		rows = append(rows, ui.Row(
			ui.Cb(p.Sprintf("bar_price_qty", it.PriceAMD, qty), "bar:noop"),
			ui.Cb("−", "bar:rem:"+it.ID),
			ui.Cb("+", "bar:add:"+it.ID),
		))

		rows = append(rows, ui.Row(ui.Cb(" ", "bar:noop")))
	}

	if withBackToCategories {
		rows = append(rows,
			ui.Row(ui.Cb(p.Sprintf("bar_btn_back_to_categories"), "bar:cat:back")),
		)
	}

	rows = append(rows,
		ui.Row(ui.Cb(p.Sprintf("bar_btn_cart"), "bar:cart"), ui.Cb(p.Sprintf("bar_btn_clear"), "bar:clear")),
		ui.Row(ui.Cb(p.Sprintf("bar_btn_checkout"), "bar:checkout")),
	)

	return ui.Inline(rows...)
}

// удобный хелпер: меню для конкретной категории
func renderCategoryMenu(d botengine.Deps, s *types.Session, catID string) (string, tgbotapi.InlineKeyboardMarkup) {
	items := itemsByCategory(catID)
	header := fmt.Sprintf("📂 %s", safeHTML(categoryNameByID(catID)))
	txt := renderMenuForItems(d, s, items, header)
	kb := menuKeyboardForItems(d, s, items, true)
	return txt, kb
}

// ---------- корзина ----------
func addToCart(s *types.Session, id string, delta int) (newQty int, changed bool) {
	if findItem(id) == nil {
		return qtyInCart(s, id), false
	}
	cart := getCart(s)
	before := cart[id]
	after := before + delta
	if after <= 0 {
		delete(cart, id)
		after = 0
	} else {
		cart[id] = after
	}
	s.Data[keyCart] = cart
	return after, after != before
}

func removeItem(s *types.Session, id string) { c := getCart(s); delete(c, id); s.Data[keyCart] = c }
func clearCart(s *types.Session)              { s.Data[keyCart] = map[string]int{} }
func isCartEmpty(s *types.Session) bool       { return len(getCart(s)) == 0 }

func getCart(s *types.Session) map[string]int {
	if s == nil { return map[string]int{} }
	if s.Data == nil { s.Data = map[string]interface{}{} }
	raw, ok := s.Data[keyCart]
	if !ok || raw == nil {
		out := map[string]int{}
		s.Data[keyCart] = out
		return out
	}
	switch m := raw.(type) {
	case map[string]int:
		return m
	case map[string]int64:
		out := make(map[string]int, len(m))
		for k, v := range m { out[k] = int(v) }
		s.Data[keyCart] = out; return out
	case map[string]float64:
		out := make(map[string]int, len(m))
		for k, v := range m { out[k] = int(v) }
		s.Data[keyCart] = out; return out
	case map[string]interface{}:
		out := make(map[string]int, len(m))
		for k, v := range m {
			switch vv := v.(type) {
			case int:
				out[k] = vv
			case int32:
				out[k] = int(vv)
			case int64:
				out[k] = int(vv)
			case uint:
				out[k] = int(vv)
			case uint32:
				out[k] = int(vv)
			case uint64:
				out[k] = int(vv)
			case float32:
				out[k] = int(vv)
			case float64:
				out[k] = int(vv)
			case json.Number:
				if n, err := vv.Int64(); err == nil { out[k] = int(n) }
			case string:
				if n, err := strconv.Atoi(vv); err == nil { out[k] = n }
			}
		}
		s.Data[keyCart] = out; return out
	case string:
		var tmp map[string]int
		if err := json.Unmarshal([]byte(m), &tmp); err == nil {
			s.Data[keyCart] = tmp; return tmp
		}
	case []byte:
		var tmp map[string]int
		if err := json.Unmarshal(m, &tmp); err == nil {
			s.Data[keyCart] = tmp; return tmp
		}
	}
	out := map[string]int{}
	s.Data[keyCart] = out
	return out
}

func cartSnapshot(s *types.Session) map[string]int {
	c := getCart(s)
	cp := make(map[string]int, len(c))
	for k, v := range c { cp[k] = v }
	return cp
}

func findItem(id string) *Item {
	it, ok := barItemsByID[id]
	if !ok {
		return nil
	}
	cp := it
	return &cp
}

func cartTotalAMD(s *types.Session) int {
	total := 0
	for id, q := range getCart(s) {
		if it := findItem(id); it != nil { total += it.PriceAMD * q }
	}
	return total
}

func renderCartText(s *types.Session, d botengine.Deps) string {
	p := d.Printer(s.Lang)
	if isCartEmpty(s) { return p.Sprintf("bar_cart_empty") }
	var b strings.Builder
	b.WriteString("🧺 <b>")
	b.WriteString(p.Sprintf("bar_cart_title"))
	b.WriteString("</b>\n")
	items := cartSnapshot(s)
	ids := sortedKeys(items)
	for _, id := range ids {
		it := findItem(id); if it == nil { continue }
		qty := items[id]
		b.WriteString(p.Sprintf("bar_line_item", it.ShortName, qty, it.PriceAMD*qty) + "\n")
	}
	b.WriteString(p.Sprintf("bar_cart_total", cartTotalAMD(s)))
	return b.String()
}

// ---------- меню ----------
func renderMenuCompact(d botengine.Deps, s *types.Session) string {
	p := d.Printer(s.Lang)
	var b strings.Builder
	b.WriteString(p.Sprintf("bar_menu_title") + "\n")

	for _, it := range allItems() {
		qty := qtyInCart(s, it.ID)
		if qty > 0 {
			b.WriteString(fmt.Sprintf("• %s — %d AMD (×%d)\n", it.ShortName, it.PriceAMD, qty))
		} else {
			b.WriteString(fmt.Sprintf("• %s — %d AMD\n", it.ShortName, it.PriceAMD))
		}
	}
	b.WriteString(p.Sprintf("bar_menu_hint"))
	return b.String()
}


func menuKeyboardCompact(d botengine.Deps, s *types.Session) tgbotapi.InlineKeyboardMarkup {
	p := d.Printer(s.Lang)
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, it := range allItems() {
		qty := qtyInCart(s, it.ID)

		rows = append(rows, ui.Row(
			ui.Cb(it.ShortName, "bar:noop"),
			// ui.Cb(p.Sprintf("bar_btn_photo"), "bar:peek:"+it.ID),
		))

		rows = append(rows, ui.Row(
			ui.Cb(p.Sprintf("bar_price_qty", it.PriceAMD, qty), "bar:noop"),
			ui.Cb("−", "bar:rem:"+it.ID),
			ui.Cb("+", "bar:add:"+it.ID),
		))
		rows = append(rows, ui.Row(
			ui.Cb(p.Sprintf("bar_price_qty", it.PriceAMD, qty), "bar:noop"),
			// ui.Cb(p.Sprintf("bar_btn_photo"), "bar:peek:"+it.ID),
		))
		rows = append(rows, ui.Row(ui.Cb(" ", "bar:noop")))
	}
	rows = append(rows,
		ui.Row(ui.Cb(p.Sprintf("bar_btn_cart"), "bar:cart"), ui.Cb(p.Sprintf("bar_btn_clear"), "bar:clear")),
		ui.Row(ui.Cb(p.Sprintf("bar_btn_checkout"), "bar:checkout")),
	)
	return ui.Inline(rows...)
}

func sortedKeys(m map[string]int) []string { ks := make([]string, 0, len(m)); for k := range m { ks = append(ks, k) }; sort.Strings(ks); return ks }
func qtyInCart(s *types.Session, id string) int { if s == nil { return 0 }; return getCart(s)[id] }

// func resetBarState(s *types.Session) {
// 	if s.Data == nil { s.Data = map[string]interface{}{} }
// 	delete(s.Data, keyCart)
// 	delete(s.Data, keyBuyer)
// 	delete(s.Data, keyCurrency)
// 	delete(s.Data, keyServe)
// 	delete(s.Data, keyZone)
// 	delete(s.Data, keyOrderID)
// 	delete(s.Data, keyNotes)
// 	delete(s.Data, keyAwaitNotes)
// }

// ---------- персист утилиты ----------

func loadBarStateFromUser(s *types.Session) *BarState {
	if s == nil { return nil }
	raw := s.Data[keyBarStateBlob]
	if raw == nil { return nil }
	switch v := raw.(type) {
	case *BarState:
		return v
	case BarState:
		cp := v
		return &cp
	case map[string]interface{}:
		b, _ := json.Marshal(v)
		var st BarState
		if json.Unmarshal(b, &st) == nil { return &st }
	case string:
		var st BarState
		if json.Unmarshal([]byte(v), &st) == nil { return &st }
	case []byte:
		var st BarState
		if json.Unmarshal(v, &st) == nil { return &st }
	}
	return nil
}

func persistBarStateToUser(s *types.Session) {
	if s == nil { return }
	st := BarState{
		Cart:      cartSnapshot(s), // ← было getCart(s)
		Buyer:     strings.TrimSpace(fmt.Sprint(s.Data[keyBuyer])),
		Serve:     strings.TrimSpace(fmt.Sprint(s.Data[keyServe])),
		Zone:      strings.TrimSpace(fmt.Sprint(s.Data[keyZone])),
		Currency:  func() string { if c := strings.TrimSpace(fmt.Sprint(s.Data[keyCurrency])); c != "" && c != "<nil>" { return c }; return "AMD" }(),
		Notes:     strings.TrimSpace(fmt.Sprint(s.Data[keyNotes])),
		OrderID:   strings.TrimSpace(fmt.Sprint(s.Data[keyOrderID])),
		UpdatedAt: time.Now().Unix(),
	}
	s.Data[keyBarStateBlob] = st
}

func restoreBarStateIntoSession(s *types.Session) {
	if s == nil { return }

	// тянем снепшот; помним, что s.Data на новом апдейте может быть пустым
	st := loadBarStateFromUser(s)

	// гарантируем наличие карты данных
	if s.Data == nil {
		s.Data = map[string]interface{}{}
	}

	// если снепшота нет — просто убедимся, что корзина есть
	if st == nil {
		if _, ok := s.Data[keyCart]; !ok || s.Data[keyCart] == nil {
			s.Data[keyCart] = map[string]int{}
		}
		return
	}

	// 1) корзина: не держим ссылку на сохранённую map, всегда копируем
	//    если в оперативной корзине пусто — поднимаем из снепшота
	cur := getCart(s)
	if len(cur) == 0 && st.Cart != nil {
		cp := make(map[string]int, len(st.Cart))
		for k, v := range st.Cart { cp[k] = v }
		s.Data[keyCart] = cp
	} else if cur == nil { // страховка на случай nil
		s.Data[keyCart] = map[string]int{}
	}

	// 2) простые поля — подставляем, если есть
	if st.Buyer    != "" { s.Data[keyBuyer]    = st.Buyer }
	if st.Serve    != "" { s.Data[keyServe]    = st.Serve }
	if st.Zone     != "" { s.Data[keyZone]     = st.Zone }
	if st.Currency != "" { s.Data[keyCurrency] = st.Currency }
	if st.Notes    != "" { s.Data[keyNotes]    = st.Notes }
	if st.OrderID  != "" { s.Data[keyOrderID]  = st.OrderID }

	// дефолт валюты, если так и не проставили
	if c, ok := s.Data[keyCurrency]; !ok || fmt.Sprint(c) == "" || fmt.Sprint(c) == "<nil>" {
		s.Data[keyCurrency] = "AMD"
	}
}


// ---------- базовая гидрация на каждом апдейте ----------
func ensureSessionHydrated(s *types.Session) {
	if s == nil { return }
	restoreBarStateIntoSession(s)
	// ещё раз гарантируем дефолты
	if s.Data == nil { s.Data = map[string]interface{}{} }
	if _, ok := s.Data[keyCart]; !ok || s.Data[keyCart] == nil {
		s.Data[keyCart] = map[string]int{}
	}
	if _, ok := s.Data[keyCurrency]; !ok || s.Data[keyCurrency] == nil || fmt.Sprint(s.Data[keyCurrency]) == "<nil>" {
		s.Data[keyCurrency] = "AMD"
	}
}


func appendCommandHistory(s *types.Session, cmd, payload string) {
	if s == nil { return }
	now := time.Now().Unix()
	entry := CommandEntry{Cmd: cmd, AtUnix: now, Payload: payload}
	var list []CommandEntry
	switch v := s.Data[keyCmdHistory].(type) {
	case []CommandEntry:
		list = v
	case []interface{}:
		b, _ := json.Marshal(v)
		_ = json.Unmarshal(b, &list)
	case string:
		_ = json.Unmarshal([]byte(v), &list)
	case []byte:
		_ = json.Unmarshal(v, &list)
	}
	list = append(list, entry)
	if len(list) > 200 { list = list[len(list)-200:] }
	s.Data[keyCmdHistory] = list
}

func appendOrderHistory(s *types.Session, ord BarOrder) {
	if s == nil { return }
	var list []BarOrder
	switch v := s.Data[keyOrdersHistory].(type) {
	case []BarOrder:
		list = v
	case []interface{}:
		b, _ := json.Marshal(v)
		_ = json.Unmarshal(b, &list)
	case string:
		_ = json.Unmarshal([]byte(v), &list)
	case []byte:
		_ = json.Unmarshal(v, &list)
	}
	list = append(list, ord)
	if len(list) > 200 { list = list[len(list)-200:] }
	s.Data[keyOrdersHistory] = list
	// обновим снимок
	s.Data[keyOrderID] = ord.OrderID
	persistBarStateToUser(s)
}

// ---------- утилиты ----------
func safeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func generateOrderID() string {
	// Простой монотонный ID на основе времени
	return fmt.Sprintf("ORD-%d", time.Now().Unix())
}

// Кнопка в админском сообщении
func adminOrderKeyboard(userChatID int64, serve, zone, label string) tgbotapi.InlineKeyboardMarkup {
	var payload string
	if serve == "pickup" {
		payload = fmt.Sprintf("bar:done:%d:p", userChatID)
	} else {
		zc := "z"
		switch zone {
		case "coworking": zc = "zcw"
		case "cafe":      zc = "zcf"
		case "street":    zc = "zst"
		}
		payload = fmt.Sprintf("bar:done:%d:%s", userChatID, zc)
	}
	return ui.Inline(ui.Row(ui.Cb(label, payload)))
}

func parseDonePayload(data string) (userID int64, serve, zone string, ok bool) {
	parts := strings.Split(data, ":")
	if len(parts) < 4 { return 0, "", "", false }
	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil { return 0, "", "", false }
	switch parts[3] {
	case "p":   return id, "pickup", "", true
	case "zcw": return id, "tozone", "coworking", true
	case "zcf": return id, "tozone", "cafe", true
	case "zst": return id, "tozone", "street", true
	case "z":   return id, "tozone", "", true
	}
	return id, "", "", false
}

func zoneLabel(p func(string, ...any) string, zone string) string {
	switch zone {
	case "coworking": return p("bar_zone_coworking_name")
	case "cafe":      return p("bar_zone_cafe_name")
	case "street":    return p("bar_zone_street_name")
	default:          return ""
	}
}

func readyText(p func(string, ...any) string, serve, zone string) string {
	switch serve {
	case "pickup":
		return p("bar_ready_pickup")
	case "tozone":
		zl := zoneLabel(p, zone)
		if zl == "" { return p("bar_ready_tozone_generic") }
		return p("bar_ready_tozone_zone", zl)
	default:
		return p("bar_ready_generic")
	}
}

func safeEditMenu(d botengine.Deps, chatID int64, messageID int, txt string, kb tgbotapi.InlineKeyboardMarkup) {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, txt)
	edit.ParseMode = "HTML"
	edit.ReplyMarkup = &kb
	if _, err := d.Bot.Send(edit); err != nil {
		log.Printf("[bar] edit text failed: %v (fallback to markup)", err)
		_, _ = d.Bot.Send(tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, kb))
	}
}

// ---------- шаги ----------
func prompt(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ensureSessionHydrated(s)
	p := d.Printer(s.Lang)

	// восстановим состояние пользователя и запишем команду
	if ev.Kind == botengine.EventCommand && ev.Command == "bar" {
		restoreBarStateIntoSession(s)
		appendCommandHistory(s, "/bar", "")
	}

	if s.Data == nil {
		s.Data = map[string]interface{}{}
	}
	if _, ok := s.Data[keyCart]; !ok {
		s.Data[keyCart] = map[string]int{}
	}
	if _, ok := s.Data[keyCurrency]; !ok {
		s.Data[keyCurrency] = "AMD"
	}

	// грузим категории и меню один раз
	// if err := ensureBarDataLoaded(ctx, d); err != nil {
	// 	log.Printf("[bar] failed to load bar catalog: %v", err)
	// }

	err := loadBarCategoriesFromAPI(ctx, d)
	if err != nil {
		log.Printf("[bar] failed to load bar catalog: %v", err)
	}

	prods, err := fetchProducts(ctx, d) // ← []ProductItem
	if err != nil {
		barMenuErr = err
		log.Printf("[bar] fetchProducts failed: %v", err)
	}

	buildBarMenuFromProducts(prods) 

	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_welcome"))
	log.Printf("[bar] session started for chat %d", s.ChatID)

	// сначала показываем категории
	txt := renderCategoriesText(d, s)
	kb := categoriesKeyboard(d, s)
	_ = ui.SendHTML(d.Bot, s.ChatID, txt, kb)

	persistBarStateToUser(s)
	return stepHandle, nil
}


func handle(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ensureSessionHydrated(s)
	p := d.Printer(s.Lang)
	if ev.Kind != botengine.EventCallback {
		return stepHandle, nil
	}

	// if err := ensureBarDataLoaded(ctx, d); err != nil {
	// 	log.Printf("[bar] failed to load bar catalog in handle: %v", err)
	// }

	data := strings.TrimSpace(ev.CallbackData)
	log.Printf("[bar] cb data=%q chat=%d msg=%d inline=%q", data, s.ChatID, ev.MessageID, ev.InlineMessageID)

	switch {
	case strings.HasPrefix(data, "bar:cat:"):
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "")

		catID := strings.TrimPrefix(data, "bar:cat:")

		// показать всё меню
		if catID == "all" {
			if s.Data != nil {
				delete(s.Data, keyCurrentCategory)
			}
			txt := renderMenuCompact(d, s)
			kb := menuKeyboardCompact(d, s)

			if ev.MessageID != 0 {
				safeEditMenu(d, s.ChatID, ev.MessageID, txt, kb)
			} else {
				_ = ui.SendHTML(d.Bot, s.ChatID, txt, kb)
			}
			return stepHandle, nil
		}

		// вернуться к списку категорий
		if catID == "back" {
			if s.Data != nil {
				delete(s.Data, keyCurrentCategory)
			}
			txt := renderCategoriesText(d, s)
			kb := categoriesKeyboard(d, s)

			if ev.MessageID != 0 {
				safeEditMenu(d, s.ChatID, ev.MessageID, txt, kb)
			} else {
				_ = ui.SendHTML(d.Bot, s.ChatID, txt, kb)
			}
			return stepHandle, nil
		}

		// открыть конкретную категорию
		if s.Data == nil {
			s.Data = map[string]interface{}{}
		}
		s.Data[keyCurrentCategory] = catID

		txt, kb := renderCategoryMenu(d, s, catID)

		if ev.MessageID != 0 {
			safeEditMenu(d, s.ChatID, ev.MessageID, txt, kb)
		} else {
			_ = ui.SendHTML(d.Bot, s.ChatID, txt, kb)
		}
		return stepHandle, nil
	case strings.HasPrefix(data, "bar:done:"):
		// 1) мгновенно гасим спиннер
		if err := ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "Принято"); err != nil {
			log.Printf("[bar] answerCallback failed: %v", err)
		}

		// 2) парсим payload
		userID, serve, zone, ok := parseDonePayload(data)
		if !ok {
			_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "Некорректные данные кнопки")
			return stepHandle, nil
		}

		// 3) уведомление гостю
		txt := readyText(func(key string, a ...any) string { return p.Sprintf(key, a...) }, serve, zone)
		var sendErr error
		for i := 0; i < 2; i++ {
			_, sendErr = d.Bot.Send(tgbotapi.NewMessage(userID, txt))
			if sendErr == nil { break }
			time.Sleep(500 * time.Millisecond)
		}

		// 4) выключаем кнопку в админском сообщении
		doneKB := ui.Inline(ui.Row(ui.Cb(p.Sprintf("bar_admin_issued_label"), "bar:noop")))
		if ev.MessageID != 0 && s.ChatID != 0 {
			if _, err := d.Bot.Send(tgbotapi.NewEditMessageReplyMarkup(s.ChatID, ev.MessageID, doneKB)); err != nil {
				log.Printf("[bar] edit admin markup failed (chat=%d msg=%d): %v", s.ChatID, ev.MessageID, err)
			}
		} else if ev.InlineMessageID != "" {
			cfg := tgbotapi.NewEditMessageReplyMarkup(0, 0, doneKB)
			cfg.InlineMessageID = ev.InlineMessageID
			if _, err := d.Bot.Send(cfg); err != nil {
				log.Printf("[bar] edit inline admin markup failed: %v", err)
			}
		}

		// 5) заметка в админском чате
		if ev.MessageID != 0 && s.ChatID != 0 {
			var note tgbotapi.MessageConfig
			if sendErr != nil {
				note = tgbotapi.NewMessage(s.ChatID, p.Sprintf("bar_admin_notify_fail", userID, sendErr))
			} else {
				note = tgbotapi.NewMessage(s.ChatID, p.Sprintf("bar_admin_user_notified"))
			}
			note.ReplyToMessageID = ev.MessageID
			if _, err := d.Bot.Send(note); err != nil {
				log.Printf("[bar] post admin ack failed: %v", err)
			}
		}
		return stepHandle, nil

	case strings.HasPrefix(data, "bar:serve:"):
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "")
		if strings.HasSuffix(data, ":pickup") {
			if s.Data == nil { s.Data = map[string]interface{}{} }
			s.Data[keyServe] = "pickup"
			delete(s.Data, keyZone)
			persistBarStateToUser(s)
			presentConfirm(d, s); return stepConfirm, nil
		}
		if strings.HasSuffix(data, ":tozone") {
			if s.Data == nil { s.Data = map[string]interface{}{} }
			s.Data[keyServe] = "tozone"
			persistBarStateToUser(s)
			_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_ask_zone"), zoneKeyboard(d, s))
			return stepAskZone, nil
		}
		return stepHandle, nil

	case strings.HasPrefix(data, "bar:zone:"):
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "")
		if s.Data == nil { s.Data = map[string]interface{}{} }
		switch {
		case strings.HasSuffix(data, ":coworking"): s.Data[keyZone] = "coworking"
		case strings.HasSuffix(data, ":cafe"):      s.Data[keyZone] = "cafe"
		case strings.HasSuffix(data, ":street"):    s.Data[keyZone] = "street"
		}
		persistBarStateToUser(s)
		presentConfirm(d, s); return stepConfirm, nil

	case strings.HasPrefix(data, "bar:add:"):
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, p.Sprintf("bar_added"))
		id := strings.TrimPrefix(data, "bar:add:")
		_, changed := addToCart(s, id, +1)
		if changed {
			persistBarStateToUser(s)
			if ev.MessageID != 0 {
				if cat, _ := s.Data[keyCurrentCategory].(string); cat != "" {
					txt, kb := renderCategoryMenu(d, s, cat)
					safeEditMenu(d, s.ChatID, ev.MessageID, txt, kb)
				} else {
					kb := menuKeyboardCompact(d, s)
					txt := renderMenuCompact(d, s)
					safeEditMenu(d, s.ChatID, ev.MessageID, txt, kb)
				}
			}
		}
		return stepHandle, nil

	case strings.HasPrefix(data, "bar:rem:"):
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, p.Sprintf("bar_removed"))
		id := strings.TrimPrefix(data, "bar:rem:")
		_, changed := addToCart(s, id, -1)
		if changed {
			persistBarStateToUser(s)
			if ev.MessageID != 0 {
				if cat, _ := s.Data[keyCurrentCategory].(string); cat != "" {
					txt, kb := renderCategoryMenu(d, s, cat)
					safeEditMenu(d, s.ChatID, ev.MessageID, txt, kb)
				} else {
					kb := menuKeyboardCompact(d, s)
					txt := renderMenuCompact(d, s)
					safeEditMenu(d, s.ChatID, ev.MessageID, txt, kb)
				}
			}
		}
		return stepHandle, nil


	case strings.HasPrefix(data, "bar:peek:"):
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "")
		id := strings.TrimPrefix(data, "bar:peek:")
		if it := findItem(id); it != nil && strings.TrimSpace(it.PhotoURL) != "" {
			pm := tgbotapi.NewPhoto(s.ChatID, tgbotapi.FileURL(it.PhotoURL))
			pm.Caption = fmt.Sprintf("%s — %d AMD", it.ShortName, it.PriceAMD)
			msg, _ := d.Bot.Send(pm)
			go func(chatID int64, messageID int) {
				time.Sleep(8 * time.Second)
				_, _ = d.Bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID))
			}(s.ChatID, msg.MessageID)
		}
		return stepHandle, nil

	case data == "bar:cart":
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "")
		txt := renderCartText(s, d)
		kb := cartKeyboard(d, s)
		_ = ui.SendHTML(d.Bot, s.ChatID, txt, kb)
		return stepHandle, nil

	case strings.HasPrefix(data, "bar:rmitem:"):
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, p.Sprintf("bar_removed"))
		id := strings.TrimPrefix(data, "bar:rmitem:")
		if _, changed := addToCart(s, id, -1); changed {
			persistBarStateToUser(s)
		}
		txt := renderCartText(s, d)
		kb := cartKeyboard(d, s)
		_ = ui.SendHTML(d.Bot, s.ChatID, txt, kb)
		return stepHandle, nil

	case data == "bar:clear":
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, p.Sprintf("bar_cart_cleared"))
		clearCart(s)
		persistBarStateToUser(s)
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_cart_cleared"))
		if ev.MessageID != 0 {
			kb := menuKeyboardCompact(d, s)
			txt := renderMenuCompact(d, s)
			safeEditMenu(d, s.ChatID, ev.MessageID, txt, kb)
		}
		return stepHandle, nil

	case data == "bar:checkout":
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "")
		if isCartEmpty(s) {
			_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_cart_empty"))
			return stepHandle, nil
		}
		buyer := strings.TrimSpace(fmt.Sprint(s.Data[keyBuyer]))
		if buyer != "" && buyer != "<nil>" {
			_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_ask_serve"), serveKeyboard(d, s))
			return stepAskServe, nil
		}
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_ask_name"))
		return stepAskName, nil

	case data == "bar:noop":
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "")
		return stepHandle, nil
	}
	return stepHandle, nil
}

// ---------- имя ----------
func askName(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ensureSessionHydrated(s)
	p := d.Printer(s.Lang)
	if ev.Kind != botengine.EventText {
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_ask_name_hint"))
		return stepAskName, nil
	}
	name := strings.TrimSpace(ev.Text)
	if name == "" {
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_ask_name_hint"))
		return stepAskName, nil
	}
	if s.Data == nil { s.Data = map[string]interface{}{} }
	s.Data[keyBuyer] = name
	persistBarStateToUser(s)
	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_ask_serve"), serveKeyboard(d, s))
	return stepAskServe, nil
}

// ---------- подача ----------
func askServe(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ensureSessionHydrated(s)
	p := d.Printer(s.Lang)
	if ev.Kind != botengine.EventCallback {
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_ask_serve"), serveKeyboard(d, s))
		return stepAskServe, nil
	}
	_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "")
	switch ev.CallbackData {
	case "bar:serve:pickup":
		if s.Data == nil { s.Data = map[string]interface{}{} }
		s.Data[keyServe] = "pickup"; delete(s.Data, keyZone)
		persistBarStateToUser(s)
		presentConfirm(d, s); return stepConfirm, nil
	case "bar:serve:tozone":
		if s.Data == nil { s.Data = map[string]interface{}{} }
		s.Data[keyServe] = "tozone"
		persistBarStateToUser(s)
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_ask_zone"), zoneKeyboard(d, s))
		return stepAskZone, nil
	case "bar:cancel":
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_order_cancelled"))
		clearCart(s)
		delete(s.Data, keyBuyer)
		delete(s.Data, keyServe)
		delete(s.Data, keyZone)
		delete(s.Data, keyOrderID)
		delete(s.Data, keyNotes)
		delete(s.Data, keyAwaitNotes)
		persistBarStateToUser(s)
		s.Flow, s.Step = "", ""
		return stepDone, nil
	default:
		return stepAskServe, nil
	}
}

// ---------- зона ----------
func askZone(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ensureSessionHydrated(s)
	p := d.Printer(s.Lang)
	if ev.Kind != botengine.EventCallback {
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_ask_zone"), zoneKeyboard(d, s))
		return stepAskZone, nil
	}
	_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "")
	data := strings.TrimSpace(ev.CallbackData)
	switch data {
	case "bar:zone:coworking":
		if s.Data == nil { s.Data = map[string]interface{}{} }
		s.Data[keyZone] = "coworking"; persistBarStateToUser(s); presentConfirm(d, s); return stepConfirm, nil
	case "bar:zone:cafe":
		if s.Data == nil { s.Data = map[string]interface{}{} }
		s.Data[keyZone] = "cafe"; persistBarStateToUser(s); presentConfirm(d, s); return stepConfirm, nil
	case "bar:zone:street":
		if s.Data == nil { s.Data = map[string]interface{}{} }
		s.Data[keyZone] = "street"; persistBarStateToUser(s); presentConfirm(d, s); return stepConfirm, nil
	case "bar:cancel":
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_order_cancelled"))
		clearCart(s); delete(s.Data, keyBuyer); delete(s.Data, keyServe); delete(s.Data, keyZone); delete(s.Data, keyNotes); delete(s.Data, keyAwaitNotes)
		persistBarStateToUser(s)
		s.Flow, s.Step = "", ""
		return stepDone, nil
	default:
		return stepAskZone, nil
	}
}

// ---------- подтверждение ----------
func presentConfirm(d botengine.Deps, s *types.Session) {
	ensureSessionHydrated(s)
	p := d.Printer(s.Lang)
	buyer := fmt.Sprint(s.Data[keyBuyer])
	summary := orderServeSummary(d, s)

	var b strings.Builder
	b.WriteString(renderCartText(s, d))
	b.WriteString("\n\n")
	b.WriteString(p.Sprintf("bar_buyer_is", buyer))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("📍 Подача: <b>%s</b>", summary))

	if notesRaw, ok := s.Data[keyNotes]; ok {
		notes := strings.TrimSpace(fmt.Sprint(notesRaw))
		if notes != "" {
			b.WriteString("\n")
			b.WriteString(p.Sprintf("bar_comment_label") + " ") // ← пробел
			b.WriteString(safeHTML(notes))
		}
	}

	b.WriteString("\n")
	b.WriteString(p.Sprintf("bar_contact_hint", baristaContact))

	kb := confirmKeyboard(d, s)
	_ = ui.SendHTML(d.Bot, s.ChatID, b.String(), kb)
}

func orderServeSummary(d botengine.Deps, s *types.Session) string {
	ensureSessionHydrated(s)
	p := d.Printer(s.Lang)
	serve := fmt.Sprint(s.Data[keyServe])
	// zone  := fmt.Sprint(s.Data[keyZone])
	switch serve {
	case "pickup":
		return p.Sprintf("bar_serve_pickup_label")
	// case "tozone":
	// 	if zone == "" { return p.Sprintf("bar_serve_tozone_label") }
	// 	return p.Sprintf("bar_serve_tozone_with_label", zoneLabel(func(key string, a ...any) string { return p.Sprintf(key, a...) }, zone))
	default:
		return p.Sprintf("bar_not_specified")
	}
}

func barHome(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ensureSessionHydrated(s)
	p := d.Printer(s.Lang)

    // Если хочешь кнопку "Full menu" — см. пункт 3 ниже (cmd-роутер).
    // Без кнопки просто отправляем текст с /menu.
    _ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_overview"), nil)

    // После обзора — в твой первый шаг (как раньше было)
    return stepAskServe, nil
}

func confirm(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	ensureSessionHydrated(s)
	p := d.Printer(s.Lang)

	// 1) если пришёл текст и мы ждём комментарий — сохраняем
	if ev.Kind == botengine.EventText && fmt.Sprint(s.Data[keyAwaitNotes]) == "1" {
		notes := strings.TrimSpace(ev.Text)
		if len(notes) > 300 { notes = notes[:300] }
		if s.Data == nil { s.Data = map[string]interface{}{} }
		s.Data[keyNotes] = notes
		delete(s.Data, keyAwaitNotes)
		persistBarStateToUser(s)

		ack, _ := d.Bot.Send(tgbotapi.NewMessage(s.ChatID, p.Sprintf("bar_notes_saved")))
		go func(chatID int64, msgID int) {
			time.Sleep(2 * time.Second)
			_, _ = d.Bot.Request(tgbotapi.NewDeleteMessage(chatID, msgID))
		}(s.ChatID, ack.MessageID)

		presentConfirm(d, s)
		return stepConfirm, nil
	}

	// 2) обычные колбэки
	if ev.Kind != botengine.EventCallback {
		return stepConfirm, nil
	}
	_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "")

	switch ev.CallbackData {
	case "bar:confirm":
		buyer := fmt.Sprint(s.Data[keyBuyer])
		items := cartSnapshot(s)
		total := cartTotalAMD(s)
		oid := strings.TrimSpace(fmt.Sprint(s.Data[keyOrderID]))
		if oid == "" || oid == "<nil>" {
			oid = generateOrderID()
			if s.Data == nil { s.Data = map[string]interface{}{} }
			s.Data[keyOrderID] = oid
		}

		// текст в админский чат
		var b strings.Builder
		b.WriteString("🔔 ")
		b.WriteString(baristaMention)
		b.WriteString("\n")
		b.WriteString(p.Sprintf("bar_admin_new_order_title") + "\n")
		b.WriteString(p.Sprintf("bar_admin_order_no", oid) + "\n")
		b.WriteString(p.Sprintf("bar_admin_name", buyer) + "\n")
		b.WriteString("— — — — — — —\n")
		ids := sortedKeys(items)
		for _, id := range ids {
			it := findItem(id); if it == nil { continue }
			qty := items[id]
			b.WriteString(p.Sprintf("bar_line_item", it.ShortName, qty, it.PriceAMD*qty) + "\n")
		}
		b.WriteString(p.Sprintf("bar_cart_total", total) + "\n")
		b.WriteString(p.Sprintf("bar_admin_serve_line", orderServeSummary(d, s)) + "\n")
		if notesRaw, ok := s.Data[keyNotes]; ok {
			notes := strings.TrimSpace(fmt.Sprint(notesRaw))
			if notes != "" {
				b.WriteString(p.Sprintf("bar_comment_label") + " " + safeHTML(notes) + "\n")
			}
		}
		b.WriteString("— — — — — — —\n")
		b.WriteString(p.Sprintf("bar_admin_questions_title") + "\n")
		b.WriteString(p.Sprintf("bar_admin_q_delivery")    + "\n")
		b.WriteString(p.Sprintf("bar_admin_q_disposables") + "\n")
		b.WriteString(p.Sprintf("bar_admin_q_time")        + "\n")
		b.WriteString(p.Sprintf("bar_admin_q_payment")     + "\n")
		b.WriteString(p.Sprintf("bar_admin_contact_line", baristaContact) + "\n")
		b.WriteString(p.Sprintf("bar_admin_contact_meta", ev.FromUserName, s.ChatID) + "\n")

		// история заказа
		var orderItems []BarOrderItem
		for _, id := range ids {
			it := findItem(id); if it == nil { continue }
			qty := items[id]
			orderItems = append(orderItems, BarOrderItem{
				ID:       id,
				Title:    it.ShortName,
				Qty:      qty,
				PriceAMD: it.PriceAMD,
				SumAMD:   it.PriceAMD * qty,
			})
		}
		ord := BarOrder{
			OrderID:   oid,
			Buyer:     buyer,
			Items:     orderItems,
			TotalAMD:  total,
			Serve:     fmt.Sprint(s.Data[keyServe]),
			Zone:      fmt.Sprint(s.Data[keyZone]),
			Notes:     strings.TrimSpace(fmt.Sprint(s.Data[keyNotes])),
			CreatedAt: time.Now().Unix(),
		}
		appendOrderHistory(s, ord)
		persistBarStateToUser(s)

		// отправка в админский чат с fallback
		targetChat := d.Cfg.OrdersChatId
		if targetChat == 0 { targetChat = d.Cfg.AdminChatId }

		serve := fmt.Sprint(s.Data[keyServe])
		zone  := fmt.Sprint(s.Data[keyZone])

		msg := tgbotapi.NewMessage(targetChat, b.String())
		msg.ParseMode  = "HTML"
		msg.ReplyMarkup = adminOrderKeyboard(s.ChatID, serve, zone, p.Sprintf("bar_admin_ready_btn"))
		if _, err := d.Bot.Send(msg); err != nil {
			log.Printf("[bar] FAILED to send order to chat=%d: %v", targetChat, err)
			if targetChat != d.Cfg.AdminChatId && d.Cfg.AdminChatId != 0 {
				fb := tgbotapi.NewMessage(d.Cfg.AdminChatId,
					fmt.Sprintf("⚠️ Не удалось отправить заказ в чат %d: %v\n\n%s", targetChat, err, b.String()))
				fb.ParseMode = "HTML"
				_, _ = d.Bot.Send(fb)
			}
		}

		// подтверждение пользователю
		var contactLink string
		if ev.FromUserName != "" {
			contactLink = fmt.Sprintf("<a href=\"https://t.me/%s\">@%s</a>", ev.FromUserName, ev.FromUserName)
		} else {
			contactLink = fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>", s.ChatID, p.Sprintf("bar_open_chat"))
		}
		confTxt := p.Sprintf("bar_order_sent") + "\n" +
			p.Sprintf("bar_order_number_label", oid) + "\n" +
			p.Sprintf("bar_order_customer_label", buyer) + "\n" +
			p.Sprintf("bar_chat_label", contactLink) + "\n\n" +
			p.Sprintf("bar_contact_hint", baristaContact)
		_ = ui.SendHTML(d.Bot, s.ChatID, confTxt)

		// очистка оперативного состояния
		clearCart(s)
		delete(s.Data, keyBuyer)
		delete(s.Data, keyServe)
		delete(s.Data, keyZone)
		delete(s.Data, keyOrderID)
		delete(s.Data, keyNotes)
		delete(s.Data, keyAwaitNotes)
		persistBarStateToUser(s)

		s.Flow, s.Step = "", ""
		return stepDone, nil

	case "bar:notes":
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, p.Sprintf("bar_notes_toast_prompt"))
		if s.Data == nil { s.Data = map[string]interface{}{} }
		s.Data[keyAwaitNotes] = "1"
		persistBarStateToUser(s)
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_notes_enter"), notesCancelKeyboard(d, s))
		return stepConfirm, nil

	case "bar:notes:clear":
		delete(s.Data, keyNotes)
		persistBarStateToUser(s)
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, p.Sprintf("bar_notes_deleted"))
		presentConfirm(d, s)
		return stepConfirm, nil

	case "bar:notes:cancel":
		delete(s.Data, keyAwaitNotes)
		persistBarStateToUser(s)
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, p.Sprintf("bar_notes_unchanged"))
		presentConfirm(d, s)
		return stepConfirm, nil

	case "bar:cancel":
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_order_cancelled"))
		clearCart(s)
		delete(s.Data, keyBuyer)
		delete(s.Data, keyNotes)
		delete(s.Data, keyAwaitNotes)
		persistBarStateToUser(s)
		s.Flow, s.Step = "", ""
		return stepDone, nil
	}

	return stepConfirm, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return stepDone, nil
}

func cartKeyboard(d botengine.Deps, s *types.Session) tgbotapi.InlineKeyboardMarkup {
	p := d.Printer(s.Lang)
	items := cartSnapshot(s)
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, id := range sortedKeys(items) {
		it := findItem(id); if it == nil { continue }
		rows = append(rows, ui.Row(ui.Cb(fmt.Sprintf("❌ %s (×%d)", it.ShortName, items[id]), "bar:rmitem:"+id)))
	}
	rows = append(rows,
		ui.Row(ui.Cb(p.Sprintf("bar_btn_clear"), "bar:clear")),
		ui.Row(ui.Cb(p.Sprintf("bar_btn_checkout"), "bar:checkout")),
	)
	return ui.Inline(rows...)
}

func confirmKeyboard(d botengine.Deps, s *types.Session) tgbotapi.InlineKeyboardMarkup {
	if s == nil {
		p := d.Printer("ru")
		var rows [][]tgbotapi.InlineKeyboardButton
		rows = append(rows, ui.Row(ui.Cb(p.Sprintf("bar_btn_confirm"), "bar:confirm")))
		rows = append(rows, ui.Row(ui.Cb(p.Sprintf("bar_btn_add_note"), "bar:notes")))
		rows = append(rows, ui.Row(ui.Cb(p.Sprintf("bar_btn_cancel"), "bar:cancel")))
		return ui.Inline(rows...)
	}
	p := d.Printer(s.Lang)
	hasNotes := false
	if s.Data != nil {
		if n, ok := s.Data[keyNotes]; ok && strings.TrimSpace(fmt.Sprint(n)) != "" {
			hasNotes = true
		}
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, ui.Row(ui.Cb(p.Sprintf("bar_btn_confirm"), "bar:confirm")))
	if hasNotes {
		rows = append(rows,
			ui.Row(ui.Cb(p.Sprintf("bar_btn_edit_note"), "bar:notes")),
			ui.Row(ui.Cb(p.Sprintf("bar_btn_delete_note"), "bar:notes:clear")),
		)
	} else {
		rows = append(rows, ui.Row(ui.Cb(p.Sprintf("bar_btn_add_note"), "bar:notes")))
	}
	rows = append(rows, ui.Row(ui.Cb(p.Sprintf("bar_btn_cancel"), "bar:cancel")))
	return ui.Inline(rows...)
}

func notesCancelKeyboard(d botengine.Deps, s *types.Session) tgbotapi.InlineKeyboardMarkup {
	p := d.Printer(s.Lang)
	return ui.Inline(ui.Row(ui.Cb(p.Sprintf("bar_btn_back"), "bar:notes:cancel")))
}

func zoneKeyboard(d botengine.Deps, s *types.Session) tgbotapi.InlineKeyboardMarkup {
	p := d.Printer(s.Lang)
	return ui.Inline(
		ui.Row(ui.Cb("💻 "+p.Sprintf("bar_zone_coworking_name"), "bar:zone:coworking")),
		ui.Row(ui.Cb("☕️ "+p.Sprintf("bar_zone_cafe_name"),      "bar:zone:cafe")),
		ui.Row(ui.Cb("🌳 "+p.Sprintf("bar_zone_street_name"),    "bar:zone:street")),
	)
}

func serveKeyboard(d botengine.Deps, s *types.Session) tgbotapi.InlineKeyboardMarkup {
	p := d.Printer(s.Lang)
	return ui.Inline(
		ui.Row(ui.Cb(p.Sprintf("bar_serve_pickup_btn"), "bar:serve:pickup")),
		ui.Row(ui.Cb(p.Sprintf("bar_serve_tozone_btn"), "bar:serve:tozone")),
		ui.Row(ui.Cb(p.Sprintf("bar_btn_cancel"), "bar:cancel")),
	)
}
