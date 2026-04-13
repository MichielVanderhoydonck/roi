package tui

import (
	"testing"
)

func TestCreateFinOpsForm(t *testing.T) {
	f := createFinOpsForm()
	if f == nil {
		t.Fatal("expected form to be non-nil")
	}
}

func TestGetFinOpsContext(t *testing.T) {
	tests := []struct {
		key      string
		contains string
	}{
		{"oldBill", "average monthly cloud infrastructure bill"},
		{"newBill", "target or actual monthly cloud bill"},
		{"invalid", "infrastructure savings"},
	}

	for _, tt := range tests {
		ctx := getFinOpsContext(tt.key)
		if !contains(ctx, tt.contains) {
			t.Errorf("context for %s should contain %q, but got %q", tt.key, tt.contains, ctx)
		}
	}
}

func TestGetFinOpsFormula(t *testing.T) {
	formula := getFinOpsFormula(nil)
	if !contains(formula, "Previous Monthly Bill") {
		t.Errorf("expected placeholder in formula, got %s", formula)
	}
}
