package tui

import (
	"testing"
	tea "charm.land/bubbletea/v2"
)

func TestApp_CalculateResultDispatch(t *testing.T) {
	a := NewApp()

	if a.resultText != "" {
		t.Errorf("expected empty resultText, got %s", a.resultText)
	}

	a.activeCalc = NewProductivityCalculator()
	a.activeForm = a.activeCalc.CreateForm()
	a.resultText = a.activeCalc.CalculateResult(a.activeForm)
	
	if a.resultText == "" {
		t.Error("expected resultText to be populated after productivity calculation")
	}

	a.resultText = ""
	a.activeCalc = NewReliabilityCalculator()
	a.activeForm = a.activeCalc.CreateForm()
	a.resultText = a.activeCalc.CalculateResult(a.activeForm)
	
	if a.resultText == "" {
		t.Error("expected resultText to be populated after reliability calculation")
	}
}

func TestApp_ResetState(t *testing.T) {
	a := NewApp()
	a.focus = focusForm
	a.activeForm.State = 1 // huh.StateCompleted
	
	handled, _ := a.handleFormKey(tea.KeyPressMsg{Text: "esc"})
	
	if !handled {
		t.Error("expected handleFormKey to handle esc")
	}
	
	if a.activeForm.State != 0 {
		t.Errorf("expected form state to be Normal, got %v", a.activeForm.State)
	}
	if a.focus != focusMenu {
		t.Errorf("expected focus to be menu, got %v", a.focus)
	}
}
