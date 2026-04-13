package tui

import (
	"fmt"
	"strconv"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
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

func (a *App) calcReliabilityResult() {
	om, _ := time.ParseDuration(a.relForm.GetString("oldMTTR"))
	nm, _ := time.ParseDuration(a.relForm.GetString("newMTTR"))
	inc, _ := strconv.Atoi(a.relForm.GetString("incidents"))
	dc, _ := strconv.ParseFloat(a.relForm.GetString("downtimeCost"), 64)

	res := a.relService.Calculate(service.ReliabilityInput{
		OldMTTR:          om,
		NewMTTR:          nm,
		IncidentsPerYear: inc,
		DowntimeCost:     dc,
	})

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Warning)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success)
	a.resultText = fmt.Sprintf("%s\n\nTotal Downtime avoided: %s\nDowntime Savings:       %s",
		titleStyle.Render("=== Reliability ROI Results ==="),
		res.TimeSaved.String(),
		valStyle.Render(fmt.Sprintf("$%.2f", res.DowntimeSavings)))
}
