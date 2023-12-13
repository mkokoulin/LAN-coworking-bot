package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GuestsSheetService struct {
	spreadsheetId string
	readRange string
	srv *sheets.Service
}

type Guest struct {
	Telegram string `json:"telegram" mapstructure:"telegram"`
	FirstName string `json:"firstName" mapstructure:"firstName"`
	LastName string `json:"lastName" mapstructure:"lastName"`
	Datetime string `json:"datetime" mapstructure:"datetime"`
	Commands []string `json:"commands" mapstructure:"commands"`
}

func NewGuestSheets(ctx context.Context, googleClient *http.Client, spreadsheetId, readRange string) (*GuestsSheetService, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &GuestsSheetService{
		spreadsheetId,
		readRange,
		srv,
	}, nil
}

func (ESS *GuestsSheetService) GetGuests(ctx context.Context) ([]Guest, error) {
	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, "guests!A1:Z10000").Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return nil, fmt.Errorf("%v", err)
	}

	guests := []Guest{}

	header := []string{}

	for i, val := range res.Values {
		if i == 0 {
			for _, v := range val {
				header = append(header, v.(string))
			}
		} else {
			rawCoworker := map[string]string{}

			for i, v := range header {
				if len(val) - 1 < i {
					rawCoworker[v] = ""
				} else {
					rawCoworker[v] = val[i].(string)
				}
			}

			guest := Guest{}

			err := mapstructure.Decode(rawCoworker, &guest)
			if err != nil {
				return guests, err
			}

			guests = append(guests, guest)
		}
	}

	return guests, nil
}


func (ESS *GuestsSheetService) GetGuest(ctx context.Context, guest string) (Guest, error) {
	g := Guest{}

	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, ESS.readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return g, fmt.Errorf("%v", err)
	}

	guests, err := ESS.GetGuests(ctx)
	if err != nil {
		return g, fmt.Errorf("%v", err)
	}

	for _, v := range guests {
		if v.Telegram == guest {
			g = v
			break
		}
	}

	return g, nil
}

func (ESS *GuestsSheetService) CreateGuest(ctx context.Context, readRange string, guest Guest) (Guest, error) {
	resGuest, err := ESS.GetGuest(ctx, guest.Telegram)
	if err != nil {
		return Guest{}, fmt.Errorf("%v", err)
	}

	newGuest := guest.FirstName + guest.LastName + guest.Telegram
	oldGuest := resGuest.FirstName + resGuest.LastName + resGuest.Telegram

	if newGuest == oldGuest {
		return resGuest, nil
	}

	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return Guest{}, fmt.Errorf("%v", err)
	}

	rowNumber := len(res.Values) + 2
	sheet := readRange[:strings.IndexByte(readRange, '!')]
	updateRowRange := fmt.Sprintf("%s!A%d:Z%d", sheet, rowNumber, rowNumber)

	now := time.Now()
	guest.Datetime = now.Format(time.RFC3339)

	row := &sheets.ValueRange{
		Values: [][]interface{}{{
			guest.FirstName,
			guest.LastName,
			guest.Telegram,
			guest.Datetime,
		}},
	}

	_, err = ESS.srv.Spreadsheets.Values.Update(ESS.spreadsheetId, updateRowRange, row).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return Guest{}, fmt.Errorf("%v", err)
	}

	return guest, nil
}