package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	teav1 "github.com/charmbracelet/bubbletea"
)

func TestFormatFormulaValue(t *testing.T) {
	// Test empty value (placeholder)
	res := formatFormulaValue("", "Placeholder")
	if res == "" {
		t.Error("expected non-empty string")
	}

	// Test non-empty value
	res = formatFormulaValue("100", "Placeholder")
	if res == "" {
		t.Error("expected non-empty string")
	}
}

func TestMapV2MsgToV1(t *testing.T) {
	tests := []struct {
		v2msg    tea.Msg
		expected teav1.Msg
		name     string
	}{
		{
			name:     "Enter",
			v2msg:    tea.KeyPressMsg{Text: "enter"},
			expected: teav1.KeyMsg{Type: teav1.KeyEnter},
		},
		{
			name:     "Esc",
			v2msg:    tea.KeyPressMsg{Text: "esc"},
			expected: teav1.KeyMsg{Type: teav1.KeyEsc},
		},
		{
			name:     "WindowSize",
			v2msg:    tea.WindowSizeMsg{Width: 100, Height: 50},
			expected: teav1.WindowSizeMsg{Width: 100, Height: 50},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapV2MsgToV1(tt.v2msg)
			// Simple check since comparing interface types with complex structs can be tricky
			if tt.name == "WindowSize" {
				if g, ok := got.(teav1.WindowSizeMsg); !ok || g.Width != 100 {
					t.Errorf("got %v, expected %v", got, tt.expected)
				}
			} else {
				if g, ok := got.(teav1.KeyMsg); !ok || g.Type != tt.expected.(teav1.KeyMsg).Type {
					t.Errorf("got %v, expected %v", got, tt.expected)
				}
			}
		})
	}
}

func TestValidators(t *testing.T) {
	t.Run("validateDuration", func(t *testing.T) {
		if err := validateDuration("4h"); err != nil {
			t.Errorf("expected nil error for 4h, got %v", err)
		}
		if err := validateDuration("invalid"); err == nil {
			t.Error("expected error for invalid duration, got nil")
		}
	})

	t.Run("validateInt", func(t *testing.T) {
		if err := validateInt("100"); err != nil {
			t.Errorf("expected nil error for 100, got %v", err)
		}
		if err := validateInt("-10"); err == nil {
			t.Error("expected error for negative int, got nil")
		}
		if err := validateInt("abc"); err == nil {
			t.Error("expected error for non-int, got nil")
		}
	})

	t.Run("validateFloat", func(t *testing.T) {
		if err := validateFloat("75.5"); err != nil {
			t.Errorf("expected nil error for 75.5, got %v", err)
		}
		if err := validateFloat("-1.0"); err == nil {
			t.Error("expected error for negative float, got nil")
		}
		if err := validateFloat("abc"); err == nil {
			t.Error("expected error for non-float, got nil")
		}
	})
}
