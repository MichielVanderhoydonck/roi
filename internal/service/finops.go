package service

import (
)

// FinOpsInput holds the parameters for calculating FinOps ROI.
type FinOpsInput struct {
	OldMonthlyBill float64
	NewMonthlyBill float64
}

// FinOpsResult holds the output of the FinOps ROI calculation.
type FinOpsResult struct {
	AnnualSavings float64
}

// FinOpsCalculator defines the interface for calculating FinOps ROI.
type FinOpsCalculator interface {
	Calculate(input FinOpsInput) FinOpsResult
}

// FinOpsService implements the FinOpsCalculator interface.
type FinOpsService struct{}

// NewFinOpsService creates a new FinOpsService.
func NewFinOpsService() *FinOpsService {
	return &FinOpsService{}
}

// Calculate computes the FinOps ROI.
// Cloud Savings = (Old Monthly Bill - New Monthly Bill) * 12
func (s *FinOpsService) Calculate(input FinOpsInput) FinOpsResult {
	monthlySavings := input.OldMonthlyBill - input.NewMonthlyBill
	if monthlySavings < 0 {
		monthlySavings = 0
	}

	annualSavings := monthlySavings * 12

	return FinOpsResult{
		AnnualSavings: annualSavings,
	}
}
