package tui

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"charm.land/huh/v2"
)

type CostOfDelayCalculator struct {
	service        *service.CostOfDelayService
	monthlyRevenue string
	daysDelayed    string
}

func NewCostOfDelayCalculator() *CostOfDelayCalculator {
	return &CostOfDelayCalculator{
		service: service.NewCostOfDelayService(),
	}
}

func (c *CostOfDelayCalculator) CreateForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("monthlyRevenue").
				Title("Estimated Monthly Revenue of Feature ($)").
				Placeholder("e.g. 300000").
				Value(&c.monthlyRevenue).
				Validate(validateFloat),
			huh.NewInput().
				Key("daysDelayed").
				Title("Days Delayed").
				Placeholder("e.g. 15").
				Value(&c.daysDelayed).
				Validate(validateFloat),
		),
	).WithShowHelp(false)
	return applyTheme(f)
}

func (c *CostOfDelayCalculator) GetContext(key string) string {
	help := map[string]string{
		"monthlyRevenue": "What is the anticipated monthly revenue this feature will generate?\nExample: 300000",
		"daysDelayed":    "How many days was the launch delayed due to bottlenecks?\nExample: 15",
	}
	if val, ok := help[key]; ok {
		return val
	}
	return "Fill in the details to calculate the cost of delay."
}

func (c *CostOfDelayCalculator) GetFormula(form *huh.Form) string {
	mr := getFormField(form, "monthlyRevenue", c.monthlyRevenue)
	dd := getFormField(form, "daysDelayed", c.daysDelayed)

	return fmt.Sprintf(`Cost of Delay

Revenue Lost ($) = 
  (%s / 30) 
  × %s days`,
		formatFormulaValue(mr, "Monthly Revenue"),
		formatFormulaValue(dd, "Days Delayed"))
}

func (c *CostOfDelayCalculator) CalculateResult(form *huh.Form) (string, Sentiment) {
	mrStr := getFormField(form, "monthlyRevenue", c.monthlyRevenue)
	ddStr := getFormField(form, "daysDelayed", c.daysDelayed)

	if mrStr == "" || ddStr == "" {
		return "", SentimentNone
	}

	mr, _ := strconv.ParseFloat(mrStr, 64)
	dd, _ := strconv.ParseFloat(ddStr, 64)

	res := c.service.Calculate(service.CostOfDelayInput{
		EstimatedMonthlyRevenue: mr,
		DaysDelayed:             dd,
	})

	sentiment := SentimentGood
	roiColor := DefaultTheme.Success
	if res.CostOfDelay > 0 {
		sentiment = SentimentBad
		roiColor = DefaultTheme.Critical
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Primary).MarginBottom(1)
	roiStyle := lipgloss.NewStyle().Foreground(roiColor).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(DefaultTheme.TextNormal)

	str := fmt.Sprintf("%s\n\n%s %s",
		titleStyle.Render("󰚌 COST OF DELAY IMPACT"),
		labelStyle.Render("Total Revenue Lost:"), roiStyle.Render(fmt.Sprintf("$%.2f", res.CostOfDelay)))

	return str, sentiment
}

func (c *CostOfDelayCalculator) Reset() {
	c.monthlyRevenue = ""
	c.daysDelayed = ""
}
