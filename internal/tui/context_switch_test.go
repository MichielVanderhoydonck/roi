package tui

import (
	"testing"
)

func TestCreateContextSwitchForm(t *testing.T) {
	f := createContextSwitchForm()
	if f == nil {
		t.Fatal("expected form to be non-nil")
	}
}

func TestGetContextSwitchContext(t *testing.T) {
	tests := []struct {
		key      string
		contains string
	}{
		{"reducedIncidents", "fewer false alerts or flaky builds"},
		{"hourlyRate", "fully loaded hourly rate"},
		{"invalid", "context switching"},
	}

	for _, tt := range tests {
		ctx := getContextSwitchContext(tt.key)
		if !contains(ctx, tt.contains) {
			t.Errorf("context for %s should contain %q, but got %q", tt.key, tt.contains, ctx)
		}
	}
}

func TestGetContextSwitchFormula(t *testing.T) {
	formula := getContextSwitchFormula(nil)
	if !contains(formula, "Reduced Incidents") {
		t.Errorf("expected placeholder in formula, got %s", formula)
	}
}
