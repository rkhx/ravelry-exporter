package ui

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
	"time"
)

func GetWindow(ctx context.Context, myApp fyne.App, run func(ctx context.Context, text string) error) (fyne.Window, *widget.RichText) {
	myWindow := myApp.NewWindow("Ravelry Exporter")

	logOutput := widget.NewRichText()
	logOutput.Wrapping = fyne.TextWrapWord

	scrollContainer := container.NewScroll(logOutput)

	statusEntry := widget.NewEntry()

	statusEntry.SetPlaceHolder("Put bundle URL here...")

	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		ShowSettingsDialog(myWindow)
	})
	settingsButton.Resize(fyne.NewSize(settingsButton.MinSize().Height*2, settingsButton.MinSize().Height))

	c := container.NewBorder(nil, nil, nil, settingsButton, statusEntry)

	btn := widget.NewButton("Start", nil)
	btn.OnTapped = func() {
		btn.Disable()
		logOutput.Segments = []widget.RichTextSegment{}
		start := time.Now()
		err := run(ctx, statusEntry.Text)
		if err != nil {
			log.Error().Err(err).Msg("Error occurred")
		} else {
			log.Info().Dur("Execution time", time.Since(start)).Msg("Export completed successfully!")
		}
		btn.Enable()
	}

	myWindow.SetContent(container.NewBorder(
		c,
		btn,
		nil,
		nil,
		scrollContainer,
	))

	myWindow.Resize(fyne.NewSize(800, 600))

	return myWindow, logOutput
}
