package tui

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"charm.land/huh/v2"
)

type FinOpsCalculator struct {
	service *service.FinOpsService
	oldBill string
	newBill string
}

func NewFinOpsCalculator() *FinOpsCalculator {
	return &FinOpsCalculator{
		service: service.NewFinOpsService(),
	}
}

func (c *FinOpsCalculator) CreateForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("oldBill").
				Title("Previous Monthly Cloud Bill ($)").
				Placeholder("e.g. 20000").
				Value(&c.oldBill).
				Validate(validateFloat),
			huh.NewInput().
				Key("newBill").
				Title("New Monthly Cloud Bill ($)").
				Placeholder("e.g. 15000").
				Value(&c.newBill).
				Validate(validateFloat),
		),
	).WithShowHelp(false)
	return applyTheme(f)
}

func (c *FinOpsCalculator) GetContext(key string) string {
	help := map[string]string{
		"oldBill": "Your average monthly cloud infrastructure bill before optimization.",
		"newBill": "Your target or actual monthly cloud bill after right-sizing or spinning down resources.",
	}
	if val, ok := help[key]; ok {
		return val
	}
	return "Fill in your cloud bills to calculate infrastructure savings."
}

func (c *FinOpsCalculator) GetFormula(form *huh.Form) string {
	ob := getFormField(form, "oldBill", c.oldBill)
	nb := getFormField(form, "newBill", c.newBill)

	return fmt.Sprintf(`FinOps ROI (Infrastructure Optimization)

Cloud Savings ($) = 
  (%s - %s) 
  × 12 months`,
		formatFormulaValue(ob, "Previous Monthly Bill"),
		formatFormulaValue(nb, "New Monthly Bill"))
}

func (c *FinOpsCalculator) CalculateResult(form *huh.Form) (string, Sentiment) {
	obStr := getFormField(form, "oldBill", c.oldBill)
	nbStr := getFormField(form, "newBill", c.newBill)

	if obStr == "" || nbStr == "" {
		return "", SentimentNone
	}

	ob, _ := strconv.ParseFloat(obStr, 64)
	nb, _ := strconv.ParseFloat(nbStr, 64)

	res := c.service.Calculate(service.FinOpsInput{
		OldMonthlyBill: ob,
		NewMonthlyBill: nb,
	})

	sentiment := SentimentGood
	roiColor := DefaultTheme.Success
	if ob <= nb {
		sentiment = SentimentBad
		roiColor = DefaultTheme.Critical
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Secondary).MarginBottom(1)
	roiStyle := lipgloss.NewStyle().Foreground(roiColor).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(DefaultTheme.TextNormal)

	str := fmt.Sprintf("%s\n\n%s %s\n%s %s",
		titleStyle.Render("󰠩 FINOPS OPTIMIZATION"),
		labelStyle.Render("Monthly Savings:"), roiStyle.Render(fmt.Sprintf("$%.2f", ob-nb)),
		labelStyle.Render("Annual Cloud Savings:"), roiStyle.Render(fmt.Sprintf("$%.2f", res.AnnualSavings)))

	return str, sentiment
}

func (c *FinOpsCalculator) Reset() {
	c.oldBill = ""
	c.newBill = ""
}
