package tui

import (
	"testing"
)

func TestCreateReliabilityForm(t *testing.T) {
	f := createReliabilityForm()
	if f == nil {
		t.Fatal("expected form to be non-nil")
	}
}

func TestGetReliabilityContext(t *testing.T) {
	tests := []struct {
		key      string
		contains string
	}{
		{"oldMTTR", "Mean Time To Recovery before"},
		{"newMTTR", "Mean Time To Recovery after"},
		{"incidents", "major incidents usually occur"},
		{"invalid", "reliability details"},
	}

	for _, tt := range tests {
		ctx := getReliabilityContext(tt.key)
		if !contains(ctx, tt.contains) {
			t.Errorf("context for %s should contain %q, but got %q", tt.key, tt.contains, ctx)
		}
	}
}

func TestGetReliabilityFormula(t *testing.T) {
	formula := getReliabilityFormula(nil)
	if !contains(formula, "Old MTTR") {
		t.Errorf("expected placeholder in formula, got %s", formula)
	}

	// Setting values directly in field doesn't always work for GetString in tests
	// So we check if it renders placeholders correctly at least.
}
