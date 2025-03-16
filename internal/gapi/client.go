package gapi

import (
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// newGoogleServices creates clients for Google Sheets and Google Drive API
func newGoogleServices(serviceAccountFileContent string, ctx context.Context) (sheetsSrv *sheets.Service, driveSrv *drive.Service, err error) {
	cred := option.WithCredentialsJSON([]byte(serviceAccountFileContent))

	sheetsSrv, err = sheets.NewService(ctx, cred)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create Sheets client: %w", err)
	}

	driveSrv, err = drive.NewService(ctx, cred)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create Drive client: %w", err)
	}

	return sheetsSrv, driveSrv, nil
}
