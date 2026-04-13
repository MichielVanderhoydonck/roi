package tui

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"github.com/charmbracelet/huh"
)

func createOnboardingForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("oldDays").
				Title("Old Days to First Deploy").
				Placeholder("e.g. 10").
				Validate(validateFloat),
			huh.NewInput().
				Key("newDays").
				Title("New Days to First Deploy").
				Placeholder("e.g. 2").
				Validate(validateFloat),
			huh.NewInput().
				Key("newHires").
				Title("Number of New Hires per Year").
				Placeholder("e.g. 20").
				Validate(validateInt),
			huh.NewInput().
				Key("dailyRate").
				Title("Daily Developer Rate ($)").
				Placeholder("e.g. 600").
				Validate(validateFloat),
		),
	).WithShowHelp(false)
	return applyTheme(f)
}

func getOnboardingContext(key string) string {
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

func getOnboardingFormula(form *huh.Form) string {
	var od, nd, nh, dr string
	if form != nil {
		od = form.GetString("oldDays")
		nd = form.GetString("newDays")
		nh = form.GetString("newHires")
		dr = form.GetString("dailyRate")
	}

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

func (a *App) calcOnboardingResult() {
	od, _ := strconv.ParseFloat(a.onboardingForm.GetString("oldDays"), 64)
	nd, _ := strconv.ParseFloat(a.onboardingForm.GetString("newDays"), 64)
	nh, _ := strconv.Atoi(a.onboardingForm.GetString("newHires"))
	dr, _ := strconv.ParseFloat(a.onboardingForm.GetString("dailyRate"), 64)

	res := a.onboardingService.Calculate(service.OnboardingInput{
		OldDays:   od,
		NewDays:   nd,
		NewHires:  nh,
		DailyRate: dr,
	})

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Primary)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success)
	a.resultText = fmt.Sprintf("%s\n\nDays Saved per Hire: %.1f\nIdle Time Savings:   %s",
		titleStyle.Render("=== Onboarding ROI Results ==="),
		res.DaysSavedPerHire,
		valStyle.Render(fmt.Sprintf("$%.2f", res.AnnualSavings)))
}
