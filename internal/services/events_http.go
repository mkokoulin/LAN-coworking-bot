package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

type EventsService struct {
	client *http.Client
	url    string
}

func NewEventsService(client *http.Client, url string) *EventsService {
	return &EventsService{
		client: client,
		url:    url,
	}
}

func (s *EventsService) ListUpcoming(ctx context.Context) ([]types.Event, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, s.url, nil)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("events http: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("events read: %w", err)
	}

	var all []types.Event
	if err := json.Unmarshal(b, &all); err != nil {
		return nil, fmt.Errorf("events json: %w", err)
	}

	out := make([]types.Event, 0, len(all))
	for _, e := range all {
		if !e.ShowForm {
			continue
		}
		if _, err := parseDate(e.Date); err == nil {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ti, _ := parseDate(out[i].Date)
		tj, _ := parseDate(out[j].Date)
		return ti.Before(tj)
	})
	return out, nil
}

func parseDate(s string) (time.Time, error) {
	for _, f := range []string{"2006-01-02", "02.01.2006"} {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("bad date: %s", s)
}
