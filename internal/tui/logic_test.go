package tui

import (
	"testing"
)

func TestApp_CalculateResultDispatch(t *testing.T) {
	a := NewApp()

	// Initial result should be empty
	if a.resultText != "" {
		t.Errorf("expected empty resultText, got %s", a.resultText)
	}

	// Test productivity dispatch
	a.calculateResult(calcProductivity)
	if a.resultText == "" {
		t.Error("expected resultText to be populated after productivity calculation")
	}

	// Test reliability dispatch
	a.resultText = ""
	a.calculateResult(calcReliability)
	if a.resultText == "" {
		t.Error("expected resultText to be populated after reliability calculation")
	}
}

func TestApp_ResetState(t *testing.T) {
	a := NewApp()

	// Pretend a form was completed
	a.prodForm.State = 1 // huh.StateCompleted
	a.resultText = "some result"

	// Reset state
	a.resetFormState()
	if a.prodForm.State != 0 { // huh.StateNormal
		t.Errorf("expected form state to be Normal, got %v", a.prodForm.State)
	}

	// resultText is NOT reset by resetFormState currently (only by updateFocus switching)
	// but we can check the form states
	if a.relForm.State != 0 {
		t.Errorf("expected relForm state to be Normal, got %v", a.relForm.State)
	}
}
