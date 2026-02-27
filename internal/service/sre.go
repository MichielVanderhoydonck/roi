package service

import (
	"github.com/MichielVanderhoydonck/roi/internal/domain"
)

// SREToilService implements the SREToilCalculator interface.
type SREToilService struct{}

// NewSREToilService creates a new SREToilService.
func NewSREToilService() *SREToilService {
	return &SREToilService{}
}

// Calculate computes the SRE Toil Eradication ROI.
// Toil Savings ($) = [(Hours of Toil per Week × 52) × (Hourly Dev Rate)] - Cost to Automate
func (s *SREToilService) Calculate(input domain.SREToilInput) domain.SREToilResult {
	annualToilHours := input.HoursPerWeek * 52
	grossSavings := annualToilHours * input.HourlyRate
	netSavings := grossSavings - input.CostToAutomate

	return domain.SREToilResult{
		AnnualSavings: netSavings,
		HoursSaved:    annualToilHours,
	}
}
