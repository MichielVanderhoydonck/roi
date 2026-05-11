package tui

import (
	"testing"
)

func TestCreateDORAAIForm(t *testing.T) {
	f := NewDORAAICalculator().CreateForm()
	if f == nil {
		t.Fatal("expected form to be non-nil")
	}
}

func TestGetDORAAIContext(t *testing.T) {
	tests := []struct {
		key      string
		contains string
	}{
		{"staffSize", "count of full-time employees"},
		{"salary", "Average fully loaded"},
		{"invalid", "calculate your ROI"},
	}

	for _, tt := range tests {
		ctx := NewDORAAICalculator().GetContext(tt.key)
		if !contains(ctx, tt.contains) {
			t.Errorf("context for %s should contain %q, but got %q", tt.key, tt.contains, ctx)
		}
	}
}

func TestGetDORAAIFormula(t *testing.T) {
	formula := NewDORAAICalculator().GetFormula(nil)
	if !contains(formula, "DORA ROI of AI-assisted") {
		t.Errorf("expected DORA ROI of AI-assisted in formula, got %s", formula)
	}
}
