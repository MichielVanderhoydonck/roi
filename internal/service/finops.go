package service

import (
	"github.com/MichielVanderhoydonck/roi/internal/domain"
)

// FinOpsService implements the FinOpsCalculator interface.
type FinOpsService struct{}

// NewFinOpsService creates a new FinOpsService.
func NewFinOpsService() *FinOpsService {
	return &FinOpsService{}
}

// Calculate computes the FinOps ROI.
// Cloud Savings = (Old Monthly Bill - New Monthly Bill) * 12
func (s *FinOpsService) Calculate(input domain.FinOpsInput) domain.FinOpsResult {
	monthlySavings := input.OldMonthlyBill - input.NewMonthlyBill
	if monthlySavings < 0 {
		monthlySavings = 0
	}

	annualSavings := monthlySavings * 12

	return domain.FinOpsResult{
		AnnualSavings: annualSavings,
	}
}
