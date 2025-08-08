package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type BotLogsSheetService struct {
	spreadsheetId string
	readRange     string
	srv           *sheets.Service
}

type BotLog struct {
	Telegram string `json:"telegram" mapstructure:"telegram"`
	Command  string `json:"command" mapstructure:"command"`
	Datetime string `json:"datetime" mapstructure:"datetime"`
}

func NewBotLogsSheets(ctx context.Context, googleClient *http.Client, spreadsheetId, readRange string) (*BotLogsSheetService, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &BotLogsSheetService{
		spreadsheetId,
		readRange,
		srv,
	}, nil
}

func (ESS *BotLogsSheetService) Log(ctx context.Context, readRange string, botLog BotLog) error {
	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return fmt.Errorf("%v", err)
	}

	rowNumber := len(res.Values) + 2
	sheet := readRange[:strings.IndexByte(readRange, '!')]
	updateRowRange := fmt.Sprintf("%s!A%d:Z%d", sheet, rowNumber, rowNumber)

	now := time.Now()
	botLog.Datetime = now.Format(time.RFC3339)

	row := &sheets.ValueRange{
		Values: [][]interface{}{{
			botLog.Telegram,
			botLog.Command,
			botLog.Datetime,
		}},
	}

	_, err = ESS.srv.Spreadsheets.Values.Update(ESS.spreadsheetId, updateRowRange, row).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}