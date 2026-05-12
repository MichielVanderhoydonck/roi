package tui

import (
	"testing"
)

func TestFormatFormulaValue(t *testing.T) {
	// Test empty value (placeholder)
	res := formatFormulaValue("", "Placeholder")
	if res == "" {
		t.Error("expected non-empty string")
	}

	// Test non-empty value
	res = formatFormulaValue("100", "Placeholder")
	if res == "" {
		t.Error("expected non-empty string")
	}
}

func TestValidators(t *testing.T) {
	t.Run("validateDuration", func(t *testing.T) {
		if err := validateDuration("4h"); err != nil {
			t.Errorf("expected nil error for 4h, got %v", err)
		}
		if err := validateDuration("invalid"); err == nil {
			t.Error("expected error for invalid duration, got nil")
		}
	})

	t.Run("validateInt", func(t *testing.T) {
		if err := validateInt("100"); err != nil {
			t.Errorf("expected nil error for 100, got %v", err)
		}
		if err := validateInt("-10"); err == nil {
			t.Error("expected error for negative int, got nil")
		}
		if err := validateInt("abc"); err == nil {
			t.Error("expected error for non-int, got nil")
		}
	})

	t.Run("validateFloat", func(t *testing.T) {
		if err := validateFloat("75.5"); err != nil {
			t.Errorf("expected nil error for 75.5, got %v", err)
		}
		if err := validateFloat("-1.0"); err == nil {
			t.Error("expected error for negative float, got nil")
		}
		if err := validateFloat("abc"); err == nil {
			t.Error("expected error for non-float, got nil")
		}
	})
}
