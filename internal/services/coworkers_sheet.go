package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type coworkersSheetService struct {
	spreadsheetID string
	readRange     string
	srv           *sheets.Service
}

type Coworker struct {
	Secret   string `json:"secret" mapstructure:"secret"`
	Telegram string `json:"telegram" mapstructure:"telegram"`
}

// Конструктор возвращает ИНТЕРФЕЙС, не *интерфейс.
func NewCoworkersSheets(ctx context.Context, googleClient *http.Client, spreadsheetID, readRange string) (types.CoworkersSheetsService, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, fmt.Errorf("sheets.NewService: %w", err)
	}
	return &coworkersSheetService{
		spreadsheetID: spreadsheetID,
		readRange:     readRange,
		srv:           srv,
	}, nil
}

// ValidateSecret — true, если код есть среди «свободных» секретов (строка с пустым Telegram).
func (s *coworkersSheetService) ValidateSecret(ctx context.Context, code string) (bool, error) {
	if code == "" {
		return false, errors.New("empty code")
	}
	coworkers, err := s.GetCoworkers(ctx)
	if err != nil {
		return false, err
	}
	for _, c := range coworkers {
		if c.Secret == code && c.Telegram == "" {
			return true, nil
		}
	}
	return false, nil
}

// GetUnusedSecrets — список всех секретов без закреплённого Telegram.
func (s *coworkersSheetService) GetUnusedSecrets(ctx context.Context) ([]string, error) {
	coworkers, err := s.GetCoworkers(ctx)
	if err != nil {
		return nil, err
	}
	secrets := make([]string, 0, len(coworkers))
	for _, c := range coworkers {
		if c.Telegram == "" {
			secrets = append(secrets, c.Secret)
		}
	}
	return secrets, nil
}

// GetCoworkers — читает лист master!A1:Z1000 и мапит по заголовкам.
func (s *coworkersSheetService) GetCoworkers(ctx context.Context) ([]Coworker, error) {
	const rng = "master!A1:Z1000"
	res, err := s.srv.Spreadsheets.Values.Get(s.spreadsheetID, rng).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("sheets get %s: %w", rng, err)
	}
	if res.HTTPStatusCode != 200 {
		return nil, fmt.Errorf("sheets get %s: http %d", rng, res.HTTPStatusCode)
	}

	var (
		coworkers []Coworker
		header    []string
	)

	for i, row := range res.Values {
		// строка заголовков
		if i == 0 {
			header = header[:0]
			for _, v := range row {
				if s, ok := v.(string); ok {
					header = append(header, s)
				} else {
					header = append(header, fmt.Sprint(v))
				}
			}
			continue
		}

		// данные
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

		var cw Coworker
		if err := mapstructure.Decode(raw, &cw); err != nil {
			return coworkers, fmt.Errorf("decode coworker: %w", err)
		}
		coworkers = append(coworkers, cw)
	}

	return coworkers, nil
}

// GetCoworker — ищет по username в кэше таблицы (простая фильтрация после GetCoworkers).
func (s *coworkersSheetService) GetCoworker(ctx context.Context, telegram string) (Coworker, error) {
	var zero Coworker
	coworkers, err := s.GetCoworkers(ctx)
	if err != nil {
		return zero, err
	}
	for _, c := range coworkers {
		if c.Telegram == telegram {
			return c, nil
		}
	}
	return zero, fmt.Errorf("coworker not found: %s", telegram)
}

// UpdateCoworker — обновляет строку (по совпадению первого столбца = secret).
func (s *coworkersSheetService) UpdateCoworker(ctx context.Context, coworker Coworker) error {
	readRange := "master!2:1000"

	res, err := s.srv.Spreadsheets.Values.Get(s.spreadsheetID, readRange).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("sheets get %s: %w", readRange, err)
	}
	if res.HTTPStatusCode != 200 {
		return fmt.Errorf("sheets get %s: http %d", readRange, res.HTTPStatusCode)
	}

	rowNumber := 0
	for i, row := range res.Values {
		if len(row) > 0 && fmt.Sprint(row[0]) == coworker.Secret {
			rowNumber = i + 2 // +2, потому что диапазон начинается со 2-й строки
			break
		}
	}
	if rowNumber == 0 {
		return fmt.Errorf("row for secret not found: %s", coworker.Secret)
	}

	updateRange := fmt.Sprintf("A%d:Z%d", rowNumber, rowNumber)
	values := &sheets.ValueRange{
		Values: [][]interface{}{{
			coworker.Secret,
			coworker.Telegram,
		}},
	}

	_, err = s.srv.Spreadsheets.Values.
		Update(s.spreadsheetID, updateRange, values).
		ValueInputOption("USER_ENTERED").
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("sheets update %s: %w", updateRange, err)
	}
	return nil
}
