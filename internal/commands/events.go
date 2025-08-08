package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/locales"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

func stripHTML(input string) string {
	re := regexp.MustCompile(`<.*?>`)
	return strings.TrimSpace(re.ReplaceAllString(input, ""))
}

func parseEventDate(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02",  // ISO
		"02.01.2006",  // RU
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized date format: %s", s)
}

func EventsCommand(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, services types.Services, state *types.ChatStorage) error {
	p := locales.Printer(state.Language)
	chatID := update.Message.Chat.ID

	resp, err := http.Get("https://shark-app-wrcei.ondigitalocean.app/api/events")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var rawEvents []types.Event
	if err := json.Unmarshal(body, &rawEvents); err != nil {
		return err
	}

	var events []types.Event
	for _, e := range rawEvents {
		if !e.ShowForm {
			continue
		}
		if _, err := parseEventDate(e.Date); err != nil {
			continue
		}
		events = append(events, e)
	}

	sort.Slice(events, func(i, j int) bool {
		dateI, _ := parseEventDate(events[i].Date)
		dateJ, _ := parseEventDate(events[j].Date)
		return dateI.Before(dateJ)
	})

	if len(events) > 5 {
		events = events[:5]
	}

	// Вводное сообщение
	intro := tgbotapi.NewMessage(chatID, p.Sprintf("events_intro"))
	intro.ParseMode = "HTML"
	if _, err := bot.Send(intro); err != nil {
		return err
	}

	// Список мероприятий
	var sb strings.Builder
	for _, e := range events {
		description := stripHTML(e.Description)
		sb.WriteString(p.Sprintf("event_item", e.Date, description, e.Id))
	}

	state.CurrentCommand = ""

	list := tgbotapi.NewMessage(chatID, sb.String())
	list.ParseMode = "HTML"

	_, err = bot.Send(list)
	return err
}
