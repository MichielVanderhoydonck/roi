package tui

import (
	"fmt"
	"strconv"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"charm.land/huh/v2"
)

type ProductivityCalculator struct {
	service     *service.ProductivityService
	timeBefore  string
	timeAfter   string
	executions  string
	hourlyRate  string
	maintenance string
}

func NewProductivityCalculator() *ProductivityCalculator {
	return &ProductivityCalculator{
		service: service.NewProductivityService(),
	}
}

func (c *ProductivityCalculator) CreateForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("timeBefore").
				Title("Time spent on task BEFORE").
				Placeholder("e.g. 4h, 30m").
				Value(&c.timeBefore).
				Validate(validateDuration),
			huh.NewInput().
				Key("timeAfter").
				Title("Time spent on task AFTER").
				Placeholder("e.g. 5m, 10s").
				Value(&c.timeAfter).
				Validate(validateDuration),
			huh.NewInput().
				Key("executions").
				Title("Executions per year").
				Placeholder("e.g. 1000").
				Value(&c.executions).
				Validate(validateInt),
			huh.NewInput().
				Key("hourlyRate").
				Title("Average Developer Hourly Rate ($)").
				Placeholder("e.g. 75").
				Value(&c.hourlyRate).
				Validate(validateFloat),
			huh.NewInput().
				Key("maintenance").
				Title("Cost of Building/Maintaining Tool ($)").
				Placeholder("e.g. 1000").
				Value(&c.maintenance).
				Validate(validateFloat),
		),
	).WithShowHelp(false)
	return applyTheme(f)
}

func (c *ProductivityCalculator) GetContext(key string) string {
	help := map[string]string{
		"timeBefore":  "How long did this process take manually before your automation?\nFormat: Valid time units are \"h\", \"m\", \"s\". Examples: \"4h\", \"30m\".",
		"timeAfter":   "How long does the process take now with your internal developer platform or automation?\nFormat: Valid time units are \"h\", \"m\", \"s\".",
		"executions":  "How many times per year does your engineering team execute this task?",
		"hourlyRate":  "What is the fully loaded hourly cost of an engineer at your company?\nExample: 75",
		"maintenance": "How much does it cost annually in software licenses or engineering time to maintain this tool?\nExample: 500",
	}
	if val, ok := help[key]; ok {
		return val
	}
	return "Fill in the productivity details to calculate your ROI."
}

func (c *ProductivityCalculator) GetFormula(form *huh.Form) string {
	tb := getFormField(form, "timeBefore", c.timeBefore)
	ta := getFormField(form, "timeAfter", c.timeAfter)
	execs := getFormField(form, "executions", c.executions)
	hr := getFormField(form, "hourlyRate", c.hourlyRate)
	mc := getFormField(form, "maintenance", c.maintenance)

	return fmt.Sprintf(`Developer Productivity ROI (Time Saved)

Annual ROI ($) = 
  [(%s - %s) 
   × %s 
   × %s] 
  - %s`,
		formatFormulaValue(tb, "Time BEFORE"),
		formatFormulaValue(ta, "Time AFTER"),
		formatFormulaValue(execs, "Executions per year"),
		formatFormulaValue(hr, "Hourly Rate"),
		formatFormulaValue(mc, "Maintenance Cost"))
}

func (c *ProductivityCalculator) CalculateResult(form *huh.Form) (string, Sentiment) {
	tbStr := getFormField(form, "timeBefore", c.timeBefore)
	taStr := getFormField(form, "timeAfter", c.timeAfter)
	execsStr := getFormField(form, "executions", c.executions)
	hrStr := getFormField(form, "hourlyRate", c.hourlyRate)
	mcStr := getFormField(form, "maintenance", c.maintenance)

	if tbStr == "" || taStr == "" || execsStr == "" {
		return lipgloss.NewStyle().Foreground(DefaultTheme.TextDim).Render("Enter time values and executions to see ROI..."), SentimentNone
	}

	tb, _ := time.ParseDuration(tbStr)
	ta, _ := time.ParseDuration(taStr)
	execs, _ := strconv.Atoi(execsStr)
	hr, _ := strconv.ParseFloat(hrStr, 64)
	mc, _ := strconv.ParseFloat(mcStr, 64)

	res := c.service.Calculate(service.ProductivityInput{
		TimeBefore:        tb,
		TimeAfter:         ta,
		ExecutionsPerYear: execs,
		HourlyRate:        hr,
		MaintenanceCost:   mc,
	})

	sentiment := SentimentGood
	roiColor := DefaultTheme.Success
	if res.NetROI < 0 {
		sentiment = SentimentBad
		roiColor = DefaultTheme.Critical
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Primary).MarginBottom(1)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success).Bold(true)
	roiStyle := lipgloss.NewStyle().Foreground(roiColor).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(DefaultTheme.TextNormal)

	str := fmt.Sprintf("%s\n\n%s %s\n%s %s\n%s %s",
		titleStyle.Render("󰄬 PRODUCTIVITY IMPACT"),
		labelStyle.Render("Total Time Saved:"), valStyle.Render(res.TimeSaved.String()),
		labelStyle.Render("Gross Savings:"), valStyle.Render(fmt.Sprintf("$%.2f", res.GrossSavings)),
		labelStyle.Render("Net Annual ROI:"), roiStyle.Render(fmt.Sprintf("$%.2f", res.NetROI)))

	return str, sentiment
}

func (c *ProductivityCalculator) Reset() {
	c.timeBefore = ""
	c.timeAfter = ""
	c.executions = ""
	c.hourlyRate = ""
	c.maintenance = ""
}
