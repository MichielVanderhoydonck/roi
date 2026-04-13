package tui

import (
	"testing"
)

func TestCreateProductivityForm(t *testing.T) {
	f := createProductivityForm()
	if f == nil {
		t.Fatal("expected form to be non-nil")
	}
}

func TestGetProductivityContext(t *testing.T) {
	tests := []struct {
		key      string
		contains string
	}{
		{"timeBefore", "process take manually"},
		{"timeAfter", "with your internal developer platform"},
		{"executions", "year does your engineering team execute"},
		{"invalid", "calculate your ROI"},
	}

	for _, tt := range tests {
		ctx := getProductivityContext(tt.key)
		if !contains(ctx, tt.contains) {
			t.Errorf("context for %s should contain %q, but got %q", tt.key, tt.contains, ctx)
		}
	}
}

func TestGetProductivityFormula(t *testing.T) {
	formula := getProductivityFormula(nil)
	if !contains(formula, "Time BEFORE") {
		t.Errorf("expected placeholder in formula, got %s", formula)
	}
}
