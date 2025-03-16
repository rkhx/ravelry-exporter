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
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if os.Getenv("LOG_FORMAT") == "pretty" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	defer func() {
		if r := recover(); r != nil {
			pc, file, line, ok := runtime.Caller(3)
			if !ok {
				log.Fatal().Msgf("PANIC: %v", r)
			}
			fn := runtime.FuncForPC(pc)
			log.Fatal().Msgf("PANIC: %v in %s:%d (%s)", r, file, line, fn.Name())
		}
	}()

	start := time.Now()
	if err := run(ctx); err != nil {
		log.Fatal().Msgf("%s", err.Error())
	}
	log.Info().Dur("uptime", time.Since(start)).Msg("Application exited successfully")
}

func run(ctx context.Context) error {
	login, password := os.Getenv("RAVELRY_LOGIN"), os.Getenv("RAVELRY_PASSWORD")
	if login == "" || password == "" {
		return fmt.Errorf("RAVELRY_LOGIN or RAVELRY_PASSWORD is not set")
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
		fmt.Println("Error fetching bundles:", err)
	}

	var data []gapi.RowData

	for _, item := range bundle.BundledItems {
		aa, err := client.Bundles.GetBundleItem(ctx, item.ID)
		if err != nil {
			fmt.Println("Error fetching bundles:", err)
		}

		patternInfo, err := client.Patterns.GetPattern(ctx, aa.Item[0].ID)
		if err != nil {
			fmt.Println("Error fetching bundles:", err)
		}

		data = append(data, gapi.NewRowData(patternInfo))
	}

	columns := []interface{}{"Pattern Name", "Designer", "Gauge, needle size", "Sizes", "Recommended yarn", "Attributes"}

	spreadsheetID := os.Getenv("SPREADSHEET_ID")

	err = gapi.Spreadsheet(ctx, spreadsheetID, data, columns)
	if err != nil {
		return fmt.Errorf("error fetching spreadsheet: %w", err)
	}

	return nil
}
