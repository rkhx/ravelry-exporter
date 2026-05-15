package handlers

import (
	"context"
	"fmt"
	"github.com/rkhx/ravelry-exporter/internal/gapi"
	"github.com/rkhx/ravelry-exporter/internal/models"
	"github.com/rkhx/ravelry-exporter/internal/ravelry"
	"github.com/rkhx/ravelry-exporter/internal/ui"
	"github.com/rs/zerolog/log"
	"strings"
)

func Run(ctx context.Context, bundleUrl string) error {
	cfg, err := ui.ReadConfig()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	login := cfg.Credentials.Username
	password := cfg.Credentials.Password
	if login == "" || password == "" {
		return fmt.Errorf("RAVELRY_LOGIN or RAVELRY_PASSWORD is not set")
	}

	serviceAccountFileContent := cfg.ServiceAccount
	if serviceAccountFileContent == "" {
		return fmt.Errorf("SERVICE_ACCOUNT_FILE_CONTENT is not set")
	}

	spreadsheetID := cfg.SpreadsheetID
	if spreadsheetID == "" {
		log.Warn().Msg("SPREADSHEET_ID is not set, we'll create a new spreadsheet")
	}

	client := ravelry.NewRavelryClient("https://api.ravelry.com", login, password)

	username, err := client.Users.GetCurrentUsername(ctx)
	if err != nil {
		return fmt.Errorf("error fetching username: %w", err)
	}
	log.Info().Msg("Current username fetched successfully")

	bundles, err := client.Bundles.GetUserBundles(ctx, username)
	if err != nil || len(bundles) == 0 {
		return fmt.Errorf("error fetching bundles or no bundles found: %w", err)
	}
	log.Info().Msg("Bundles list fetched successfully")

	bundlePermalink := strings.TrimPrefix(bundleUrl, "https://www.ravelry.com/bundles/")

	var foundBundle *models.Bundle
	for _, b := range bundles {
		bundle, err := client.Bundles.GetBundleContent(ctx, username, b.ID)
		if err != nil {
			return fmt.Errorf("error fetching bundle content: %w", err)
		}
		log.Info().Int("bundle_id", b.ID).Msg("Bundle content fetched successfully")

		if bundle.Permalink == bundlePermalink {
			log.Info().Msg("Bundle was found")
			foundBundle = bundle
			break
		}
	}

	if foundBundle == nil {
		return fmt.Errorf("bundle not found")
	}

	var data []gapi.RowData

	for _, item := range foundBundle.BundledItems {
		aa, err := client.Bundles.GetBundleItem(ctx, item.ID)
		if err != nil {
			return fmt.Errorf("error fetching bundle item: %w", err)
		}

		patternInfo, err := client.Patterns.GetPattern(ctx, aa.Item[0].ID)
		if err != nil {
			return fmt.Errorf("error fetching pattern: %w", err)
		}

		data = append(data, gapi.NewRowData(patternInfo))
	}

	columns := []interface{}{"Pattern Name", "Designer", "Gauge, needle size", "Sizes", "Recommended yarn", "Attributes"}

	err = gapi.Spreadsheet(ctx, serviceAccountFileContent, spreadsheetID, data, columns)
	if err != nil {
		return fmt.Errorf("error fetching spreadsheet: %w", err)
	}

	return nil
}
