package tui

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"github.com/charmbracelet/huh"
)

func createContextSwitchForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("reducedIncidents").
				Title("Reduction in False Alerts/Flaky Builds per Year").
				Placeholder("e.g. 6000").
				Validate(validateInt),
			huh.NewInput().
				Key("hourlyRate").
				Title("Hourly Dev Rate ($)").
				Placeholder("e.g. 100").
				Validate(validateFloat),
		),
	).WithShowHelp(false)
	return applyTheme(f)
}

func getContextSwitchContext(key string) string {
	help := map[string]string{
		"reducedIncidents": "How many fewer false alerts or flaky builds do you expect per year?\nExample: 6000",
		"hourlyRate":       "What is the fully loaded hourly rate of an engineer?\nExample: 100",
	}
	if val, ok := help[key]; ok {
		return val
	}
	return "Fill in the details to calculate the cost of context switching."
}

func getContextSwitchFormula(form *huh.Form) string {
	var ri, hr string
	if form != nil {
		ri = form.GetString("reducedIncidents")
		hr = form.GetString("hourlyRate")
	}

	return fmt.Sprintf(`Context Switch Penalty Avoided

Savings ($) = 
  %s incidents 
  × 0.4 hours 
  × %s/hr`,
		formatFormulaValue(ri, "Reduced Incidents"),
		formatFormulaValue(hr, "Hourly Rate"))
}

func (a *App) calcContextSwitchResult() {
	ri, _ := strconv.Atoi(a.contextSwitchForm.GetString("reducedIncidents"))
	hr, _ := strconv.ParseFloat(a.contextSwitchForm.GetString("hourlyRate"), 64)

	res := a.contextSwitchService.Calculate(service.ContextSwitchInput{
		ReducedIncidentsPerYear: ri,
		HourlyRate:              hr,
	})

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Warning)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success)
	a.resultText = fmt.Sprintf("%s\n\nHours Saved: %.1f\nContext Switch Penalty Avoided: %s",
		titleStyle.Render("=== Context Switch ROI Results ==="),
		res.HoursSaved,
		valStyle.Render(fmt.Sprintf("$%.2f", res.AnnualSavings)))
}
