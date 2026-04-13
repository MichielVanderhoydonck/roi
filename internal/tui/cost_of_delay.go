package tui

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"github.com/charmbracelet/huh"
)

func createCostOfDelayForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("monthlyRevenue").
				Title("Estimated Monthly Revenue of Feature ($)").
				Placeholder("e.g. 300000").
				Validate(validateFloat),
			huh.NewInput().
				Key("daysDelayed").
				Title("Days Delayed").
				Placeholder("e.g. 15").
				Validate(validateFloat),
		),
	).WithShowHelp(false)
	return applyTheme(f)
}

func getCostOfDelayContext(key string) string {
	help := map[string]string{
		"monthlyRevenue": "What is the anticipated monthly revenue this feature will generate?\nExample: 300000",
		"daysDelayed":    "How many days was the launch delayed due to bottlenecks?\nExample: 15",
	}
	if val, ok := help[key]; ok {
		return val
	}
	return "Fill in the details to calculate the cost of delay."
}

func getCostOfDelayFormula(form *huh.Form) string {
	var mr, dd string
	if form != nil {
		mr = form.GetString("monthlyRevenue")
		dd = form.GetString("daysDelayed")
	}

	return fmt.Sprintf(`Cost of Delay

Revenue Lost ($) = 
  (%s / 30) 
  × %s days`,
		formatFormulaValue(mr, "Monthly Revenue"),
		formatFormulaValue(dd, "Days Delayed"))
}

func (a *App) calcCostOfDelayResult() {
	mr, _ := strconv.ParseFloat(a.costOfDelayForm.GetString("monthlyRevenue"), 64)
	dd, _ := strconv.ParseFloat(a.costOfDelayForm.GetString("daysDelayed"), 64)

	res := a.costOfDelayService.Calculate(service.CostOfDelayInput{
		EstimatedMonthlyRevenue: mr,
		DaysDelayed:             dd,
	})

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Primary)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Warning)
	a.resultText = fmt.Sprintf("%s\n\nRevenue Lost (Cost of Delay): %s",
		titleStyle.Render("=== Cost of Delay Results ==="),
		valStyle.Render(fmt.Sprintf("$%.2f", res.CostOfDelay)))
}
