package tui

import (
	"fmt"
	"strconv"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	teav1 "github.com/charmbracelet/bubbletea"
)

func validateDuration(s string) error {
	_, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration format")
	}
	return nil
}

func validateInt(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil || v < 0 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}

func validateFloat(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v < 0 {
		return fmt.Errorf("must be a positive number")
	}
	return nil
}

// formatFormulaValue highlights a value if it's present, otherwise returns a dimmed placeholder.
func formatFormulaValue(val, fallback string) string {
	if val == "" {
		return lipgloss.NewStyle().Foreground(DefaultTheme.TextDim).Render(fallback)
	}
	return lipgloss.NewStyle().Foreground(DefaultTheme.Success).Render(val)
}

// Interoperability Helpers for Bubble Tea v1 and v2

// wrapCmd wraps a v1 tea.Cmd to work with v2 tea.Msg.
func wrapCmd(cmd teav1.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}
	return func() tea.Msg {
		msg := cmd()
		if b, ok := msg.(teav1.BatchMsg); ok {
			var cmds []tea.Cmd
			for _, c := range b {
				cmds = append(cmds, wrapCmd(c))
			}
			return tea.BatchMsg(cmds)
		}

		// Map standard v1 messages back to v2 if needed.
		if _, ok := msg.(teav1.QuitMsg); ok {
			return tea.Quit()
		}
		return msg
	}
}

// mapV2MsgToV1 maps v2 tea.Msg to v1 teav1.Msg (for huh components).
func mapV2MsgToV1(msg tea.Msg) teav1.Msg {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		var k teav1.KeyMsg
		switch msg.String() {
		case "ctrl+c":
			k.Type = teav1.KeyCtrlC
		case "enter":
			k.Type = teav1.KeyEnter
		case "esc":
			k.Type = teav1.KeyEsc
		case "tab":
			k.Type = teav1.KeyTab
		case "shift+tab":
			k.Type = teav1.KeyShiftTab
		case "up":
			k.Type = teav1.KeyUp
		case "down":
			k.Type = teav1.KeyDown
		case "left":
			k.Type = teav1.KeyLeft
		case "right":
			k.Type = teav1.KeyRight
		case "backspace":
			k.Type = teav1.KeyBackspace
		case "delete":
			k.Type = teav1.KeyDelete
		case "space":
			k.Type = teav1.KeySpace
			k.Runes = []rune{' '}
		default:
			k.Type = teav1.KeyRunes
			k.Runes = []rune(msg.Text)
		}

		if msg.Mod.Contains(tea.ModAlt) {
			k.Alt = true
		}
		return k
	case tea.WindowSizeMsg:
		return teav1.WindowSizeMsg{Width: msg.Width, Height: msg.Height}
	}
	return msg
}
