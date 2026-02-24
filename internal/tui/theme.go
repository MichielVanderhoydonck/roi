package tui

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Primary      lipgloss.Color // Main brand color (e.g., active borders, selected items)
	Secondary    lipgloss.Color // Secondary accents
	Title        lipgloss.Color // Main headers
	Success      lipgloss.Color // Positive ROI values
	Warning      lipgloss.Color // Warnings or important stats
	TextNormal   lipgloss.Color // Body text
	TextDim      lipgloss.Color // De-emphasized text (help text, context)
	BorderActive lipgloss.Color // Active panel border
	BorderDim    lipgloss.Color // Inactive panel border

	HuhTheme *huh.Theme
}

var DefaultTheme = Theme{
	Primary:      lipgloss.Color("#CBA6F7"), // Mauve
	Secondary:    lipgloss.Color("#89B4FA"), // Blue
	Title:        lipgloss.Color("#F5C2E7"), // Pink
	Success:      lipgloss.Color("#A6DA95"), // Green
	Warning:      lipgloss.Color("#F9E2AF"), // Yellow
	TextNormal:   lipgloss.Color("#CDD6F4"), // Text
	TextDim:      lipgloss.Color("#A6ADC8"), // Subtext0
	BorderActive: lipgloss.Color("#CBA6F7"), // Mauve
	BorderDim:    lipgloss.Color("#45475A"), // Surface1

	HuhTheme: huh.ThemeCatppuccin(),
}

// applyTheme is a utility to apply our preferred huh theme to forms
func applyTheme(f *huh.Form) *huh.Form {
	return f.WithTheme(DefaultTheme.HuhTheme)
}
