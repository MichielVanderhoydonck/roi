package service

import (
	"testing"
)

func TestDORAAICalculate(t *testing.T) {
	service := NewDORAAIService()

	// Default inputs from the DORA ROI calculator
	input := DORAAIInput{
		StaffSize:                 500,
		Salary:                    176000,
		Revenue:                   100000000,
		DowntimeCostPerHour:       100000,
		CurrentDeploymentsPerYear: 50,
		CurrentFeaturesPerYear:    50,
		IdeaSuccessRate:           0.33,
		RevenueImpactPerFeature:   0.005,
		CurrentCfr:                0.05,
		CurrentFdrt:               4,
		TimeSavedPerDeveloper:     0.125,
		AILicenseCostPerUser:      250,
		AdditionalAICostPerUser:   80,
		AdditionalAIInfraCost:     100000,
		TrainingCostPerUser:       9600,
		TargetDeploymentsPerYear:  56,
		TargetFeaturesPerYear:     56,
		TargetCfr:                 0.06,
		JCurveDrop:                0.15,
		JCurveDuration:            3,
	}

	result := service.Calculate(input)

	// We calculate expected manually based on formulas or JS logic
	// Hard Costs = (250+80+9600) * 500 + 100000 = 9930 * 500 + 100000 = 4965000 + 100000 = 5065000
	if result.TotalHardCosts != 5065000 {
		t.Errorf("Expected TotalHardCosts 5065000, got %f", result.TotalHardCosts)
	}

	// JCurve Cost = 500 * 176000 * 0.15 * (3/12) = 88000000 * 0.15 * 0.25 = 13200000 * 0.25 = 3300000
	if result.JCurveCost != 3300000 {
		t.Errorf("Expected JCurveCost 3300000, got %f", result.JCurveCost)
	}

	if result.TotalFirstYearInvestment != 8365000 {
		t.Errorf("Expected TotalFirstYearInvestment 8365000, got %f", result.TotalFirstYearInvestment)
	}

	// Headcount reinvestment = 500 * 176000 * 0.125 = 88000000 * 0.125 = 11000000
	if result.HeadcountReinvestmentCapacity != 11000000 {
		t.Errorf("Expected HeadcountReinvestmentCapacity 11000000, got %f", result.HeadcountReinvestmentCapacity)
	}

	// Revenue extra features = (56 - 50) * 0.33 * 0.005 * 100000000 = 6 * 0.33 * 500000 = 1.98 * 500000 = 990000
	if int(result.RevenueFromExtraFeatures) != 990000 {
		t.Errorf("Expected RevenueFromExtraFeatures 990000, got %f", result.RevenueFromExtraFeatures)
	}

	// Downtime savings = (50*0.05*4*100000) - (56*0.06*4*100000) = (10*100000) - (13.44*100000) = 1000000 - 1344000 = -344000
	// Because of float precision, it might not be exactly -344000, let's round or cast to int.
	if int(result.DowntimeSavings) != -344000 {
		t.Errorf("Expected DowntimeSavings -344000, got %d", int(result.DowntimeSavings))
	}

	expectedTotalFirstYearValue := 11000000 + 990000 - 344000 // 11646000
	if int(result.TotalFirstYearValue) != expectedTotalFirstYearValue {
		t.Errorf("Expected TotalFirstYearValue %d, got %d", expectedTotalFirstYearValue, int(result.TotalFirstYearValue))
	}
}

func TestDORAAICalculateZeroes(t *testing.T) {
	service := NewDORAAIService()

	input := DORAAIInput{}
	result := service.Calculate(input)

	if result.TotalHardCosts != 0 {
		t.Errorf("Expected TotalHardCosts 0, got %f", result.TotalHardCosts)
	}
	if result.JCurveCost != 0 {
		t.Errorf("Expected JCurveCost 0, got %f", result.JCurveCost)
	}
	if result.TotalFirstYearInvestment != 0 {
		t.Errorf("Expected TotalFirstYearInvestment 0, got %f", result.TotalFirstYearInvestment)
	}
	if result.TotalFirstYearValue != 0 {
		t.Errorf("Expected TotalFirstYearValue 0, got %f", result.TotalFirstYearValue)
	}
	if result.ROI != 0 {
		t.Errorf("Expected ROI 0, got %f", result.ROI)
	}
	if result.PaybackPeriod != 0 {
		t.Errorf("Expected PaybackPeriod 0, got %f", result.PaybackPeriod)
	}
}
