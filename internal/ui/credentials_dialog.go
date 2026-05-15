package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/rkhx/ravelry-exporter/internal/formui"
	"os"
	"path/filepath"
)

type Cfg struct {
	Credentials    Credentials `json:"credentials"`
	EmailToShare   string      `json:"emailToShare" formui:"label=Email to share,required,placeholder=E-Mail to share"`
	SpreadsheetID  string      `json:"spreadsheetID" formui:"label=Spreadsheet ID,placeholder=Spreadsheet ID"`
	ServiceAccount string      `json:"serviceAccount" formui:"label=Service Account,type=textarea,placeholder=Your notes"`
}

type Credentials struct {
	Username string `json:"username" formui:"label=Username,required,placeholder=Your API login"`
	Password string `json:"password" formui:"label=Password,required,placeholder=Your API password"`
}

func NewFormItemWithHint(label string, entry *widget.Entry, hint string) *widget.FormItem {
	item := widget.NewFormItem(label, entry)
	item.HintText = hint
	return item
}

func NewRequiredEntry(placeholder string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	entry.Validator = func(s string) error {
		if len(s) == 0 {
			return errors.New("this field is required")
		}
		return nil
	}
	return entry
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".ravelry-exporter", "creds.json"), nil
}

func ReadConfig() (*Cfg, error) {
	cfgPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("could not get config path: %w", err)
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Cfg
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return &cfg, nil
}

func writeConfig(cfg *Cfg) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".ravelry-exporter")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "creds.json")
	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func ShowSettingsDialog(w fyne.Window) {
	cancelBtn := widget.NewButton("Cancel", func() {
		w.Canvas().Overlays().Remove(w.Canvas().Overlays().Top())
	})

	formData := formui.BuildFormFromStruct(Cfg{})
	form := formData.Form
	form.Refresh()

	usernameEntry := formData.Entries["Credentials.Username"]
	passwordEntry := formData.Entries["Credentials.Password"]
	emailToShareEntry := formData.Entries["EmailToShare"]
	spreadsheetIDEntry := formData.Entries["SpreadsheetID"]
	serviceAccountEntry := formData.Entries["ServiceAccount"]

	if cfg, err := ReadConfig(); err == nil {
		usernameEntry.SetText(cfg.Credentials.Username)
		passwordEntry.SetText(cfg.Credentials.Password)
		emailToShareEntry.SetText(cfg.EmailToShare)
		spreadsheetIDEntry.SetText(cfg.SpreadsheetID)
	}

	saveBtn := widget.NewButton("Save", func() {
		cfg := &Cfg{
			Credentials: Credentials{
				Username: usernameEntry.Text,
				Password: passwordEntry.Text,
			},
			EmailToShare:   emailToShareEntry.Text,
			SpreadsheetID:  spreadsheetIDEntry.Text,
			ServiceAccount: serviceAccountEntry.Text,
		}

		if err := writeConfig(cfg); err != nil {
			dialog.ShowError(err, w)
			return
		}

		info := dialog.NewInformation("Success", "Credentials saved successfully!", w)

		info.SetOnClosed(func() {
			w.Canvas().Overlays().Remove(w.Canvas().Overlays().Top())
		})
		info.Show()
	})
	saveBtn.Disable()

	usernameEntry.OnChanged = func(content string) {
		usernameEntry.Validate()
		usernameEntry.Refresh()

		if err := form.Validate(); err == nil {
			saveBtn.Enable()
		}
	}

	passwordEntry.OnChanged = func(content string) {
		passwordEntry.Validate()
		passwordEntry.Refresh()

		if err := form.Validate(); err == nil {
			saveBtn.Enable()
		}
	}

	if err := form.Validate(); err == nil {
		saveBtn.Enable()
	}

	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(500, form.MinSize().Height)),
		form)

	modal := dialog.NewCustom("App Settings", "Close", grid, w)
	modal.SetButtons([]fyne.CanvasObject{saveBtn, cancelBtn})
	modal.Show()
}
