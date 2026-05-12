package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestCreateFinOpsForm(t *testing.T) {
	f := NewFinOpsCalculator().CreateForm()
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
		ctx := NewFinOpsCalculator().GetContext(tt.key)
		if !contains(ctx, tt.contains) {
			t.Errorf("context for %s should contain %q, but got %q", tt.key, tt.contains, ctx)
		}
	}
}

func TestGetFinOpsFormula(t *testing.T) {
	formula := NewFinOpsCalculator().GetFormula(nil)
	if !contains(formula, "Previous Monthly Bill") {
		t.Errorf("expected placeholder in formula, got %s", formula)
	}
}

func TestApp_FinOps_E2EFlow(t *testing.T) {
	a := NewApp()

	// Step 1: Initialize screen layout sizes
	a.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Verify initial View content of the three panes
	content := a.View().Content
	if !contains(content, "ROI CALCULATOR") {
		t.Errorf("expected Header in view, got:\n%s", content)
	}
	if !contains(content, "Developer Productivity") {
		t.Errorf("expected default Menu items in left pane, got:\n%s", content)
	}
	if !contains(content, "Enter time values") {
		t.Errorf("expected initial placeholder prompt in right pane, got:\n%s", content)
	}

	// Step 2: Navigate to FinOps calculator at index 2 in the menu list
	a.menuList.Select(2)
	selectedItem := a.menuList.SelectedItem().(item)
	a.activeCalc = selectedItem.calc
	a.activeCalc.Reset()
	a.activeForm = a.activeCalc.CreateForm()
	a.resultText, _ = a.activeCalc.CalculateResult(a.activeForm)

	// Verify the menu and context panes update correctly for FinOps
	content = a.View().Content
	if !contains(content, "FinOps") {
		t.Errorf("expected FinOps item visible in menu, got:\n%s", content)
	}
	if !contains(content, "infrastructure bill") {
		t.Errorf("expected FinOps context string in middle pane, got:\n%s", content)
	}

	// Step 3: Switch focus to the form pane
	a.focus = focusForm

	// Step 4: Simulate filling in the form fields dynamically
	// We bind the underlying struct pointers and sync the active form to simulate completed user input
	finOpsCalc, ok := a.activeCalc.(*FinOpsCalculator)
	if !ok {
		t.Fatal("expected activeCalc to be FinOpsCalculator")
	}
	finOpsCalc.oldBill = "20000"
	finOpsCalc.newBill = "15000"
	a.activeForm = a.activeCalc.CreateForm()

	// Trigger live calculation evaluation as form updates do
	a.resultText, _ = a.activeCalc.CalculateResult(a.activeForm)

	// Step 5: Verify the dynamically calculated outcome in the Result pane
	content = a.View().Content
	if !contains(content, "FINOPS") {
		t.Errorf("expected FinOps optimization title in result pane, got:\n%s", content)
	}
	if !contains(content, "5000") {
		t.Errorf("expected Monthly Savings of 5000 in result pane, got:\n%s", content)
	}
	if !contains(content, "60000") {
		t.Errorf("expected Annual Cloud Savings of 60000 in result pane, got:\n%s", content)
	}
}
