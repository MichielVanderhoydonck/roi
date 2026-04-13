package tui

import (
	"testing"
)

func TestCreateCostOfDelayForm(t *testing.T) {
	f := createCostOfDelayForm()
	if f == nil {
		t.Fatal("expected form to be non-nil")
	}
}

func TestGetCostOfDelayContext(t *testing.T) {
	tests := []struct {
		key      string
		contains string
	}{
		{"monthlyRevenue", "anticipated monthly revenue"},
		{"daysDelayed", "How many days was the launch delayed"},
		{"invalid", "cost of delay"},
	}

	for _, tt := range tests {
		ctx := getCostOfDelayContext(tt.key)
		if !contains(ctx, tt.contains) {
			t.Errorf("context for %s should contain %q, but got %q", tt.key, tt.contains, ctx)
		}
	}
}

func TestGetCostOfDelayFormula(t *testing.T) {
	formula := getCostOfDelayFormula(nil)
	if !contains(formula, "Monthly Revenue") {
		t.Errorf("expected placeholder in formula, got %s", formula)
	}
}
