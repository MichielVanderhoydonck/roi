package service

import (
)

// SREToilInput holds the parameters for calculating SRE Toil Eradication ROI.
type SREToilInput struct {
	HoursPerWeek   float64
	HourlyRate     float64
	CostToAutomate float64
}

// SREToilResult holds the output of the SRE Toil Eradication ROI calculation.
type SREToilResult struct {
	AnnualSavings float64
	HoursSaved    float64
}

// SREToilCalculator defines the interface for calculating SRE Toil Eradication ROI.
type SREToilCalculator interface {
	Calculate(input SREToilInput) SREToilResult
}

// SREToilService implements the SREToilCalculator interface.
type SREToilService struct{}

// NewSREToilService creates a new SREToilService.
func NewSREToilService() *SREToilService {
	return &SREToilService{}
}

// Calculate computes the SRE Toil Eradication ROI.
// Toil Savings ($) = [(Hours of Toil per Week × 52) × (Hourly Dev Rate)] - Cost to Automate
func (s *SREToilService) Calculate(input SREToilInput) SREToilResult {
	annualToilHours := input.HoursPerWeek * 52
	grossSavings := annualToilHours * input.HourlyRate
	netSavings := grossSavings - input.CostToAutomate

	return SREToilResult{
		AnnualSavings: netSavings,
		HoursSaved:    annualToilHours,
	}
}
