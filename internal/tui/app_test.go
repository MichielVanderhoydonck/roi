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

	if a.prodService == nil {
		t.Error("expected prodService to be initialized")
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
func strPtr(s string) *string {
	return &s
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
