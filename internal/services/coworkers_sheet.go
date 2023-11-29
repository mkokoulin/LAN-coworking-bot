package services

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"github.com/mitchellh/mapstructure"
)

type CoworkersSheetService struct {
	spreadsheetId string
	readRange string
	srv *sheets.Service
}

type Coworker struct {
	Secret string `json:"secret" mapstructure:"secret"`
	Telegram string `json:"telegram" mapstructure:"telegram"`
}

func NewCoworkersSheets(ctx context.Context, googleClient *http.Client, spreadsheetId, readRange string) (*CoworkersSheetService, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &CoworkersSheetService{
		spreadsheetId,
		readRange,
		srv,
	}, nil
}

func (ESS *CoworkersSheetService) GetCoworkers(ctx context.Context) ([]Coworker, error) {
	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, "master!A1:Z1000").Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return nil, fmt.Errorf("%v", err)
	}

	coworkers := []Coworker{}

	header := []string{}

	for i, val := range res.Values {
		if i == 0 {
			for _, v := range val {
				header = append(header, v.(string))
			}
		} else {
			rawCoworker := map[string]string{}

			for i, v := range header {
				fmt.Println(len(val))

				if len(val) - 1 < i {
					rawCoworker[v] = ""
				} else {
					rawCoworker[v] = val[i].(string)
				}
			}

			coworker := Coworker{}

			err := mapstructure.Decode(rawCoworker, &coworker)
			if err != nil {
				return coworkers, err
			}

			coworkers = append(coworkers, coworker)
		}
	}

	return coworkers, nil
}

func (ESS *CoworkersSheetService) GetCoworker(ctx context.Context, telegram string) (Coworker, error) {
	coworker := Coworker{}

	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, ESS.readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return coworker, fmt.Errorf("%v", err)
	}

	coworkers, err := ESS.GetCoworkers(ctx)
	if err != nil {
		return coworker, fmt.Errorf("%v", err)
	}

	for _, v := range coworkers {
		if v.Telegram == telegram {
			coworker = v
			break
		}
	}

	return coworker, nil
}

func (ESS *CoworkersSheetService) GetUnusedSecrets(ctx context.Context) ([]string, error) {
	secrets := []string {}

	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, ESS.readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return secrets, fmt.Errorf("%v", err)
	}

	coworkers, err := ESS.GetCoworkers(ctx)
	if err != nil {
		return secrets, fmt.Errorf("%v", err)
	}

	for _, v := range coworkers {
		if v.Telegram == "" {
			secrets = append(secrets, v.Secret)
		}
	}

	return secrets, nil
}

func (ESS *CoworkersSheetService) UpdateCoworker(ctx context.Context, coworker Coworker) error {
	var rowNumber int
	readRange := "master!2:1000" 

	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return fmt.Errorf("%v", err)
	}

	for i, v := range res.Values {
		if v[0] == coworker.Secret {
			rowNumber = i + 2
		}
	}

	updateRowRange := fmt.Sprintf("A%d:Z%d", rowNumber, rowNumber)

	row := &sheets.ValueRange{
		Values: [][]interface{}{{
			coworker.Secret,
			coworker.Telegram,
		}},
	}

	_, err = ESS.srv.Spreadsheets.Values.Update(ESS.spreadsheetId, updateRowRange, row).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}