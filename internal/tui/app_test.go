package tui

import (
	"regexp"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestNewApp(t *testing.T) {
	a := NewApp()
	if a == nil {
		t.Fatal("expected NewApp to return non-nil")
	}

	if a.focus != focusMenu {
		t.Errorf("expected initial focus to be focusMenu, got %v", a.focus)
	}

	if a.activeCalc == nil {
		t.Error("expected activeCalc to be initialized")
	}
}

func TestApp_WindowSize(t *testing.T) {
	a := NewApp()
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	a.Update(msg)

	if a.width != 100 || a.height != 50 {
		t.Errorf("expected width/height to be 100/50, got %v/%v", a.width, a.height)
	}
}

func TestApp_Navigation(t *testing.T) {
	a := NewApp()

	// Initial state
	if a.focus != focusMenu {
		t.Errorf("initial focus: %v", a.focus)
	}

	// Press Enter to go to form
	msg := tea.KeyPressMsg{Text: "enter"}
	a.Update(msg)

	if a.focus != focusForm {
		t.Errorf("expected focusForm after enter, got %v", a.focus)
	}

	// Press Esc to return to menu
	msg = tea.KeyPressMsg{Text: "esc"}
	a.Update(msg)

	if a.focus != focusMenu {
		t.Errorf("expected focusMenu after esc, got %v", a.focus)
	}
}

// Helpers for TUI testing

func TestApp_FormPlaceholderVisibility(t *testing.T) {
	a := NewApp()

	// Send WindowSizeMsg with a typical layout size
	a.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Generate View
	v := a.View()

	// Verify full placeholder exists instead of just truncated 'e'
	if !contains(v.Content, "e.g. 4h, 30m") {
		t.Errorf("expected view content to contain full placeholder 'e.g. 4h, 30m'")
	}
}

func TestApp_ResultPersistence(t *testing.T) {
	a := NewApp()
	a.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Enter form
	a.Update(tea.KeyPressMsg{Text: "enter"})

	expected := "Valid ROI Computed Result"
	a.resultText = expected

	v := a.View()
	if !contains(v.Content, expected) {
		t.Errorf("expected view to persist result text")
	}
}

func contains(s, substr string) bool {
	cleanS := stripANSI(s)
	for i := 0; i <= len(cleanS)-len(substr); i++ {
		if cleanS[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func stripANSI(str string) string {
	re := regexp.MustCompile("[\u001B\u009B][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]")
	return re.ReplaceAllString(str, "")
}
