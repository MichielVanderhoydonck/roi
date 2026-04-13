package tui

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"github.com/charmbracelet/huh"
)

func createFinOpsForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("oldBill").
				Title("Previous Monthly Cloud Bill ($)").
				Placeholder("e.g. 20000").
				Validate(validateFloat),
			huh.NewInput().
				Key("newBill").
				Title("New Monthly Cloud Bill ($)").
				Placeholder("e.g. 15000").
				Validate(validateFloat),
		),
	).WithShowHelp(false)
	return applyTheme(f)
}

func getFinOpsContext(key string) string {
	help := map[string]string{
		"oldBill": "Your average monthly cloud infrastructure bill before optimization.",
		"newBill": "Your target or actual monthly cloud bill after right-sizing or spinning down resources.",
	}
	if val, ok := help[key]; ok {
		return val
	}
	return "Fill in your cloud bills to calculate infrastructure savings."
}

func getFinOpsFormula(form *huh.Form) string {
	var ob, nb string
	if form != nil {
		ob = form.GetString("oldBill")
		nb = form.GetString("newBill")
	}

	return fmt.Sprintf(`FinOps ROI (Infrastructure Optimization)

Cloud Savings ($) = 
  (%s - %s) 
  × 12 months`,
		formatFormulaValue(ob, "Previous Monthly Bill"),
		formatFormulaValue(nb, "New Monthly Bill"))
}

func (a *App) calcFinOpsResult() {
	ob, _ := strconv.ParseFloat(a.finForm.GetString("oldBill"), 64)
	nb, _ := strconv.ParseFloat(a.finForm.GetString("newBill"), 64)

	res := a.finService.Calculate(service.FinOpsInput{
		OldMonthlyBill: ob,
		NewMonthlyBill: nb,
	})

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Secondary)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success)
	a.resultText = fmt.Sprintf("%s\n\nAnnual Cloud Savings: %s",
		titleStyle.Render("=== FinOps ROI Results ==="),
		valStyle.Render(fmt.Sprintf("$%.2f", res.AnnualSavings)))
}
