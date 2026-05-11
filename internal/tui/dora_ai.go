package tui

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"github.com/charmbracelet/huh"
)

type DORAAICalculator struct {
	service *service.DORAAIService
}

func NewDORAAICalculator() *DORAAICalculator {
	return &DORAAICalculator{service: service.NewDORAAIService()}
}

func (c *DORAAICalculator) CreateForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Key("staffSize").Title("Technical staff size").Placeholder("500").Validate(validateFloat),
			huh.NewInput().Key("salary").Title("Average technical staff salary").Placeholder("176000").Validate(validateFloat),
			huh.NewInput().Key("revenue").Title("Product portfolio revenue").Placeholder("100000000").Validate(validateFloat),
			huh.NewInput().Key("downtimeCost").Title("Cost of downtime per hour").Placeholder("100000").Validate(validateFloat),
		).Title("Organizational metrics"),
		huh.NewGroup(
			huh.NewInput().Key("currentDeployments").Title("Current deployments per year").Placeholder("50").Validate(validateFloat),
			huh.NewInput().Key("currentFeatures").Title("Current features deployed per year").Placeholder("50").Validate(validateFloat),
			huh.NewInput().Key("ideaSuccessRate").Title("Idea success rate (e.g. 0.33 for 33%)").Placeholder("0.33").Validate(validateFloat),
			huh.NewInput().Key("revenueImpact").Title("Average revenue impact per feature").Placeholder("0.005").Validate(validateFloat),
			huh.NewInput().Key("currentCfr").Title("Current change failure rate").Placeholder("0.05").Validate(validateFloat),
			huh.NewInput().Key("currentFdrt").Title("Failed deployment recovery time (hours)").Placeholder("4").Validate(validateFloat),
		).Title("Baseline software delivery metrics"),
		huh.NewGroup(
			huh.NewInput().Key("timeSaved").Title("Net time saved per developer").Placeholder("0.125").Validate(validateFloat),
			huh.NewInput().Key("licenseCost").Title("Annual AI license cost per user").Placeholder("250").Validate(validateFloat),
			huh.NewInput().Key("addAICost").Title("Additional annual AI costs per user").Placeholder("80").Validate(validateFloat),
			huh.NewInput().Key("infraCost").Title("Additional annual AI infra costs").Placeholder("100000").Validate(validateFloat),
			huh.NewInput().Key("trainingCost").Title("Annual training costs per user").Placeholder("9600").Validate(validateFloat),
			huh.NewInput().Key("targetDeployments").Title("Target number of deployments per year").Placeholder("56").Validate(validateFloat),
			huh.NewInput().Key("targetFeatures").Title("Target number of features per year").Placeholder("56").Validate(validateFloat),
			huh.NewInput().Key("targetCfr").Title("Target change failure rate").Placeholder("0.06").Validate(validateFloat),
			huh.NewInput().Key("jCurveDrop").Title("J-Curve productivity drop").Placeholder("0.15").Validate(validateFloat),
			huh.NewInput().Key("jCurveDuration").Title("J-Curve productivity drop timeline (months)").Placeholder("3").Validate(validateFloat),
		).Title("AI estimates"),
	).WithShowHelp(false)
	return applyTheme(f)
}

func (c *DORAAICalculator) GetContext(key string) string {
	help := map[string]string{
		"staffSize":          "The count of full-time employees involved in the SDLC.",
		"salary":             "Average fully loaded technical staff salary.",
		"revenue":            "The annual revenue driven by this software.",
		"downtimeCost":       "An estimate of the cost of one hour of system outage.",
		"currentDeployments": "The total number of deployments per year.",
		"currentFeatures":    "The number of features deployed each year.",
		"ideaSuccessRate":    "The percentage of deployed features that increase revenue (e.g. 0.33).",
		"revenueImpact":      "The average revenue increase for each successful feature (e.g. 0.005).",
		"currentCfr":         "The current change failure rate.",
		"currentFdrt":        "The average hours to restore service after a change failure.",
		"timeSaved":          "Net productivity boost per developer (e.g. 0.125).",
		"licenseCost":        "Annual price per user of an AI subscription.",
		"addAICost":          "Additional per-user costs like API or token costs.",
		"infraCost":          "New AI-related infrastructure costs.",
		"trainingCost":       "Training and enablement costs for each employee.",
		"targetDeployments":  "Expected deployments per year while using AI.",
		"targetFeatures":     "Expected features deployed per year.",
		"targetCfr":          "Target percentage of changes resulting in degraded service.",
		"jCurveDrop":         "Temporary productivity decrease during the AI learning phase.",
		"jCurveDuration":     "Length of the productivity decrease in months.",
	}
	if val, ok := help[key]; ok {
		return val
	}
	return "Fill in the DORA AI metrics to calculate your ROI."
}

func (c *DORAAICalculator) GetFormula(form *huh.Form) string {
	return `DORA ROI of AI-assisted Software Development

First Year Benefit = Total First Year Value - Total First Year Investment
ROI = First Year Benefit / Total First Year Investment

Total First Year Investment includes:
- Hard costs (Tooling and training)
- J-Curve cost

Total First Year Value includes:
- Headcount reinvestment capacity
- Revenue from extra feature deployments
- Downtime impact (savings/costs)`
}

func (c *DORAAICalculator) CalculateResult(form *huh.Form) string {
	getFloat := func(key string) float64 {
		val, _ := strconv.ParseFloat(form.GetString(key), 64)
		return val
	}

	input := service.DORAAIInput{
		StaffSize:                 getFloat("staffSize"),
		Salary:                    getFloat("salary"),
		Revenue:                   getFloat("revenue"),
		DowntimeCostPerHour:       getFloat("downtimeCost"),
		CurrentDeploymentsPerYear: getFloat("currentDeployments"),
		CurrentFeaturesPerYear:    getFloat("currentFeatures"),
		IdeaSuccessRate:           getFloat("ideaSuccessRate"),
		RevenueImpactPerFeature:   getFloat("revenueImpact"),
		CurrentCfr:                getFloat("currentCfr"),
		CurrentFdrt:               getFloat("currentFdrt"),
		TimeSavedPerDeveloper:     getFloat("timeSaved"),
		AILicenseCostPerUser:      getFloat("licenseCost"),
		AdditionalAICostPerUser:   getFloat("addAICost"),
		AdditionalAIInfraCost:     getFloat("infraCost"),
		TrainingCostPerUser:       getFloat("trainingCost"),
		TargetDeploymentsPerYear:  getFloat("targetDeployments"),
		TargetFeaturesPerYear:     getFloat("targetFeatures"),
		TargetCfr:                 getFloat("targetCfr"),
		JCurveDrop:                getFloat("jCurveDrop"),
		JCurveDuration:            getFloat("jCurveDuration"),
	}

	res := c.service.Calculate(input)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Primary)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success)

	paybackStr := "N/A"
	if res.PaybackPeriod != 0 {
		paybackStr = fmt.Sprintf("%.1f years", res.PaybackPeriod)
	}

	return fmt.Sprintf("%s\n\nReturn on investment (ROI): %s\nPayback period: %s\nFirst year benefit: %s\n\nTotal first year investment: $%.0f\nTotal first year value: $%.0f",
		titleStyle.Render("=== Calculated ROI ==="),
		valStyle.Render(fmt.Sprintf("%.1f%%", res.ROI*100)),
		valStyle.Render(paybackStr),
		valStyle.Render(fmt.Sprintf("$%.0f", res.FirstYearBenefit)),
		res.TotalFirstYearInvestment,
		res.TotalFirstYearValue,
	)
}
