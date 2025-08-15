package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type botLogsSheetService struct {
	spreadsheetID string
	readRange     string
	srv           *sheets.Service
}

// Конструктор возвращает ИНТЕРФЕЙС, а не *интерфейс.
func NewBotLogsSheets(ctx context.Context, googleClient *http.Client, spreadsheetID, readRange string) (types.BotLogsSheetsService, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, fmt.Errorf("sheets.NewService: %w", err)
	}
	return &botLogsSheetService{
		spreadsheetID: spreadsheetID,
		readRange:     readRange,
		srv:           srv,
	}, nil
}

// Log — добавляет запись в конец диапазона.
func (s *botLogsSheetService) Log(ctx context.Context, readRange string, botLog types.BotLog) error {
	if strings.TrimSpace(readRange) == "" {
		readRange = s.readRange
	}
	if strings.TrimSpace(readRange) == "" {
		return fmt.Errorf("readRange is empty")
	}

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

	// Следующая строка: после заголовка (если он есть)
	rowNumber := len(res.Values) + 1
	if rowNumber < 2 {
		rowNumber = 2
	}
	targetRange := fmt.Sprintf("%s!A%d:Z%d", sheet, rowNumber, rowNumber)

	botLog.Datetime = time.Now().Format(time.RFC3339)

	values := &sheets.ValueRange{
		Values: [][]interface{}{{
			botLog.Telegram,
			botLog.Command,
			botLog.Datetime,
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
