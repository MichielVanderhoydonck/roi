package tui

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"charm.land/huh/v2"
)

type SRECalculator struct {
	service        *service.SREToilService
	hoursPerWeek   string
	hourlyRate     string
	costToAutomate string
}

func NewSRECalculator() *SRECalculator {
	return &SRECalculator{
		service: service.NewSREToilService(),
	}
}

func (c *SRECalculator) CreateForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("hoursPerWeek").
				Title("Hours of Toil per Week").
				Placeholder("e.g. 5").
				Value(&c.hoursPerWeek).
				Validate(validateFloat),
			huh.NewInput().
				Key("hourlyRate").
				Title("Average Developer Hourly Rate ($)").
				Placeholder("e.g. 75").
				Value(&c.hourlyRate).
				Validate(validateFloat),
			huh.NewInput().
				Key("costToAutomate").
				Title("Cost to Automate ($)").
				Placeholder("e.g. 1500").
				Value(&c.costToAutomate).
				Validate(validateFloat),
		),
	).WithShowHelp(false)
	return applyTheme(f)
}

func (c *SRECalculator) GetContext(key string) string {
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

func (c *SRECalculator) GetFormula(form *huh.Form) string {
	hpw := getFormField(form, "hoursPerWeek", c.hoursPerWeek)
	hr := getFormField(form, "hourlyRate", c.hourlyRate)
	cta := getFormField(form, "costToAutomate", c.costToAutomate)

	return fmt.Sprintf(`SRE Toil Eradication ROI

Annual ROI ($) = 
  [(%s × 52 weeks) 
   × %s] 
  - %s`,
		formatFormulaValue(hpw, "Hours per Week"),
		formatFormulaValue(hr, "Hourly Rate"),
		formatFormulaValue(cta, "Cost to Automate"))
}

func (c *SRECalculator) CalculateResult(form *huh.Form) (string, Sentiment) {
	hpwStr := getFormField(form, "hoursPerWeek", c.hoursPerWeek)
	hrStr := getFormField(form, "hourlyRate", c.hourlyRate)
	ctaStr := getFormField(form, "costToAutomate", c.costToAutomate)

	if hpwStr == "" || hrStr == "" {
		return "", SentimentNone
	}

	hpw, _ := strconv.ParseFloat(hpwStr, 64)
	hr, _ := strconv.ParseFloat(hrStr, 64)
	cta, _ := strconv.ParseFloat(ctaStr, 64)

	res := c.service.Calculate(service.SREToilInput{
		HoursPerWeek:   hpw,
		HourlyRate:     hr,
		CostToAutomate: cta,
	})

	sentiment := SentimentGood
	roiColor := DefaultTheme.Success
	if res.AnnualSavings < 0 {
		sentiment = SentimentBad
		roiColor = DefaultTheme.Critical
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Primary).MarginBottom(1)
	roiStyle := lipgloss.NewStyle().Foreground(roiColor).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(DefaultTheme.TextNormal)

	str := fmt.Sprintf("%s\n\n%s %.1f\n%s %s",
		titleStyle.Render("󱔗 SRE TOIL ERADICATION"),
		labelStyle.Render("Total Annual Hours Saved:"), res.HoursSaved,
		labelStyle.Render("Annual Net Savings:"), roiStyle.Render(fmt.Sprintf("$%.2f", res.AnnualSavings)))

	return str, sentiment
}

func (c *SRECalculator) Reset() {
	c.hoursPerWeek = ""
	c.hourlyRate = ""
	c.costToAutomate = ""
}
