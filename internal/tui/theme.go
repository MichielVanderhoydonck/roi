package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"charm.land/huh/v2"
)

type Theme struct {
	Primary      color.Color // Main brand color (e.g., active borders, selected items)
	Secondary    color.Color // Secondary accents
	Title        color.Color // Main headers
	Success      color.Color // Positive ROI values
	Warning      color.Color // Warnings or important stats
	Critical     color.Color // Negative or bad ROI values
	TextNormal   color.Color // Body text
	TextDim      color.Color // De-emphasized text (help text, context)
	BorderActive color.Color // Active panel border
	BorderDim    color.Color // Inactive panel border

	HuhTheme huh.Theme
}

var DefaultTheme = Theme{
	Primary:      lipgloss.Color("#CBA6F7"), // Mauve
	Secondary:    lipgloss.Color("#89B4FA"), // Blue
	Title:        lipgloss.Color("#F5C2E7"), // Pink
	Success:      lipgloss.Color("#A6E3A1"), // Green
	Warning:      lipgloss.Color("#F9E2AF"), // Yellow
	Critical:     lipgloss.Color("#F38BA8"), // Red
	TextNormal:   lipgloss.Color("#CDD6F4"), // Text
	TextDim:      lipgloss.Color("#7F849C"), // Overlay1
	BorderActive: lipgloss.Color("#CBA6F7"), // Mauve
	BorderDim:    lipgloss.Color("#313244"), // Surface0

	HuhTheme: huh.ThemeFunc(huh.ThemeCatppuccin),
}

// applyTheme is a utility to apply our preferred huh theme to forms
func applyTheme(f *huh.Form) *huh.Form {
	return f.WithTheme(DefaultTheme.HuhTheme)
}
