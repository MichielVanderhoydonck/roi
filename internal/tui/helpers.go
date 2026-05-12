package tui

import (
	"fmt"
	"strconv"
	"time"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

func validateDuration(s string) error {
	_, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration format")
	}
	return nil
}

func validateInt(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil || v < 0 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}

func validateFloat(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v < 0 {
		return fmt.Errorf("must be a positive number")
	}
	return nil
}

// formatFormulaValue highlights a value if it's present, otherwise returns a dimmed placeholder.
func formatFormulaValue(val, fallback string) string {
	if val == "" {
		return lipgloss.NewStyle().Foreground(DefaultTheme.TextDim).Render(fallback)
	}
	return lipgloss.NewStyle().Foreground(DefaultTheme.Success).Render(val)
}

// getFormField retrieves the live field string from the form if available, falling back to the persistent struct value.
func getFormField(form *huh.Form, key string, fallback string) string {
	if form != nil {
		if val := form.GetString(key); val != "" {
			return val
		}
	}
	return fallback
}
