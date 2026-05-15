package main

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/rkhx/ravelry-exporter/internal/handlers"
	"github.com/rkhx/ravelry-exporter/internal/ui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

func credsFileExists() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf("cannot determine user home directory: %w", err)
	}

	credsPath := filepath.Join(homeDir, ".ravelry-exporter", "creds.json")
	_, err = os.Stat(credsPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("error checking creds file: %w", err)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	defer func() {
		if r := recover(); r != nil {
			pc, file, line, ok := runtime.Caller(3)
			fnName := "unknown"
			if ok {
				fnName = runtime.FuncForPC(pc).Name()
			}

			log.Fatal().
				Str("file", file).
				Int("line", line).
				Str("function", fnName).
				Msgf("PANIC: %v", r)
		}
	}()

	mainApp := app.New()
	mainApp.Settings().SetTheme(&ui.CustomTheme{Theme: theme.DefaultTheme()})

	mainWindow, logContainer := ui.GetWindow(ctx, mainApp, func(ctx context.Context, bundleUrl string) error {
		return handlers.Run(ctx, bundleUrl)
	})

	if ok, err := credsFileExists(); !ok || err != nil {
		ui.ShowSettingsDialog(mainWindow)
	}

	log.Logger = zerolog.New(ui.GetUIReporter(logContainer)).With().Timestamp().Logger()
	mainWindow.ShowAndRun()
}
