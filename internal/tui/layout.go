package tui

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/huh"
)

func (a *App) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(DefaultTheme.Title).
		Bold(true).
		Padding(0, 1)

	focusText := " [ Menu ]"
	if a.focus == focusForm {
		focusText = " [ Calculator ]"
	}

	return lipgloss.JoinHorizontal(lipgloss.Left,
		headerStyle.Render("ROI Calculator"),
		lipgloss.NewStyle().Foreground(DefaultTheme.TextDim).Render(focusText),
	)
}

func (a *App) renderFooter() string {
	footerStyle := lipgloss.NewStyle().
		Foreground(DefaultTheme.TextDim).
		Padding(0, 1)

	var helpStr string
	if a.focus == focusMenu {
		helpStr = "(esc: close • enter/tab/→: select • ctrl+c: quit)"
	} else {
		activeForm := a.getActiveForm()
		if activeForm != nil && activeForm.State == huh.StateCompleted {
			helpStr = "(esc: back to menu • enter: edit form • ctrl+r: clear form • ctrl+c: quit)"
		} else {
			helpStr = "(esc: back to menu • ctrl+r: clear form • ctrl+c: quit)"
		}
	}
	return footerStyle.Render(helpStr)
}

func (a *App) calculateLayoutWidths() (left int, right int) {
	left = max(30, min(45, a.width*30/100))
	right = max(0, a.width-left)
	return left, right
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

	return style.
		Width(contentWidth).
		Height(panelHeight).
		Render(a.menuList.View())
}

func (a *App) renderRightPanel(width, height int) string {
	formulaStr, contextStr := a.getCurrentPanelStrings()

	bottomHeight := max(6, height*25/100)
	topHeight := max(0, height-bottomHeight)

	activeForm := a.getActiveForm()
	topContent := a.renderTopContent(formulaStr, activeForm)
	bottomContent := a.renderBottomContent(contextStr)

	borderColor := DefaultTheme.BorderDim
	if a.focus == focusForm {
		borderColor = DefaultTheme.BorderActive
	}

	topStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)

	bottomStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(DefaultTheme.BorderDim)

	topW, topH := topStyle.GetFrameSize()
	topBox := topStyle.
		Width(max(0, width-topW)).
		Height(max(0, topHeight-topH)).
		Render(topContent)

	bottomW, bottomH := bottomStyle.GetFrameSize()
	bottomBox := bottomStyle.
		Width(max(0, width-bottomW)).
		Height(max(0, bottomHeight-bottomH)).
		Render(bottomContent)

	return lipgloss.JoinVertical(lipgloss.Left, topBox, bottomBox)
}

func (a *App) getCurrentPanelStrings() (formulaStr, contextStr string) {
	selectedItem := a.menuList.SelectedItem().(item)
	activeForm := a.getActiveForm()
	var fieldKey string

	if activeForm.State != huh.StateCompleted {
		if field := activeForm.GetFocusedField(); field != nil {
			fieldKey = field.GetKey()
		}
	}

	switch selectedItem.calc {
	case calcProductivity:
		formulaStr = getProductivityFormula(activeForm)
		if activeForm.State != huh.StateCompleted {
			contextStr = getProductivityContext(fieldKey)
		}
	case calcReliability:
		formulaStr = getReliabilityFormula(activeForm)
		if activeForm.State != huh.StateCompleted {
			contextStr = getReliabilityContext(fieldKey)
		}
	case calcFinOps:
		formulaStr = getFinOpsFormula(activeForm)
		if activeForm.State != huh.StateCompleted {
			contextStr = getFinOpsContext(fieldKey)
		}
	case calcSRE:
		formulaStr = getSREFormula(activeForm)
		if activeForm.State != huh.StateCompleted {
			contextStr = getSREContext(fieldKey)
		}
	case calcOnboarding:
		formulaStr = getOnboardingFormula(activeForm)
		if activeForm.State != huh.StateCompleted {
			contextStr = getOnboardingContext(fieldKey)
		}
	case calcContextSwitch:
		formulaStr = getContextSwitchFormula(activeForm)
		if activeForm.State != huh.StateCompleted {
			contextStr = getContextSwitchContext(fieldKey)
		}
	case calcCostOfDelay:
		formulaStr = getCostOfDelayFormula(activeForm)
		if activeForm.State != huh.StateCompleted {
			contextStr = getCostOfDelayContext(fieldKey)
		}
	}
	return formulaStr, contextStr
}

func (a *App) renderTopContent(formulaStr string, activeForm *huh.Form) string {
	formulaStyle := lipgloss.NewStyle().
		Foreground(DefaultTheme.Title).
		Bold(true).
		MarginBottom(1)

	content := formulaStyle.Render(formulaStr) + "\n\n" + activeForm.View()

	if activeForm.State == huh.StateCompleted {
		resultBox := lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(DefaultTheme.Primary).
			Render(a.resultText)
		content += "\n" + resultBox
	}
	return content
}

func (a *App) renderBottomContent(contextStr string) string {
	contextHeader := lipgloss.NewStyle().
		Bold(true).
		Foreground(DefaultTheme.TextDim).
		Render("FIELD CONTEXT") + "\n"

	contextStyle := lipgloss.NewStyle().
		Foreground(DefaultTheme.TextNormal)

	return contextHeader + contextStyle.Render(contextStr)
}
