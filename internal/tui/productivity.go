package tui

import (
	"fmt"
	"strconv"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"github.com/charmbracelet/huh"
)

func createProductivityForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("timeBefore").
				Title("Time spent on task BEFORE").
				Placeholder("e.g. 4h, 30m").
				Validate(validateDuration),
			huh.NewInput().
				Key("timeAfter").
				Title("Time spent on task AFTER").
				Placeholder("e.g. 5m, 10s").
				Validate(validateDuration),
			huh.NewInput().
				Key("executions").
				Title("Executions per year").
				Placeholder("e.g. 1000").
				Validate(validateInt),
			huh.NewInput().
				Key("hourlyRate").
				Title("Average Developer Hourly Rate ($)").
				Placeholder("e.g. 75").
				Validate(validateFloat),
			huh.NewInput().
				Key("maintenance").
				Title("Cost of Building/Maintaining Tool ($)").
				Placeholder("e.g. 1000").
				Validate(validateFloat),
		),
	).WithShowHelp(false)
	return applyTheme(f)
}

func getProductivityContext(key string) string {
	help := map[string]string{
		"timeBefore":  "How long did this process take manually before your automation?\nFormat: Valid time units are \"h\", \"m\", \"s\". Examples: \"4h\", \"30m\".",
		"timeAfter":   "How long does the process take now with your internal developer platform or automation?\nFormat: Valid time units are \"h\", \"m\", \"s\".",
		"executions":  "How many times per year does your engineering team execute this task?",
		"hourlyRate":  "What is the fully loaded hourly cost of an engineer at your company?\nExample: 75",
		"maintenance": "How much does it cost annually in software licenses or engineering time to maintain this tool?\nExample: 500",
	}
	if val, ok := help[key]; ok {
		return val
	}
	return "Fill in the productivity details to calculate your ROI."
}

func getProductivityFormula(form *huh.Form) string {
	var tb, ta, execs, hr, mc string
	if form != nil {
		tb = form.GetString("timeBefore")
		ta = form.GetString("timeAfter")
		execs = form.GetString("executions")
		hr = form.GetString("hourlyRate")
		mc = form.GetString("maintenance")
	}

	return fmt.Sprintf(`Developer Productivity ROI (Time Saved)

Annual ROI ($) = 
  [(%s - %s) 
   × %s 
   × %s] 
  - %s`,
		formatFormulaValue(tb, "Time BEFORE"),
		formatFormulaValue(ta, "Time AFTER"),
		formatFormulaValue(execs, "Executions per year"),
		formatFormulaValue(hr, "Hourly Rate"),
		formatFormulaValue(mc, "Maintenance Cost"))
}


func (a *App) calcProductivityResult() {
	tb, _ := time.ParseDuration(a.prodForm.GetString("timeBefore"))
	ta, _ := time.ParseDuration(a.prodForm.GetString("timeAfter"))
	execs, _ := strconv.Atoi(a.prodForm.GetString("executions"))
	hr, _ := strconv.ParseFloat(a.prodForm.GetString("hourlyRate"), 64)
	mc, _ := strconv.ParseFloat(a.prodForm.GetString("maintenance"), 64)

	res := a.prodService.Calculate(service.ProductivityInput{
		TimeBefore:        tb,
		TimeAfter:         ta,
		ExecutionsPerYear: execs,
		HourlyRate:        hr,
		MaintenanceCost:   mc,
	})

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Primary)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success)
	a.resultText = fmt.Sprintf("%s\n\nTotal Time Saved: %s\nGross Savings:    %s\nNet ROI:          %s",
		titleStyle.Render("=== Productivity ROI Results ==="),
		res.TimeSaved.String(),
		valStyle.Render(fmt.Sprintf("$%.2f", res.GrossSavings)),
		valStyle.Render(fmt.Sprintf("$%.2f", res.NetROI)))
}
