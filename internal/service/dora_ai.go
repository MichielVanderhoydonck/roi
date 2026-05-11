package service

// DORAAIInput holds the parameters for calculating DORA AI ROI.
type DORAAIInput struct {
	StaffSize                 float64
	Salary                    float64
	Revenue                   float64
	DowntimeCostPerHour       float64
	CurrentDeploymentsPerYear float64
	CurrentFeaturesPerYear    float64
	IdeaSuccessRate           float64
	RevenueImpactPerFeature   float64
	CurrentCfr                float64
	CurrentFdrt               float64
	TimeSavedPerDeveloper     float64
	AILicenseCostPerUser      float64
	AdditionalAICostPerUser   float64
	AdditionalAIInfraCost     float64
	TrainingCostPerUser       float64
	TargetDeploymentsPerYear  float64
	TargetFeaturesPerYear     float64
	TargetCfr                 float64
	JCurveDrop                float64
	JCurveDuration            float64
}

// DORAAIResult holds the output of the DORA AI ROI calculation.
type DORAAIResult struct {
	TotalHardCosts                float64
	JCurveCost                    float64
	TotalFirstYearInvestment      float64
	HeadcountReinvestmentCapacity float64
	RevenueFromExtraFeatures      float64
	DowntimeSavings               float64
	TotalFirstYearValue           float64
	FirstYearBenefit              float64
	ROI                           float64
	PaybackPeriod                 float64
}

// DORAAICalculator defines the interface for calculating DORA AI ROI.
type DORAAICalculator interface {
	Calculate(input DORAAIInput) DORAAIResult
}

// DORAAIService implements the DORAAICalculator interface.
type DORAAIService struct{}

// NewDORAAIService creates a new DORAAIService.
func NewDORAAIService() *DORAAIService {
	return &DORAAIService{}
}

// Calculate computes the DORA AI ROI based on the official metrics.
func (s *DORAAIService) Calculate(input DORAAIInput) DORAAIResult {
	totalHardCosts := (input.AILicenseCostPerUser+input.AdditionalAICostPerUser+input.TrainingCostPerUser)*input.StaffSize + input.AdditionalAIInfraCost
	jCurveCost := input.StaffSize * input.Salary * input.JCurveDrop * (input.JCurveDuration / 12.0)
	totalFirstYearInvestment := totalHardCosts + jCurveCost

	headcountReinvestmentCapacity := input.StaffSize * input.Salary * input.TimeSavedPerDeveloper
	revenueFromExtraFeatures := (input.TargetFeaturesPerYear - input.CurrentFeaturesPerYear) * input.IdeaSuccessRate * input.RevenueImpactPerFeature * input.Revenue
	downtimeSavings := input.CurrentDeploymentsPerYear*input.CurrentCfr*input.CurrentFdrt*input.DowntimeCostPerHour - input.TargetDeploymentsPerYear*input.TargetCfr*input.CurrentFdrt*input.DowntimeCostPerHour

	totalFirstYearValue := headcountReinvestmentCapacity + revenueFromExtraFeatures + downtimeSavings
	firstYearBenefit := totalFirstYearValue - totalFirstYearInvestment

	roi := 0.0
	if totalFirstYearInvestment != 0 {
		roi = firstYearBenefit / totalFirstYearInvestment
	}

	paybackPeriod := 0.0
	if totalFirstYearValue > 0 {
		paybackPeriod = totalFirstYearInvestment / totalFirstYearValue
	}

	return DORAAIResult{
		TotalHardCosts:                totalHardCosts,
		JCurveCost:                    jCurveCost,
		TotalFirstYearInvestment:      totalFirstYearInvestment,
		HeadcountReinvestmentCapacity: headcountReinvestmentCapacity,
		RevenueFromExtraFeatures:      revenueFromExtraFeatures,
		DowntimeSavings:               downtimeSavings,
		TotalFirstYearValue:           totalFirstYearValue,
		FirstYearBenefit:              firstYearBenefit,
		ROI:                           roi,
		PaybackPeriod:                 paybackPeriod,
	}
}
