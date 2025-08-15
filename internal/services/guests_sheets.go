package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type guestsSheetService struct {
	spreadsheetID string
	readRange     string
	srv           *sheets.Service
}

// Конструктор возвращает ИНТЕРФЕЙС (не *интерфейс).
func NewGuestSheets(ctx context.Context, googleClient *http.Client, spreadsheetID, readRange string) (types.GuestSheetsService, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, fmt.Errorf("sheets.NewService: %w", err)
	}
	return &guestsSheetService{
		spreadsheetID: spreadsheetID,
		readRange:     readRange,
		srv:           srv,
	}, nil
}

// GetGuests — читает гостей из листа guests!A1:Z10000, мапит по заголовкам.
func (s *guestsSheetService) GetGuests(ctx context.Context) ([]types.Guest, error) {
	const rng = "guests!A1:Z10000"

	res, err := s.srv.Spreadsheets.Values.Get(s.spreadsheetID, rng).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("sheets get %s: %w", rng, err)
	}
	if res.HTTPStatusCode != 200 {
		return nil, fmt.Errorf("sheets get %s: http %d", rng, res.HTTPStatusCode)
	}

	var (
		guests []types.Guest
		header []string
	)

	for i, row := range res.Values {
		if i == 0 {
			header = header[:0]
			for _, v := range row {
				if str, ok := v.(string); ok {
					header = append(header, str)
				} else {
					header = append(header, fmt.Sprint(v))
				}
			}
			continue
		}

		raw := map[string]string{}
		for idx, key := range header {
			if idx < len(row) {
				if str, ok := row[idx].(string); ok {
					raw[key] = str
				} else {
					raw[key] = fmt.Sprint(row[idx])
				}
			} else {
				raw[key] = ""
			}
		}

		var g types.Guest
		if err := mapstructure.Decode(raw, &g); err != nil {
			return guests, fmt.Errorf("decode guest: %w", err)
		}
		guests = append(guests, g)
	}

	return guests, nil
}

// GetGuest — ищет гостя по Telegram.
func (s *guestsSheetService) GetGuest(ctx context.Context, telegram string) (types.Guest, error) {
	var zero types.Guest

	guests, err := s.GetGuests(ctx)
	if err != nil {
		return zero, err
	}
	for _, g := range guests {
	 if g.Telegram == telegram {
			return g, nil
		}
	}
	return zero, fmt.Errorf("guest not found: %s", telegram)
}

// CreateGuest — создаёт запись, если её ещё нет.
// readRange можно передать явно; если пусто — используется s.readRange.
func (s *guestsSheetService) CreateGuest(ctx context.Context, readRange string, guest types.Guest) error {
	return s.addGuest(ctx, readRange, guest)
}

// AddGuest — алиас к CreateGuest для совместимости с интерфейсом.
func (s *guestsSheetService) AddGuest(ctx context.Context, readRange string, guest types.Guest) error {
	return s.addGuest(ctx, readRange, guest)
}

// внутренняя реализация добавления гостя
func (s *guestsSheetService) addGuest(ctx context.Context, readRange string, guest types.Guest) error {
	if strings.TrimSpace(readRange) == "" {
		readRange = s.readRange
	}
	if strings.TrimSpace(readRange) == "" {
		readRange = "guests!A1:Z10000"
	}

	// Не дублируем запись
	if existing, err := s.GetGuest(ctx, guest.Telegram); err == nil {
		if existing.FirstName == guest.FirstName &&
			existing.LastName == guest.LastName &&
			existing.Telegram == guest.Telegram {
			return nil
		}
	}

	// Текущее число строк
	res, err := s.srv.Spreadsheets.Values.Get(s.spreadsheetID, readRange).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("sheets get %s: %w", readRange, err)
	}
	if res.HTTPStatusCode != 200 {
		return fmt.Errorf("sheets get %s: http %d", readRange, res.HTTPStatusCode)
	}

	// Имя листа до '!'
	sheet := readRange
	if i := strings.IndexByte(readRange, '!'); i > 0 {
		sheet = readRange[:i]
	}

	// Следующая строка (минимум 2 — после шапки)
	rowNumber := len(res.Values) + 1
	if rowNumber < 2 {
		rowNumber = 2
	}
	targetRange := fmt.Sprintf("%s!A%d:Z%d", sheet, rowNumber, rowNumber)

	guest.Datetime = time.Now().Format(time.RFC3339)

	values := &sheets.ValueRange{
		Values: [][]interface{}{{
			guest.FirstName,
			guest.LastName,
			guest.Telegram,
			guest.Datetime,
		}},
	}

	_, err = s.srv.Spreadsheets.Values.
		Update(s.spreadsheetID, targetRange, values).
		ValueInputOption("USER_ENTERED").
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("sheets update %s: %w", targetRange, err)
	}
	return nil
}
