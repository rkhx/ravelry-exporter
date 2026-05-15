package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type CustomTheme struct {
	fyne.Theme
}

func (c *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return c.Theme.Size(name)
}

func (c *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameButton {
		return color.RGBA{R: 144, G: 213, B: 255, A: 255}
	}

	if name == theme.ColorNameBackground {
		return color.RGBA{R: 237, G: 232, B: 208, A: 255}
	}

	if name == theme.ColorNameForeground {
		return color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}

	if name == theme.ColorNameInputBorder {
		return color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}

	return c.Theme.Color(name, variant)
}
