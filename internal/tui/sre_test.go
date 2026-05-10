package tui

import (
	"testing"
)

func TestCreateSREForm(t *testing.T) {
	f := NewSRECalculator().CreateForm()
	if f == nil {
		t.Fatal("expected form to be non-nil")
	}
}

func TestGetSREContext(t *testing.T) {
	tests := []struct {
		key      string
		contains string
	}{
		{"hoursPerWeek", "How many hours per week does the team spend"},
		{"hourlyRate", "fully loaded hourly cost"},
		{"invalid", "SRE toil details"},
	}

	for _, tt := range tests {
		ctx := NewSRECalculator().GetContext(tt.key)
		if !contains(ctx, tt.contains) {
			t.Errorf("context for %s should contain %q, but got %q", tt.key, tt.contains, ctx)
		}
	}
}

func TestGetSREFormula(t *testing.T) {
	formula := NewSRECalculator().GetFormula(nil)
	if !contains(formula, "Hours per Week") {
		t.Errorf("expected placeholder in formula, got %s", formula)
	}
}
