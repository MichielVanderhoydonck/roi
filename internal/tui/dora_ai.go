package tui

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"charm.land/huh/v2"
)

type DORAAICalculator struct {
	service            *service.DORAAIService
	staffSize          string
	salary             string
	revenue            string
	downtimeCost       string
	currentDeployments string
	currentFeatures    string
	ideaSuccessRate    string
	revenueImpact      string
	currentCfr         string
	currentFdrt        string
	timeSaved          string
	licenseCost        string
	addAICost          string
	infraCost          string
	trainingCost       string
	targetDeployments  string
	targetFeatures     string
	targetCfr          string
	jCurveDrop         string
	jCurveDuration     string
}

func NewDORAAICalculator() *DORAAICalculator {
	return &DORAAICalculator{
		service: service.NewDORAAIService(),
	}
}

func (c *DORAAICalculator) CreateForm() *huh.Form {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Key("staffSize").Title("Technical staff size").Placeholder("500").Value(&c.staffSize).Validate(validateFloat),
			huh.NewInput().Key("salary").Title("Average technical staff salary").Placeholder("176000").Value(&c.salary).Validate(validateFloat),
			huh.NewInput().Key("revenue").Title("Product portfolio revenue").Placeholder("100000000").Value(&c.revenue).Validate(validateFloat),
			huh.NewInput().Key("downtimeCost").Title("Cost of downtime per hour").Placeholder("100000").Value(&c.downtimeCost).Validate(validateFloat),
		).Title("Organizational metrics"),
		huh.NewGroup(
			huh.NewInput().Key("currentDeployments").Title("Current deployments per year").Placeholder("50").Value(&c.currentDeployments).Validate(validateFloat),
			huh.NewInput().Key("currentFeatures").Title("Current features deployed per year").Placeholder("50").Value(&c.currentFeatures).Validate(validateFloat),
			huh.NewInput().Key("ideaSuccessRate").Title("Idea success rate (e.g. 0.33 for 33%)").Placeholder("0.33").Value(&c.ideaSuccessRate).Validate(validateFloat),
			huh.NewInput().Key("revenueImpact").Title("Average revenue impact per feature").Placeholder("0.005").Value(&c.revenueImpact).Validate(validateFloat),
			huh.NewInput().Key("currentCfr").Title("Current change failure rate").Placeholder("0.05").Value(&c.currentCfr).Validate(validateFloat),
			huh.NewInput().Key("currentFdrt").Title("Failed deployment recovery time (hours)").Placeholder("4").Value(&c.currentFdrt).Validate(validateFloat),
		).Title("Baseline software delivery metrics"),
		huh.NewGroup(
			huh.NewInput().Key("timeSaved").Title("Net time saved per developer").Placeholder("0.125").Value(&c.timeSaved).Validate(validateFloat),
			huh.NewInput().Key("licenseCost").Title("Annual AI license cost per user").Placeholder("250").Value(&c.licenseCost).Validate(validateFloat),
			huh.NewInput().Key("addAICost").Title("Additional annual AI costs per user").Placeholder("80").Value(&c.addAICost).Validate(validateFloat),
			huh.NewInput().Key("infraCost").Title("Additional annual AI infra costs").Placeholder("100000").Value(&c.infraCost).Validate(validateFloat),
			huh.NewInput().Key("trainingCost").Title("Annual training costs per user").Placeholder("9600").Value(&c.trainingCost).Validate(validateFloat),
			huh.NewInput().Key("targetDeployments").Title("Target number of deployments per year").Placeholder("56").Value(&c.targetDeployments).Validate(validateFloat),
			huh.NewInput().Key("targetFeatures").Title("Target number of features per year").Placeholder("56").Value(&c.targetFeatures).Validate(validateFloat),
			huh.NewInput().Key("targetCfr").Title("Target change failure rate").Placeholder("0.06").Value(&c.targetCfr).Validate(validateFloat),
			huh.NewInput().Key("jCurveDrop").Title("J-Curve productivity drop").Placeholder("0.15").Value(&c.jCurveDrop).Validate(validateFloat),
			huh.NewInput().Key("jCurveDuration").Title("J-Curve productivity drop timeline (months)").Placeholder("3").Value(&c.jCurveDuration).Validate(validateFloat),
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

func (c *DORAAICalculator) CalculateResult(form *huh.Form) (string, Sentiment) {
	fallbacks := map[string]string{
		"staffSize":          c.staffSize,
		"salary":             c.salary,
		"revenue":            c.revenue,
		"downtimeCost":       c.downtimeCost,
		"currentDeployments": c.currentDeployments,
		"currentFeatures":    c.currentFeatures,
		"ideaSuccessRate":    c.ideaSuccessRate,
		"revenueImpact":      c.revenueImpact,
		"currentCfr":         c.currentCfr,
		"currentFdrt":        c.currentFdrt,
		"timeSaved":          c.timeSaved,
		"licenseCost":        c.licenseCost,
		"addAICost":          c.addAICost,
		"infraCost":          c.infraCost,
		"trainingCost":       c.trainingCost,
		"targetDeployments":  c.targetDeployments,
		"targetFeatures":     c.targetFeatures,
		"targetCfr":          c.targetCfr,
		"jCurveDrop":         c.jCurveDrop,
		"jCurveDuration":     c.jCurveDuration,
	}

	getFloat := func(key string) float64 {
		s := getFormField(form, key, fallbacks[key])
		if s == "" {
			return 0
		}
		val, _ := strconv.ParseFloat(s, 64)
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

	// Only calculate if we have at least some basic inputs
	if input.StaffSize == 0 || input.Salary == 0 {
		return "", SentimentNone
	}

	res := c.service.Calculate(input)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Primary).MarginBottom(1)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(DefaultTheme.TextNormal)
	statStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Secondary)

	paybackStr := "N/A"
	if res.PaybackPeriod > 0 {
		paybackStr = fmt.Sprintf("%.1f years", res.PaybackPeriod)
	}

	sentiment := SentimentGood
	roiColor := DefaultTheme.Success
	if res.ROI < 0 {
		sentiment = SentimentBad
		roiColor = DefaultTheme.Critical
	}
	roiStyle := lipgloss.NewStyle().Foreground(roiColor).Bold(true)

	str := fmt.Sprintf("%s\n\n%s %s\n%s %s\n%s %s\n\n%s\n%s $%.0f\n%s $%.0f",
		titleStyle.Render("󰄬 DORA AI IMPACT ANALYSIS"),
		labelStyle.Render("Return on Investment:"), roiStyle.Render(fmt.Sprintf("%.1f%%", res.ROI*100)),
		labelStyle.Render("Payback Period:"), valStyle.Render(paybackStr),
		labelStyle.Render("1st Year Net Benefit:"), valStyle.Render(fmt.Sprintf("$%.0f", res.FirstYearBenefit)),
		statStyle.Render("Financial Breakdown:"),
		labelStyle.Render("  Total Investment:"), res.TotalFirstYearInvestment,
		labelStyle.Render("  Total Gross Value:"), res.TotalFirstYearValue,
	)

	return str, sentiment
}

func (c *DORAAICalculator) Reset() {
	c.staffSize = ""
	c.salary = ""
	c.revenue = ""
	c.downtimeCost = ""
	c.currentDeployments = ""
	c.currentFeatures = ""
	c.ideaSuccessRate = ""
	c.revenueImpact = ""
	c.currentCfr = ""
	c.currentFdrt = ""
	c.timeSaved = ""
	c.licenseCost = ""
	c.addAICost = ""
	c.infraCost = ""
	c.trainingCost = ""
	c.targetDeployments = ""
	c.targetFeatures = ""
	c.targetCfr = ""
	c.jCurveDrop = ""
	c.jCurveDuration = ""
}
