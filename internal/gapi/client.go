package gapi

import (
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"os"
)

// newGoogleServices creates clients for Google Sheets and Google Drive API
func newGoogleServices(ctx context.Context) (*sheets.Service, *drive.Service, error) {
	content := os.Getenv("SERVICE_ACCOUNT_FILE_CONTENT")
	cred := option.WithCredentialsJSON([]byte(content))

	sheetsSrv, err := sheets.NewService(ctx, cred)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create Sheets client: %w", err)
	}

	driveSrv, err := drive.NewService(ctx, cred)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create Drive client: %w", err)
	}

	return sheetsSrv, driveSrv, nil
}
