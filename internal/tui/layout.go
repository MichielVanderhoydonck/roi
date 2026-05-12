package tui

import (
	"strings"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

func (a *App) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(DefaultTheme.Title).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1)

	focusText := "   Menu"
	if a.focus == focusForm {
		focusText = "   Calculator"
	}

	return lipgloss.JoinHorizontal(lipgloss.Left,
		headerStyle.Render("ROI CALCULATOR"),
		lipgloss.NewStyle().Foreground(DefaultTheme.TextDim).Render(focusText),
	)
}

func (a *App) renderFooter() string {
	footerStyle := lipgloss.NewStyle().
		Foreground(DefaultTheme.TextDim).
		Padding(0, 1)

	var helpStr string
	if a.focus == focusMenu {
		helpStr = "esc: quit • enter: select • j/k: navigate • ctrl+c: quit"
	} else {
		helpStr = "esc: menu • tab/enter: next • shift+tab: prev • ctrl+r: reset • ctrl+c: quit"
	}
	return footerStyle.Render("─── " + helpStr)
}

func (a *App) calculateLayoutWidths() (left, middle, right int) {
	left = max(30, min(45, a.width*30/100))
	remaining := max(0, a.width-left)
	middle = remaining * 55 / 100
	right = max(0, remaining-middle)
	return left, middle, right
}

func (a *App) renderLeftPanel(width, height int) string {
	borderColor := DefaultTheme.BorderDim
	if a.focus == focusMenu {
		borderColor = DefaultTheme.BorderActive
	}

	style := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)

	w, h := style.GetFrameSize()
	contentWidth := max(0, width-w)
	panelHeight := max(0, height-h)

	a.menuList.SetSize(contentWidth, panelHeight)

	return style.
		Width(width).
		Height(height).
		Render(a.menuList.View())
}

func (a *App) renderCalculationPanel(width, height int, formulaStr string) string {
	borderColor := DefaultTheme.BorderDim
	if a.focus == focusForm {
		borderColor = DefaultTheme.BorderActive
	}

	style := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)

	w, h := style.GetFrameSize()
	contentWidth := max(0, width-w)
	panelHeight := max(0, height-h)

	activeForm := a.getActiveForm()
	content := a.renderTopContent(contentWidth, formulaStr, activeForm)

	content = lipgloss.NewStyle().MaxHeight(panelHeight).Render(content)

	return style.
		Width(width).
		Height(height).
		Render(content)
}

func (a *App) renderResultPanel(width, height int) string {
	hasResult := a.resultText != "" && !strings.Contains(a.resultText, "Enter time values")
	borderColor := DefaultTheme.BorderDim
	if hasResult {
		switch a.resultSentiment {
		case SentimentGood:
			borderColor = DefaultTheme.Success
		case SentimentBad:
			borderColor = DefaultTheme.Critical
		}
	}

	style := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)

	w, _ := style.GetFrameSize()
	contentWidth := max(0, width-w)

	var content string
	if hasResult {
		content = lipgloss.NewStyle().
			Width(contentWidth).
			Render(a.resultText)
	} else {
		header := lipgloss.NewStyle().
			Bold(true).
			Foreground(DefaultTheme.TextDim).
			Render("󰄬 CALCULATION RESULT")

		var placeholderText string
		if a.resultText != "" {
			placeholderText = a.resultText
		} else {
			placeholderText = lipgloss.NewStyle().
				Foreground(DefaultTheme.TextDim).
				Render("Results will appear here\nas you fill in the form...")
		}

		content = lipgloss.NewStyle().
			Width(contentWidth).
			Render(header + "\n\n" + placeholderText)
	}

	return style.
		Width(width).
		Height(height).
		Render(content)
}

func (a *App) getCurrentPanelStrings() (formulaStr, contextStr string) {
	activeForm := a.getActiveForm()
	var fieldKey string

	if activeForm.State != huh.StateCompleted {
		if field := activeForm.GetFocusedField(); field != nil {
			fieldKey = field.GetKey()
		}
	}

	formulaStr = a.activeCalc.GetFormula(activeForm)
	contextStr = a.activeCalc.GetContext(fieldKey)

	return formulaStr, contextStr
}

func (a *App) renderTopContent(width int, formulaStr string, activeForm *huh.Form) string {
	formulaStyle := lipgloss.NewStyle().
		Foreground(DefaultTheme.Title).
		Bold(true).
		MarginBottom(1)

	formView := lipgloss.NewStyle().
		Width(width).
		MaxWidth(width).
		Render(activeForm.WithWidth(width).View())

	return formulaStyle.Render(formulaStr) + "\n\n" + formView
}

func (a *App) renderContextPanel(width int, contextStr string) string {
	style := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(DefaultTheme.BorderDim)

	w, _ := style.GetFrameSize()
	contentWidth := max(0, width-w)

	contextHeader := lipgloss.NewStyle().
		Bold(true).
		Foreground(DefaultTheme.TextDim).
		Render(" CONTEXT")

	contextStyle := lipgloss.NewStyle().
		Foreground(DefaultTheme.TextNormal).
		Width(contentWidth)

	content := contextHeader + "\n" + contextStyle.Render(contextStr)

	return style.
		Width(width).
		Render(content)
}
