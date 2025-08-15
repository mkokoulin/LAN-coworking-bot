package flow

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
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

var (
	stepAskServe types.Step = "bar:ask_serve"
	stepAskZone  types.Step = "bar:ask_zone"
	// остальные шаги (stepHandle, stepAskName, stepConfirm, stepDone) объявлены в проекте
)

const (
	baristaContact = "@LAN_Barista" // для текста клиенту
	baristaMention = "@lan_barista" // для пинга в админском чате
)

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
			case int: out[k] = vv
			case int64: out[k] = int(vv)
			case float64: out[k] = int(vv)
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
	for _, it := range getMenu() {
		if it.ID == id { cp := it; return &cp }
	}
	return nil
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
		b.WriteString(p.Sprintf("bar_line_item", it.Title, qty, it.PriceAMD*qty) + "\n")
	}
	b.WriteString(p.Sprintf("bar_cart_total", cartTotalAMD(s)))
	return b.String()
}

// ---------- меню ----------
func renderMenuCompact(d botengine.Deps, s *types.Session) string {
	p := d.Printer(s.Lang)
	var b strings.Builder
	b.WriteString(p.Sprintf("bar_menu_title") + "\n")
	for _, it := range getMenu() {
		qty := qtyInCart(s, it.ID)
		if qty > 0 {
			b.WriteString(fmt.Sprintf("• %s — %d AMD (×%d)\n", it.Title, it.PriceAMD, qty))
		} else {
			b.WriteString(fmt.Sprintf("• %s — %d AMD\n", it.Title, it.PriceAMD))
		}
	}
	b.WriteString(p.Sprintf("bar_menu_hint"))
	return b.String()
}

func menuKeyboardCompact(d botengine.Deps, s *types.Session) tgbotapi.InlineKeyboardMarkup {
	p := d.Printer(s.Lang)
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, it := range getMenu() {
		qty := qtyInCart(s, it.ID)
		rows = append(rows, ui.Row(
			ui.Cb("−", "bar:rem:"+it.ID),
			ui.Cb(it.Title, "bar:noop"),
			ui.Cb("+", "bar:add:"+it.ID),
		))
		rows = append(rows, ui.Row(
			ui.Cb(p.Sprintf("bar_price_qty", it.PriceAMD, qty), "bar:noop"),
			ui.Cb(p.Sprintf("bar_btn_photo"), "bar:peek:"+it.ID),
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

func resetBarState(s *types.Session) {
	if s.Data == nil { s.Data = map[string]interface{}{} }
	delete(s.Data, keyCart)
	delete(s.Data, keyBuyer)
	delete(s.Data, keyCurrency)
	delete(s.Data, keyServe)
	delete(s.Data, keyZone)
	delete(s.Data, keyOrderID)
	delete(s.Data, keyNotes)
	delete(s.Data, keyAwaitNotes)
}

// ---------- утилиты ----------
func safeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
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
	p := d.Printer(s.Lang)
	if ev.Kind == botengine.EventCommand && ev.Command == "bar" { resetBarState(s) }
	if s.Data == nil { s.Data = map[string]interface{}{} }
	if _, ok := s.Data[keyCart]; !ok { s.Data[keyCart] = map[string]int{} }
	if _, ok := s.Data[keyCurrency]; !ok { s.Data[keyCurrency] = "AMD" }

	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_welcome"))
	log.Printf("[bar] session started for chat %d", s.ChatID)
	txt := renderMenuCompact(d, s)
	kb  := menuKeyboardCompact(d, s)
	_ = ui.SendHTML(d.Bot, s.ChatID, txt, kb)
	return stepHandle, nil
}

func handle(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)
	if ev.Kind != botengine.EventCallback { return stepHandle, nil }

	data := strings.TrimSpace(ev.CallbackData)
	log.Printf("[bar] cb data=%q chat=%d msg=%d inline=%q", data, s.ChatID, ev.MessageID, ev.InlineMessageID)

	switch {
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

			// 3) пробуем отправить уведомление гостю (с небольшим ретраем)
			txt := readyText(func(key string, a ...any) string { return p.Sprintf(key, a...) }, serve, zone)
			var sendErr error
			for i := 0; i < 2; i++ {
				_, sendErr = d.Bot.Send(tgbotapi.NewMessage(userID, txt))
				if sendErr == nil {
					break
				}
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

			// 5) короткая заметка в админском чате (успех/ошибка), чтобы было видно, что произошло
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
			s.Data[keyServe] = "pickup"; delete(s.Data, keyZone)
			presentConfirm(d, s); return stepConfirm, nil
		}
		if strings.HasSuffix(data, ":tozone") {
			if s.Data == nil { s.Data = map[string]interface{}{} }
			s.Data[keyServe] = "tozone"
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
		presentConfirm(d, s); return stepConfirm, nil

	case strings.HasPrefix(data, "bar:add:"):
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, p.Sprintf("bar_added"))
		id := strings.TrimPrefix(data, "bar:add:")
		_, changed := addToCart(s, id, +1)
		if changed && ev.MessageID != 0 {
			kb := menuKeyboardCompact(d, s)
			txt := renderMenuCompact(d, s)
			safeEditMenu(d, s.ChatID, ev.MessageID, txt, kb)
		}
		return stepHandle, nil

	case strings.HasPrefix(data, "bar:rem:"):
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, p.Sprintf("bar_removed"))
		id := strings.TrimPrefix(data, "bar:rem:")
		_, changed := addToCart(s, id, -1)
		if changed && ev.MessageID != 0 {
			txt := renderCartText(s, d)
			kb := cartKeyboard(d, s)
			safeEditMenu(d, s.ChatID, ev.MessageID, txt, kb)
		}
		return stepHandle, nil

	case strings.HasPrefix(data, "bar:peek:"):
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "")
		id := strings.TrimPrefix(data, "bar:peek:")
		if it := findItem(id); it != nil && strings.TrimSpace(it.PhotoURL) != "" {
			pm := tgbotapi.NewPhoto(s.ChatID, tgbotapi.FileURL(it.PhotoURL))
			pm.Caption = fmt.Sprintf("%s — %d AMD", it.Title, it.PriceAMD)
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
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "")
		id := strings.TrimPrefix(data, "bar:rmitem:")
		removeItem(s, id)
		txt := renderCartText(s, d)
		kb := cartKeyboard(d, s)
		_ = ui.SendHTML(d.Bot, s.ChatID, txt, kb)
		return stepHandle, nil

	case data == "bar:clear":
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, p.Sprintf("bar_cart_cleared"))
		clearCart(s)
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
	_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_ask_serve"), serveKeyboard(d, s))
	return stepAskServe, nil
}

// ---------- подача ----------
func askServe(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)
	if ev.Kind != botengine.EventCallback {
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_ask_serve_hint"), serveKeyboard(d, s))
		return stepAskServe, nil
	}
	_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, "")
	switch ev.CallbackData {
	case "bar:serve:pickup":
		if s.Data == nil { s.Data = map[string]interface{}{} }
		s.Data[keyServe] = "pickup"; delete(s.Data, keyZone)
		presentConfirm(d, s); return stepConfirm, nil
	case "bar:serve:tozone":
		if s.Data == nil { s.Data = map[string]interface{}{} }
		s.Data[keyServe] = "tozone"
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
		s.Flow, s.Step = "", ""
		return stepDone, nil
	default:
		return stepAskServe, nil
	}
}

// ---------- зона ----------
func askZone(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
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
		s.Data[keyZone] = "coworking"; presentConfirm(d, s); return stepConfirm, nil
	case "bar:zone:cafe":
		if s.Data == nil { s.Data = map[string]interface{}{} }
		s.Data[keyZone] = "cafe"; presentConfirm(d, s); return stepConfirm, nil
	case "bar:zone:street":
		if s.Data == nil { s.Data = map[string]interface{}{} }
		s.Data[keyZone] = "street"; presentConfirm(d, s); return stepConfirm, nil
	case "bar:cancel":
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_order_cancelled"))
		clearCart(s); delete(s.Data, keyBuyer); delete(s.Data, keyServe); delete(s.Data, keyZone); delete(s.Data, keyNotes); delete(s.Data, keyAwaitNotes)
		s.Flow, s.Step = "", ""
		return stepDone, nil
	default:
		return stepAskZone, nil
	}
}

// ---------- подтверждение ----------
func presentConfirm(d botengine.Deps, s *types.Session) {
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
			b.WriteString(p.Sprintf("bar_comment_label"))
			b.WriteString(safeHTML(notes))
		}
	}

	b.WriteString("\n")
	b.WriteString(p.Sprintf("bar_contact_hint", baristaContact))

	kb := confirmKeyboard(d, s)
	_ = ui.SendHTML(d.Bot, s.ChatID, b.String(), kb)
}

func orderServeSummary(d botengine.Deps, s *types.Session) string {
	p := d.Printer(s.Lang)
	serve := fmt.Sprint(s.Data[keyServe])
	zone  := fmt.Sprint(s.Data[keyZone])
	switch serve {
	case "pickup":
		return p.Sprintf("bar_serve_pickup_label")
	case "tozone":
		if zone == "" { return p.Sprintf("bar_serve_tozone_label") }
		return p.Sprintf("bar_serve_tozone_with_label", zoneLabel(func(key string, a ...any) string { return p.Sprintf(key, a...) }, zone))
	default:
		return p.Sprintf("bar_not_specified")
	}
}

func confirm(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	p := d.Printer(s.Lang)

	// 1) если пришёл текст и мы ждём комментарий — сохраняем
	if ev.Kind == botengine.EventText && fmt.Sprint(s.Data[keyAwaitNotes]) == "1" {
		notes := strings.TrimSpace(ev.Text)
		if len(notes) > 300 { notes = notes[:300] }
		if s.Data == nil { s.Data = map[string]interface{}{} }
		s.Data[keyNotes] = notes
		delete(s.Data, keyAwaitNotes)

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
		oid := fmt.Sprint(s.Data[keyOrderID])
		if strings.TrimSpace(oid) == "" || oid == "<nil>" {
			// oid = generateCatOrderID()
			if s.Data == nil { s.Data = map[string]interface{}{} }
			s.Data[keyOrderID] = oid
		}

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
			b.WriteString(p.Sprintf("bar_line_item", it.Title, qty, it.PriceAMD*qty) + "\n")
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

		targetChat := d.Cfg.OrdersChatId
		if targetChat == 0 { targetChat = d.Cfg.AdminChatId }

		serve := fmt.Sprint(s.Data[keyServe])
		zone  := fmt.Sprint(s.Data[keyZone])

		msg := tgbotapi.NewMessage(targetChat, b.String())
		msg.ParseMode  = "HTML"
		msg.ReplyMarkup = adminOrderKeyboard(s.ChatID, serve, zone, p.Sprintf("bar_admin_ready_btn"))
		_, _ = d.Bot.Send(msg) // обработка ошибок — как у тебя

		// ссылка «открыть чат» — локализация подписи
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

		clearCart(s)
		delete(s.Data, keyBuyer)
		delete(s.Data, keyServe)
		delete(s.Data, keyZone)
		delete(s.Data, keyOrderID)
		delete(s.Data, keyNotes)
		delete(s.Data, keyAwaitNotes)
		s.Flow, s.Step = "", ""
		return stepDone, nil

	case "bar:notes":
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, p.Sprintf("bar_notes_toast_prompt"))
		if s.Data == nil { s.Data = map[string]interface{}{} }
		s.Data[keyAwaitNotes] = "1"
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_notes_enter"), notesCancelKeyboard(d, s))
		return stepConfirm, nil

	case "bar:notes:clear":
		delete(s.Data, keyNotes)
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, p.Sprintf("bar_notes_deleted"))
		presentConfirm(d, s)
		return stepConfirm, nil

	case "bar:notes:cancel":
		delete(s.Data, keyAwaitNotes)
		_ = ui.AnswerCallback(d.Bot, ev.CallbackQueryID, p.Sprintf("bar_notes_unchanged"))
		presentConfirm(d, s)
		return stepConfirm, nil

	case "bar:cancel":
		_ = ui.SendHTML(d.Bot, s.ChatID, p.Sprintf("bar_order_cancelled"))
		clearCart(s)
		delete(s.Data, keyBuyer)
		delete(s.Data, keyNotes)
		delete(s.Data, keyAwaitNotes)
		s.Flow, s.Step = "", ""
		return stepDone, nil
	}

	return stepConfirm, nil
}

func done(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	return stepDone, nil
}

// ---------- генератор номера ----------
// var catTitles = []string{"Кот", "Котэ", "Господин Кот", "Сэр Мур"}
// var catNames  = []string{"Барсик", "Сметана", "Кекс", "Пельмень", "Жмых", "Пончик", "Граф Лапкин", "Мурчестер", "Васаби", "Шпрот"}
// var catColors = []string{"рыжий", "серебристо-полосатый", "трёхцветный", "чёрный как эспрессо", "снежный", "дымчатый", "в горошек (почти)"}
// var catTraits = []string{"охотник на коробки", "надзиратель за кружками", "специалист по колбасе", "мурчательный", "шуршолог", "прыг-скок", "дример на подоконнике"}
// var tokenRunes = []rune("23456789ABCDEFGHJKLMNPQRSTUVWXYZ")

func randInt(n int64) int { x, _ := rand.Int(rand.Reader, big.NewInt(n)); return int(x.Int64()) }
func pick[T any](arr []T) T { return arr[randInt(int64(len(arr)))] }
// func shortToken(n int) string { b := make([]rune, n); for i := range b { b[i] = tokenRunes[randInt(int64(len(tokenRunes)))] }; return string(b) }
// func generateCatOrderID() string {
// 	title := pick(catTitles); name := pick(catNames); color := pick(catColors); trait := pick(catTraits); code := shortToken(4)
// 	return fmt.Sprintf("%s %s — %s, %s • %s", title, name, color, trait, code)
// }

// ---------- клавиатуры ----------
func itemKeyboard(d botengine.Deps, s *types.Session, id string, qty int) tgbotapi.InlineKeyboardMarkup {
	p := d.Printer(s.Lang)
	return ui.Inline(ui.Row(
		ui.Cb("−", "bar:rem:"+id),
		ui.Cb(fmt.Sprintf(p.Sprintf("bar_in_cart_label"), qty), "bar:noop"),
		ui.Cb("+", "bar:add:"+id),
	))
}

func cartKeyboard(d botengine.Deps, s *types.Session) tgbotapi.InlineKeyboardMarkup {
	p := d.Printer(s.Lang)
	items := cartSnapshot(s)
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, id := range sortedKeys(items) {
		it := findItem(id); if it == nil { continue }
		rows = append(rows, ui.Row(ui.Cb(fmt.Sprintf("❌ %s (×%d)", it.Title, items[id]), "bar:rmitem:"+id)))
	}
	rows = append(rows,
		ui.Row(ui.Cb(p.Sprintf("bar_btn_clear"), "bar:clear")),
		ui.Row(ui.Cb(p.Sprintf("bar_btn_checkout"), "bar:checkout")),
	)
	return ui.Inline(rows...)
}

func confirmKeyboard(d botengine.Deps, s *types.Session) tgbotapi.InlineKeyboardMarkup {
	p := d.Printer(s.Lang)
	hasNotes := false
	if s != nil && s.Data != nil {
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
