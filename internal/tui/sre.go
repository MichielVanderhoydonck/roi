package tui

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"github.com/charmbracelet/huh"
)

func createSREForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("hoursPerWeek").
				Title("Hours of Toil per Week").
				Placeholder("e.g. 5").
				Validate(validateFloat),
			huh.NewInput().
				Key("hourlyRate").
				Title("Average Developer Hourly Rate ($)").
				Placeholder("e.g. 75").
				Validate(validateFloat),
			huh.NewInput().
				Key("costToAutomate").
				Title("Cost to Automate ($)").
				Placeholder("e.g. 1500").
				Validate(validateFloat),
		),
	).WithShowHelp(false)
	return applyTheme(f)
}

func getSREContext(key string) string {
	help := map[string]string{
		"hoursPerWeek":   "How many hours per week does the team spend on manual, repetitive work (toil)?\nExample: 5",
		"hourlyRate":     "What is the fully loaded hourly cost of an engineer at your company?\nExample: 75",
		"costToAutomate": "How much does it cost in engineering time or licenses to automate this toil?\nExample: 1500",
	}
	if val, ok := help[key]; ok {
		return val
	}
	return "Fill in the SRE toil details to calculate your ROI."
}

func getSREFormula(form *huh.Form) string {
	var hpw, hr, cta string
	if form != nil {
		hpw = form.GetString("hoursPerWeek")
		hr = form.GetString("hourlyRate")
		cta = form.GetString("costToAutomate")
	}

	return fmt.Sprintf(`SRE Toil Eradication ROI

Annual ROI ($) = 
  [(%s × 52 weeks) 
   × %s] 
  - %s`,
		formatFormulaValue(hpw, "Hours per Week"),
		formatFormulaValue(hr, "Hourly Rate"),
		formatFormulaValue(cta, "Cost to Automate"))
}

func (a *App) calcSREResult() {
	hpw, _ := strconv.ParseFloat(a.sreForm.GetString("hoursPerWeek"), 64)
	hr, _ := strconv.ParseFloat(a.sreForm.GetString("hourlyRate"), 64)
	cta, _ := strconv.ParseFloat(a.sreForm.GetString("costToAutomate"), 64)

	res := a.sreService.Calculate(service.SREToilInput{
		HoursPerWeek:   hpw,
		HourlyRate:     hr,
		CostToAutomate: cta,
	})

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Primary)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success)
	a.resultText = fmt.Sprintf("%s\n\nTotal Hours Saved: %.1f\nNet Savings:       %s",
		titleStyle.Render("=== SRE Toil ROI Results ==="),
		res.HoursSaved,
		valStyle.Render(fmt.Sprintf("$%.2f", res.AnnualSavings)))
}
