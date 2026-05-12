package tui

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"charm.land/huh/v2"
)

type ContextSwitchCalculator struct {
	service          *service.ContextSwitchService
	reducedIncidents string
	hourlyRate       string
}

func NewContextSwitchCalculator() *ContextSwitchCalculator {
	return &ContextSwitchCalculator{
		service: service.NewContextSwitchService(),
	}
}

func (c *ContextSwitchCalculator) CreateForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("reducedIncidents").
				Title("Reduction in False Alerts/Flaky Builds per Year").
				Placeholder("e.g. 6000").
				Value(&c.reducedIncidents).
				Validate(validateInt),
			huh.NewInput().
				Key("hourlyRate").
				Title("Hourly Dev Rate ($)").
				Placeholder("e.g. 100").
				Value(&c.hourlyRate).
				Validate(validateFloat),
		),
	).WithShowHelp(false)
	return applyTheme(f)
}

func (c *ContextSwitchCalculator) GetContext(key string) string {
	help := map[string]string{
		"reducedIncidents": "How many fewer false alerts or flaky builds do you expect per year?\nExample: 6000",
		"hourlyRate":       "What is the fully loaded hourly rate of an engineer?\nExample: 100",
	}
	if val, ok := help[key]; ok {
		return val
	}
	return "Fill in the details to calculate the cost of context switching."
}

func (c *ContextSwitchCalculator) GetFormula(form *huh.Form) string {
	ri := getFormField(form, "reducedIncidents", c.reducedIncidents)
	hr := getFormField(form, "hourlyRate", c.hourlyRate)

	return fmt.Sprintf(`Context Switch Penalty Avoided

Savings ($) = 
  %s incidents 
  × 0.4 hours 
  × %s/hr`,
		formatFormulaValue(ri, "Reduced Incidents"),
		formatFormulaValue(hr, "Hourly Rate"))
}

func (c *ContextSwitchCalculator) CalculateResult(form *huh.Form) (string, Sentiment) {
	riStr := getFormField(form, "reducedIncidents", c.reducedIncidents)
	hrStr := getFormField(form, "hourlyRate", c.hourlyRate)

	if riStr == "" || hrStr == "" {
		return "", SentimentNone
	}

	ri, _ := strconv.Atoi(riStr)
	hr, _ := strconv.ParseFloat(hrStr, 64)

	res := c.service.Calculate(service.ContextSwitchInput{
		ReducedIncidentsPerYear: ri,
		HourlyRate:              hr,
	})

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Warning).MarginBottom(1)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(DefaultTheme.TextNormal)

	str := fmt.Sprintf("%s\n\n%s %.1f\n%s %s",
		titleStyle.Render("󱔗 CONTEXT SWITCH PENALTY"),
		labelStyle.Render("Total Hours Saved:"), res.HoursSaved,
		labelStyle.Render("Annual Savings Avoided:"), valStyle.Render(fmt.Sprintf("$%.2f", res.AnnualSavings)))

	return str, SentimentGood
}

func (c *ContextSwitchCalculator) Reset() {
	c.reducedIncidents = ""
	c.hourlyRate = ""
}
