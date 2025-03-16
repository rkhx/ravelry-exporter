// Main pack
package main

import (
	"context"
	"fmt"
	"github.com/rkhx/ravelry-exporter/internal/gapi"
	"github.com/rkhx/ravelry-exporter/internal/ravelry"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	defer func() {
		if r := recover(); r != nil {
			pc, file, line, ok := runtime.Caller(3)
			if !ok {
				log.Fatalf("PANIC: %v", r)
			}
			fn := runtime.FuncForPC(pc)
			log.Fatalf("PANIC: %v in %s:%d (%s)", r, file, line, fn.Name())
		}
	}()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
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

	gapi.Spreadsheet(ctx, spreadsheetID, data, columns)
	return nil
}
