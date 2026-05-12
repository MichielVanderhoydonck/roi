package tui

import (
	"fmt"
	"strconv"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"charm.land/huh/v2"
)

type ReliabilityCalculator struct {
	service      *service.ReliabilityService
	oldMTTR      string
	newMTTR      string
	incidents    string
	downtimeCost string
}

func NewReliabilityCalculator() *ReliabilityCalculator {
	return &ReliabilityCalculator{
		service: service.NewReliabilityService(),
	}
}

func (c *ReliabilityCalculator) CreateForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("oldMTTR").
				Title("Old MTTR").
				Placeholder("e.g. 2h, 45m").
				Value(&c.oldMTTR).
				Validate(validateDuration),
			huh.NewInput().
				Key("newMTTR").
				Title("New MTTR").
				Placeholder("e.g. 30m, 5m").
				Value(&c.newMTTR).
				Validate(validateDuration),
			huh.NewInput().
				Key("incidents").
				Title("Number of incidents per year").
				Placeholder("e.g. 10").
				Value(&c.incidents).
				Validate(validateInt),
			huh.NewInput().
				Key("downtimeCost").
				Title("Cost of Downtime per Hour ($)").
				Placeholder("e.g. 50000").
				Value(&c.downtimeCost).
				Validate(validateFloat),
		),
	).WithShowHelp(false)
	return applyTheme(f)
}

func (c *ReliabilityCalculator) GetContext(key string) string {
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

func (c *ReliabilityCalculator) GetFormula(form *huh.Form) string {
	om := getFormField(form, "oldMTTR", c.oldMTTR)
	nm := getFormField(form, "newMTTR", c.newMTTR)
	inc := getFormField(form, "incidents", c.incidents)
	dc := getFormField(form, "downtimeCost", c.downtimeCost)

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

func (c *ReliabilityCalculator) CalculateResult(form *huh.Form) (string, Sentiment) {
	omStr := getFormField(form, "oldMTTR", c.oldMTTR)
	nmStr := getFormField(form, "newMTTR", c.newMTTR)
	incStr := getFormField(form, "incidents", c.incidents)
	dcStr := getFormField(form, "downtimeCost", c.downtimeCost)

	if omStr == "" || nmStr == "" || incStr == "" {
		return "", SentimentNone
	}

	om, _ := time.ParseDuration(omStr)
	nm, _ := time.ParseDuration(nmStr)
	inc, _ := strconv.Atoi(incStr)
	dc, _ := strconv.ParseFloat(dcStr, 64)

	res := c.service.Calculate(service.ReliabilityInput{
		OldMTTR:          om,
		NewMTTR:          nm,
		IncidentsPerYear: inc,
		DowntimeCost:     dc,
	})

	sentiment := SentimentGood
	roiColor := DefaultTheme.Success
	if nm >= om {
		sentiment = SentimentBad
		roiColor = DefaultTheme.Critical
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Warning).MarginBottom(1)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success).Bold(true)
	roiStyle := lipgloss.NewStyle().Foreground(roiColor).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(DefaultTheme.TextNormal)

	str := fmt.Sprintf("%s\n\n%s %s\n%s %s",
		titleStyle.Render("󰚌 RELIABILITY IMPROVEMENT"),
		labelStyle.Render("Total Downtime Avoided:"), valStyle.Render(res.TimeSaved.String()),
		labelStyle.Render("Annual Downtime Savings:"), roiStyle.Render(fmt.Sprintf("$%.2f", res.DowntimeSavings)))

	return str, sentiment
}

func (c *ReliabilityCalculator) Reset() {
	c.oldMTTR = ""
	c.newMTTR = ""
	c.incidents = ""
	c.downtimeCost = ""
}
