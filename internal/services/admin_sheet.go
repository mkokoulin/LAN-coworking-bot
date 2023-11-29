package services

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type AdminSheetService struct {
	spreadsheetId string
	readRange string
	srv *sheets.Service
}

func NewAdminSheets(ctx context.Context, googleClient *http.Client, spreadsheetId, readRange string) (*AdminSheetService, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &AdminSheetService{
		spreadsheetId,
		readRange,
		srv,
	}, nil
}

func (ESS *AdminSheetService) GetSecrets(ctx context.Context) ([]string, error) {
	res, err := ESS.srv.Spreadsheets.Values.Get(ESS.spreadsheetId, ESS.readRange).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return nil, fmt.Errorf("%v", err)
	}

	secrets := []string{}

	for _, val := range res.Values {
		if len(val) == 1 {
			secrets = append(secrets, fmt.Sprintf("%v", val[0]))
		}
	}

	return secrets, nil
}