package tui

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

func createReliabilityForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("oldMTTR").
				Title("Old MTTR").
				Placeholder("e.g. 2h, 45m").
				Validate(validateDuration),
			huh.NewInput().
				Key("newMTTR").
				Title("New MTTR").
				Placeholder("e.g. 30m, 5m").
				Validate(validateDuration),
			huh.NewInput().
				Key("incidents").
				Title("Number of incidents per year").
				Placeholder("e.g. 10").
				Validate(validateInt),
			huh.NewInput().
				Key("downtimeCost").
				Title("Cost of Downtime per Hour ($)").
				Placeholder("e.g. 50000").
				Validate(validateFloat),
		),
	).WithShowHelp(false)
	return applyTheme(f)
}

func getReliabilityContext(key string) string {
	help := map[string]string{
		"oldMTTR":      "Mean Time To Recovery before improvements.\nFormat: \"2h\", \"45m\".",
		"newMTTR":      "Mean Time To Recovery after implementing automated rollbacks and better observability.",
		"incidents":    "How many major incidents usually occur per year?",
		"downtimeCost": "Cost per hour of downtime, including lost revenue and SLA penalties.",
	}
	if val, ok := help[key]; ok {
		return val
	}
	return "Fill in the reliability details to calculate the cost of downtime avoided."
}

func getReliabilityFormula(form *huh.Form) string {
	var om, nm, inc, dc string
	if form != nil {
		om = form.GetString("oldMTTR")
		nm = form.GetString("newMTTR")
		inc = form.GetString("incidents")
		dc = form.GetString("downtimeCost")
	}

	return fmt.Sprintf(`Reliability ROI (Cost of Downtime Avoided)

Downtime Savings ($) = 
  (%s - %s) 
  × %s 
  × %s`,
		formatFormulaValue(om, "Old MTTR"),
		formatFormulaValue(nm, "New MTTR"),
		formatFormulaValue(inc, "Incidents per year"),
		formatFormulaValue(dc, "Cost of Downtime per Hour"))
}
