// Main pack
package main

import (
	"context"
	"fmt"
	"github.com/rkhx/ravelry-exporter/internal/gapi"
	"github.com/rkhx/ravelry-exporter/internal/ravelry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339, NoColor: false})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	defer func() {
		if r := recover(); r != nil {
			pc, file, line, ok := runtime.Caller(3)
			fnName := "unknown"
			if ok {
				fnName = runtime.FuncForPC(pc).Name()
			}

			log.Error().
				Str("file", file).
				Int("line", line).
				Str("function", fnName).
				Msgf("PANIC: %v", r)
		}
	}()

	start := time.Now()
	if err := run(ctx); err != nil {
		log.Error().Msgf("%s", err.Error())
	}
	log.Info().Dur("uptime", time.Since(start)).Msg("Application exited successfully")
}

func run(ctx context.Context) error {
	login, password := os.Getenv("RAVELRY_LOGIN"), os.Getenv("RAVELRY_PASSWORD")
	if login == "" || password == "" {
		return fmt.Errorf("RAVELRY_LOGIN or RAVELRY_PASSWORD is not set")
	}

	serviceAccountFileContent := os.Getenv("SERVICE_ACCOUNT_FILE_CONTENT")
	if serviceAccountFileContent == "" {
		return fmt.Errorf("SERVICE_ACCOUNT_FILE_CONTENT is not set")
	}

	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	if spreadsheetID == "" {
		log.Warn().Msg("SPREADSHEET_ID is not set, we'll create a new spreadsheet")
	}

	client := ravelry.NewRavelryClient("https://api.ravelry.com", login, password)

	username, err := client.Users.GetCurrentUsername(ctx)
	if err != nil {
		return fmt.Errorf("error fetching username: %w", err)
	}

	bundles, err := client.Bundles.GetUserBundles(ctx, username)
	if err != nil || len(bundles) == 0 {
		return fmt.Errorf("error fetching bundles or no bundles found: %w", err)
	}

	id := bundles[0].ID
	bundle, err := client.Bundles.GetBundleContent(ctx, username, id)
	if err != nil {
		return fmt.Errorf("error fetching bundle content: %w", err)
	}

	var data []gapi.RowData

	for _, item := range bundle.BundledItems {
		aa, err := client.Bundles.GetBundleItem(ctx, item.ID)
		if err != nil {
			return fmt.Errorf("error fetching bundle item: %w", err)
		}

		patternInfo, err := client.Patterns.GetPattern(ctx, aa.Item[0].ID)
		if err != nil {
			return fmt.Errorf("Error fetching pattern: %w", err)
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
