package tui

import (
	"fmt"

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
