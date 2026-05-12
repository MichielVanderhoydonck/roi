package tui

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"charm.land/huh/v2"
)

type OnboardingCalculator struct {
	service   *service.OnboardingService
	oldDays   string
	newDays   string
	newHires  string
	dailyRate string
}

func NewOnboardingCalculator() *OnboardingCalculator {
	return &OnboardingCalculator{
		service: service.NewOnboardingService(),
	}
}

func (c *OnboardingCalculator) CreateForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("oldDays").
				Title("Old Days to First Deploy").
				Placeholder("e.g. 10").
				Value(&c.oldDays).
				Validate(validateFloat),
			huh.NewInput().
				Key("newDays").
				Title("New Days to First Deploy").
				Placeholder("e.g. 2").
				Value(&c.newDays).
				Validate(validateFloat),
			huh.NewInput().
				Key("newHires").
				Title("Number of New Hires per Year").
				Placeholder("e.g. 20").
				Value(&c.newHires).
				Validate(validateInt),
			huh.NewInput().
				Key("dailyRate").
				Title("Daily Developer Rate ($)").
				Placeholder("e.g. 600").
				Value(&c.dailyRate).
				Validate(validateFloat),
		),
	).WithShowHelp(false)
	return applyTheme(f)
}

func (c *OnboardingCalculator) GetContext(key string) string {
	help := map[string]string{
		"oldDays":   "Historically, how many days did it take for a new hire to make their first production deployment?\nExample: 10",
		"newDays":   "With your improvements (e.g. automated environments), how many days does it take now?\nExample: 2",
		"newHires":  "How many engineers do you plan to hire this year?\nExample: 20",
		"dailyRate": "What is the fully loaded daily cost of an engineer?\nExample: 600",
	}
	if val, ok := help[key]; ok {
		return val
	}
	return "Fill in the onboarding details to calculate your ROI."
}

func (c *OnboardingCalculator) GetFormula(form *huh.Form) string {
	od := getFormField(form, "oldDays", c.oldDays)
	nd := getFormField(form, "newDays", c.newDays)
	nh := getFormField(form, "newHires", c.newHires)
	dr := getFormField(form, "dailyRate", c.dailyRate)

	return fmt.Sprintf(`Onboarding ROI (Time to First Value)

Onboarding Savings ($) = 
  (%s days - %s days) 
  × %s hires 
  × %s`,
		formatFormulaValue(od, "Old Days"),
		formatFormulaValue(nd, "New Days"),
		formatFormulaValue(nh, "New Hires"),
		formatFormulaValue(dr, "Daily Rate"))
}

func (c *OnboardingCalculator) CalculateResult(form *huh.Form) (string, Sentiment) {
	odStr := getFormField(form, "oldDays", c.oldDays)
	ndStr := getFormField(form, "newDays", c.newDays)
	nhStr := getFormField(form, "newHires", c.newHires)
	drStr := getFormField(form, "dailyRate", c.dailyRate)

	if odStr == "" || ndStr == "" || nhStr == "" {
		return "", SentimentNone
	}

	od, _ := strconv.ParseFloat(odStr, 64)
	nd, _ := strconv.ParseFloat(ndStr, 64)
	nh, _ := strconv.Atoi(nhStr)
	dr, _ := strconv.ParseFloat(drStr, 64)

	res := c.service.Calculate(service.OnboardingInput{
		OldDays:   od,
		NewDays:   nd,
		NewHires:  nh,
		DailyRate: dr,
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
		titleStyle.Render("󱔗 ONBOARDING EFFICIENCY"),
		labelStyle.Render("Days Saved per Hire:"), res.DaysSavedPerHire,
		labelStyle.Render("Annual Idle Time Savings:"), roiStyle.Render(fmt.Sprintf("$%.2f", res.AnnualSavings)))

	return str, sentiment
}

func (c *OnboardingCalculator) Reset() {
	c.oldDays = ""
	c.newDays = ""
	c.newHires = ""
	c.dailyRate = ""
}
