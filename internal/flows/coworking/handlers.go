package flows

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"github.com/mkokoulin/LAN-coworking-bot/internal/ui"
)

const (
	dataAuthStatus = "coworking_auth_status" // pending|approved|rejected
	dataFullName   = "coworking_full_name"
	dataPhone      = "coworking_phone"
	dataTariff     = "coworking_tariff"
	dataReqType    = "coworking_request_type" // new|relink
)

func coworkingHome(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if ev.Kind == botengine.EventCallback {
		switch ev.CallbackData {
		case "cw:new":
			s.Step = CoworkingNewName
			return botengine.InternalContinue, nil

		case "cw:returning":
			s.Step = CoworkingReturningPhone
			return botengine.InternalContinue, nil

		case "cw:confirm":
			s.Step = CoworkingConfirm
			return botengine.InternalContinue, nil

		case "cw:back:tariff":
			s.Step = CoworkingNewTariff
			return botengine.InternalContinue, nil

		case "cw:tariff:hour":
			ensureData(s)
			s.Data[dataTariff] = "hour"
			s.Step = CoworkingConfirm
			return botengine.InternalContinue, nil

		case "cw:tariff:half_day":
			ensureData(s)
			s.Data[dataTariff] = "half_day"
			s.Step = CoworkingConfirm
			return botengine.InternalContinue, nil

		case "cw:tariff:day":
			ensureData(s)
			s.Data[dataTariff] = "day"
			s.Step = CoworkingConfirm
			return botengine.InternalContinue, nil

		case "cw:tariff:week":
			ensureData(s)
			s.Data[dataTariff] = "week"
			s.Step = CoworkingConfirm
			return botengine.InternalContinue, nil

		case "cw:tariff:month":
			ensureData(s)
			s.Data[dataTariff] = "month"
			s.Step = CoworkingConfirm
			return botengine.InternalContinue, nil
		case "cw:profile":
			s.Step = CoworkingProfile
			return botengine.InternalContinue, nil

		case "cw:tariff":
			s.Step = CoworkingTariff
			return botengine.InternalContinue, nil
		
		case "cw:home":
			s.Step = CoworkingHome
			return botengine.InternalContinue, nil
		}
	}

	authStatus := strData(s, dataAuthStatus)
	if authStatus == string(types.RegistrationApproved) {
		text := "✅ Вы авторизованы в LAN Coworking.\n\nВыберите действие:"

		kb := ui.Inline(
			ui.Row(
				ui.Cb("👤 Мой профиль", "cw:profile"),
				ui.Cb("🎫 Мой тариф", "cw:tariff"),
			),
			ui.Row(
				ui.Cb("ℹ️ О коворкинге", "/about"),
			),
		)

		if err := ui.SendHTML(d.Bot, s.ChatID, text, kb); err != nil {
			return CoworkingHome, err
		}
		return CoworkingHome, nil
	}

	text := "👋 Добро пожаловать в LAN Coworking.\n\n" +
		"Чтобы сделать вход быстрым, выберите один из вариантов ниже."

	kb := ui.Inline(
		ui.Row(
			ui.Cb("🆕 Я впервые", "cw:new"),
			ui.Cb("♻️ Я уже был у вас", "cw:returning"),
		),
	)

	if err := ui.SendHTML(d.Bot, s.ChatID, text, kb); err != nil {
		return CoworkingHome, err
	}

	return CoworkingHome, nil
}

func coworkingReturningPhone(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if d.Svcs.CoworkingRegistrations == nil {
		return CoworkingHome, ui.SendText(d.Bot, s.ChatID, "Сервис регистраций не подключён.")
	}

	if ev.HasContact && ev.ContactPhone != "" {
		reg, err := d.Svcs.CoworkingRegistrations.GetLatestApprovedByPhone(ctx, ev.ContactPhone)
		if err != nil {
			return CoworkingReturningPhone, err
		}

		if reg == nil {
			text := "Я не нашёл ваш профиль.\n\nПожалуйста, пройдите быструю регистрацию."
			if err := ui.SendHTML(d.Bot, s.ChatID, text); err != nil {
				return CoworkingReturningPhone, err
			}
			s.Step = CoworkingNewName
			return botengine.InternalContinue, nil
		}

		err = d.Svcs.CoworkingRegistrations.Create(ctx, types.CoworkingRegistration{
			ChatID:           s.ChatID,
			TelegramUserID:   ev.FromUserID,
			TelegramUsername: ev.FromUserName,
			FullName:         reg.FullName,
			Phone:            ev.ContactPhone,
			TariffCode:       reg.TariffCode,
			RequestType:      types.RequestTypeRelink,
			Status:           types.RegistrationPending,
		})
		if err != nil {
			return CoworkingReturningPhone, err
		}

		ensureData(s)
		s.Data[dataFullName] = reg.FullName
		s.Data[dataPhone] = ev.ContactPhone
		s.Data[dataTariff] = reg.TariffCode
		s.Data[dataReqType] = string(types.RequestTypeRelink)
		s.Data[dataAuthStatus] = string(types.RegistrationPending)

		if err := sendAdminRelinkRequest(ctx, d, s, ev, reg.FullName, ev.ContactPhone, reg.TariffCode); err != nil {
			return CoworkingReturningPhone, err
		}

		text := "🔄 Мы нашли ваш профиль.\n\n" +
			"Запрос на повторную авторизацию отправлен администратору."
		if err := ui.SendHTML(d.Bot, s.ChatID, text, ui.RemoveKeyboard()); err != nil {
			return CoworkingReturningPhone, err
		}

		s.Step = CoworkingPending
		return CoworkingPending, nil
	}

	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("📱 Поделиться номером"),
		),
	)
	kb.ResizeKeyboard = true
	kb.OneTimeKeyboard = true

	msg := tgbotapi.NewMessage(s.ChatID, "Пожалуйста, отправьте номер телефона, который использовали при регистрации.")
	msg.ReplyMarkup = kb
	if _, err := d.Bot.Send(msg); err != nil {
		return CoworkingReturningPhone, err
	}

	return CoworkingReturningPhone, nil
}

func coworkingNewName(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if ev.Kind == botengine.EventText && strings.TrimSpace(ev.Text) != "" {
		ensureData(s)
		s.Data[dataFullName] = strings.TrimSpace(ev.Text)
		s.Step = CoworkingNewPhone
		return botengine.InternalContinue, nil
	}

	text := "Как вас зовут?\n\nОтправьте имя и фамилию одним сообщением."
	if err := ui.SendHTML(d.Bot, s.ChatID, text); err != nil {
		return CoworkingNewName, err
	}

	return CoworkingNewName, nil
}

func coworkingNewPhone(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if ev.HasContact && ev.ContactPhone != "" {
		ensureData(s)
		s.Data[dataPhone] = ev.ContactPhone
		s.Step = CoworkingNewTariff
		return botengine.InternalContinue, nil
	}

	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("📱 Поделиться номером"),
		),
	)
	kb.ResizeKeyboard = true
	kb.OneTimeKeyboard = true

	msg := tgbotapi.NewMessage(s.ChatID, "Поделитесь, пожалуйста, номером телефона.")
	msg.ReplyMarkup = kb
	if _, err := d.Bot.Send(msg); err != nil {
		return CoworkingNewPhone, err
	}

	return CoworkingNewPhone, nil
}

func coworkingNewTariff(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	text := "Выберите тариф:"
	kb := ui.Inline(
		ui.Row(
			ui.Cb("1 час", "cw:tariff:hour"),
			ui.Cb("Полдня", "cw:tariff:half_day"),
		),
		ui.Row(
			ui.Cb("1 день", "cw:tariff:day"),
			ui.Cb("1 неделя", "cw:tariff:week"),
		),
		ui.Row(
			ui.Cb("1 месяц", "cw:tariff:month"),
		),
	)

	if err := ui.SendHTML(d.Bot, s.ChatID, text, kb); err != nil {
		return CoworkingNewTariff, err
	}

	return CoworkingNewTariff, nil
}

func coworkingConfirm(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if d.Svcs.CoworkingRegistrations == nil {
		return CoworkingHome, ui.SendText(d.Bot, s.ChatID, "Сервис регистраций не подключён.")
	}

	if ev.Kind == botengine.EventCallback && ev.CallbackData == "cwconfirm:send" {
		ensureData(s)

		reg := types.CoworkingRegistration{
			ChatID:           s.ChatID,
			TelegramUserID:   ev.FromUserID,
			TelegramUsername: ev.FromUserName,
			FullName:         strData(s, dataFullName),
			Phone:            strData(s, dataPhone),
			TariffCode:       strData(s, dataTariff),
			RequestType:      types.RequestTypeNew,
			Status:           types.RegistrationPending,
		}

		if err := d.Svcs.CoworkingRegistrations.Create(ctx, reg); err != nil {
			return CoworkingConfirm, err
		}

		s.Data[dataAuthStatus] = string(types.RegistrationPending)
		s.Data[dataReqType] = string(types.RequestTypeNew)

		if err := sendAdminApprovalRequest(ctx, d, s, ev); err != nil {
			return CoworkingConfirm, err
		}

		s.Step = CoworkingPending
		return botengine.InternalContinue, nil
	}

	text := fmt.Sprintf(
		"Проверьте заявку:\n\n"+
			"👤 %s\n"+
			"📱 %s\n"+
			"🎫 %s\n\n"+
			"После отправки администратор подтвердит доступ.",
		strData(s, dataFullName),
		strData(s, dataPhone),
		renderTariffLabel(strData(s, dataTariff)),
	)

	kb := ui.Inline(
		ui.Row(
			ui.Cb("✅ Отправить", "cwconfirm:send"),
			ui.Cb("✏️ Изменить тариф", "cw:back:tariff"),
		),
	)

	if err := ui.SendHTML(d.Bot, s.ChatID, text, kb); err != nil {
		return CoworkingConfirm, err
	}

	return CoworkingConfirm, nil
}

func coworkingPending(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	status := strData(s, dataAuthStatus)

	if d.Svcs.CoworkingRegistrations != nil {
		reg, err := d.Svcs.CoworkingRegistrations.GetPendingByChatID(ctx, s.ChatID)
		if err != nil {
			return CoworkingPending, err
		}
		if reg != nil {
			status = string(reg.Status)
		}
	}

	switch status {
	case string(types.RegistrationApproved):
		text := "✅ Ваша заявка подтверждена.\n\nТеперь вам доступен функционал гостя."
		if err := ui.SendHTML(d.Bot, s.ChatID, text, ui.RemoveKeyboard()); err != nil {
			return CoworkingPending, err
		}
		s.Step = CoworkingHome
		return botengine.InternalContinue, nil

	case string(types.RegistrationRejected):
		text := "❌ Заявка пока не подтверждена.\n\nПожалуйста, подойдите к администратору на ресепшене."
		if err := ui.SendHTML(d.Bot, s.ChatID, text, ui.RemoveKeyboard()); err != nil {
			return CoworkingPending, err
		}
		s.Step = CoworkingHome
		return botengine.InternalContinue, nil
	}

	text := "⏳ Заявка отправлена администратору.\n\nОжидайте подтверждения."
	if err := ui.SendHTML(d.Bot, s.ChatID, text, ui.RemoveKeyboard()); err != nil {
		return CoworkingPending, err
	}

	return CoworkingPending, nil
}

func coworkingAdminAction(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if ev.Kind != botengine.EventCallback || !strings.HasPrefix(ev.CallbackData, "cwa:") {
		return CoworkingAdminAction, nil
	}

	parts := strings.Split(ev.CallbackData, ":")
	if len(parts) < 3 {
		_ = ui.Toast(d.Bot, ev.CallbackQueryID, "Некорректное действие")
		return CoworkingAdminAction, nil
	}

	action := parts[1]
	targetChatID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		_ = ui.Toast(d.Bot, ev.CallbackQueryID, "Некорректный chat id")
		return CoworkingAdminAction, nil
	}

	if d.State == nil {
		_ = ui.Alert(d.Bot, ev.CallbackQueryID, "State manager is not configured")
		return CoworkingAdminAction, nil
	}

	if d.Svcs.CoworkingRegistrations == nil {
		_ = ui.Alert(d.Bot, ev.CallbackQueryID, "Registrations service is not configured")
		return CoworkingAdminAction, nil
	}

	target := d.State.Get(targetChatID)
	if target == nil {
		_ = ui.Alert(d.Bot, ev.CallbackQueryID, "Пользователь не найден")
		return CoworkingAdminAction, nil
	}

	ensureData(target)

	switch action {
	case "approve":
		if err := d.Svcs.CoworkingRegistrations.UpdateStatusByChatID(
			ctx,
			targetChatID,
			types.RegistrationApproved,
			ev.FromUserID,
			"",
		); err != nil {
			return CoworkingAdminAction, err
		}

		target.Data[dataAuthStatus] = string(types.RegistrationApproved)
		d.State.Set(target.ChatID, target)

		if d.Svcs.Guests != nil {
			fullName := strData(target, dataFullName)
			firstName, lastName := splitName(fullName)

			guest := types.Guest{
				Telegram:  fmt.Sprintf("id:%d", target.UserID),
				FirstName: firstName,
				LastName:  lastName,
			}

			if err := d.Svcs.Guests.AddGuest(ctx, d.Cfg.GuestsReadRange, guest); err != nil {
				log.Printf("[coworking.admin] cannot save guest: %v", err)
			}
		}

		_, _ = d.Bot.Send(tgbotapi.NewMessage(
			target.ChatID,
			"✅ Администратор подтвердил вашу заявку.\nТеперь вам доступен функционал гостя.",
		))
		_ = ui.Toast(d.Bot, ev.CallbackQueryID, "Approved")

	case "reject":
		if err := d.Svcs.CoworkingRegistrations.UpdateStatusByChatID(
			ctx,
			targetChatID,
			types.RegistrationRejected,
			ev.FromUserID,
			"",
		); err != nil {
			return CoworkingAdminAction, err
		}

		target.Data[dataAuthStatus] = string(types.RegistrationRejected)
		d.State.Set(target.ChatID, target)

		_, _ = d.Bot.Send(tgbotapi.NewMessage(
			target.ChatID,
			"❌ Заявка не подтверждена.\nПожалуйста, обратитесь к администратору.",
		))
		_ = ui.Toast(d.Bot, ev.CallbackQueryID, "Rejected")

	default:
		_ = ui.Toast(d.Bot, ev.CallbackQueryID, "Неизвестное действие")
	}

	return CoworkingAdminAction, nil
}

func sendAdminApprovalRequest(ctx context.Context, d botengine.Deps, s *types.Session, ev botengine.Event) error {
	_ = ctx

	if d.Cfg.AdminChatId == 0 {
		return fmt.Errorf("ADMIN_CHAT_ID is empty")
	}

	text := fmt.Sprintf(
		"🆕 Новая заявка на доступ в коворкинг\n\n"+
			"ChatID: %d\n"+
			"UserID: %d\n"+
			"Username: @%s\n"+
			"Имя: %s\n"+
			"Телефон: %s\n"+
			"Тариф: %s\n"+
			"Время: %s",
		s.ChatID,
		ev.FromUserID,
		ev.FromUserName,
		strData(s, dataFullName),
		strData(s, dataPhone),
		renderTariffLabel(strData(s, dataTariff)),
		time.Now().Format(time.RFC3339),
	)

	msg := tgbotapi.NewMessage(d.Cfg.AdminChatId, text)
	msg.ReplyMarkup = ui.Inline(
		ui.Row(
			ui.Cb("✅ Approve", fmt.Sprintf("cwa:approve:%d", s.ChatID)),
			ui.Cb("❌ Reject", fmt.Sprintf("cwa:reject:%d", s.ChatID)),
		),
	)

	_, err := d.Bot.Send(msg)
	return err
}

func sendAdminRelinkRequest(
	ctx context.Context,
	d botengine.Deps,
	s *types.Session,
	ev botengine.Event,
	fullName string,
	phone string,
	tariffCode string,
) error {
	_ = ctx

	if d.Cfg.AdminChatId == 0 {
		return fmt.Errorf("ADMIN_CHAT_ID is empty")
	}

	text := fmt.Sprintf(
		"🔄 Повторная авторизация гостя\n\n"+
			"ChatID: %d\n"+
			"UserID: %d\n"+
			"Username: @%s\n"+
			"Найденный профиль: %s\n"+
			"Телефон: %s\n"+
			"Тариф: %s\n"+
			"Время: %s",
		s.ChatID,
		ev.FromUserID,
		ev.FromUserName,
		fullName,
		phone,
		renderTariffLabel(tariffCode),
		time.Now().Format(time.RFC3339),
	)

	msg := tgbotapi.NewMessage(d.Cfg.AdminChatId, text)
	msg.ReplyMarkup = ui.Inline(
		ui.Row(
			ui.Cb("✅ Approve", fmt.Sprintf("cwa:approve:%d", s.ChatID)),
			ui.Cb("❌ Reject", fmt.Sprintf("cwa:reject:%d", s.ChatID)),
		),
	)

	_, err := d.Bot.Send(msg)
	return err
}

func coworkingProfile(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if d.Svcs.CoworkingRegistrations == nil {
		return CoworkingHome, ui.SendText(d.Bot, s.ChatID, "Сервис регистраций не подключён.")
	}

	reg, err := d.Svcs.CoworkingRegistrations.GetLatestApprovedByChatID(ctx, s.ChatID)
	if err != nil {
		return CoworkingProfile, err
	}

	if reg == nil {
		text := "Не удалось найти подтверждённый профиль."
		if err := ui.SendHTML(d.Bot, s.ChatID, text); err != nil {
			return CoworkingProfile, err
		}
		s.Step = CoworkingHome
		return botengine.InternalContinue, nil
	}

	text := fmt.Sprintf(
		"👤 <b>Мой профиль</b>\n\n"+
			"Имя: %s\n"+
			"Телефон: %s\n"+
			"Telegram: @%s\n"+
			"Статус: подтверждён\n"+
			"Тариф: %s",
		reg.FullName,
		reg.Phone,
		reg.TelegramUsername,
		renderTariffLabel(reg.TariffCode),
	)

	kb := ui.Inline(
		ui.Row(
			ui.Cb("⬅️ Назад", "cw:home"),
		),
	)

	if err := ui.SendHTML(d.Bot, s.ChatID, text, kb); err != nil {
		return CoworkingProfile, err
	}

	return CoworkingProfile, nil
}

func coworkingTariff(ctx context.Context, ev botengine.Event, d botengine.Deps, s *types.Session) (types.Step, error) {
	if d.Svcs.CoworkingRegistrations == nil {
		return CoworkingHome, ui.SendText(d.Bot, s.ChatID, "Сервис регистраций не подключён.")
	}

	reg, err := d.Svcs.CoworkingRegistrations.GetLatestApprovedByChatID(ctx, s.ChatID)
	if err != nil {
		return CoworkingTariff, err
	}

	if reg == nil {
		if err := ui.SendHTML(d.Bot, s.ChatID, "Не удалось найти ваш тариф."); err != nil {
			return CoworkingTariff, err
		}
		s.Step = CoworkingHome
		return botengine.InternalContinue, nil
	}

	text := fmt.Sprintf(
		"🎫 <b>Мой тариф</b>\n\nТекущий тариф: %s",
		renderTariffLabel(reg.TariffCode),
	)

	kb := ui.Inline(
		ui.Row(
			ui.Cb("⬅️ Назад", "cw:home"),
		),
	)

	if err := ui.SendHTML(d.Bot, s.ChatID, text, kb); err != nil {
		return CoworkingTariff, err
	}

	return CoworkingTariff, nil
}

func ensureData(s *types.Session) {
	if s.Data == nil {
		s.Data = map[string]interface{}{}
	}
}

func strData(s *types.Session, key string) string {
	if s == nil || s.Data == nil {
		return ""
	}
	v, ok := s.Data[key]
	if !ok || v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func splitName(full string) (string, string) {
	full = strings.TrimSpace(full)
	if full == "" {
		return "", ""
	}

	parts := strings.Fields(full)
	if len(parts) == 1 {
		return parts[0], ""
	}

	return parts[0], strings.Join(parts[1:], " ")
}

func renderTariffLabel(code string) string {
	switch code {
	case "hour":
		return "1 час"
	case "half_day":
		return "Полдня"
	case "day":
		return "1 день"
	case "week":
		return "1 неделя"
	case "month":
		return "1 месяц"
	default:
		return code
	}
}