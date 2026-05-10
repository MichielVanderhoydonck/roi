package tui

import (
	"testing"
)

func TestCreateOnboardingForm(t *testing.T) {
	f := NewOnboardingCalculator().CreateForm()
	if f == nil {
		t.Fatal("expected form to be non-nil")
	}
}

func TestGetOnboardingContext(t *testing.T) {
	tests := []struct {
		key      string
		contains string
	}{
		{"oldDays", "days did it take for a new hire"},
		{"newHires", "engineers do you plan to hire"},
		{"invalid", "onboarding details"},
	}

	for _, tt := range tests {
		ctx := NewOnboardingCalculator().GetContext(tt.key)
		if !contains(ctx, tt.contains) {
			t.Errorf("context for %s should contain %q, but got %q", tt.key, tt.contains, ctx)
		}
	}
}

func TestGetOnboardingFormula(t *testing.T) {
	formula := NewOnboardingCalculator().GetFormula(nil)
	if !contains(formula, "Old Days") {
		t.Errorf("expected placeholder in formula, got %s", formula)
	}
}
