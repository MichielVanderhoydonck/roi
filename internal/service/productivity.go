package service

import (
	"time"

	"github.com/MichielVanderhoydonck/roi/internal/domain"
)

// ProductivityService implements the ProductivityCalculator interface.
type ProductivityService struct{}

// NewProductivityService creates a new ProductivityService.
func NewProductivityService() *ProductivityService {
	return &ProductivityService{}
}

// Calculate computes the Developer Productivity ROI.
// Annual ROI = [(Time BEFORE - Time AFTER) * Executions * HourlyRate] - MaintenanceCost
func (s *ProductivityService) Calculate(input domain.ProductivityInput) domain.ProductivityResult {
	timeSavedPerExecution := input.TimeBefore - input.TimeAfter
	if timeSavedPerExecution < 0 {
		timeSavedPerExecution = 0
	}

	totalTimeSaved := timeSavedPerExecution * time.Duration(input.ExecutionsPerYear)
	totalHoursSaved := totalTimeSaved.Hours()

	grossSavings := totalHoursSaved * input.HourlyRate
	netROI := grossSavings - input.MaintenanceCost

	return domain.ProductivityResult{
		GrossSavings: grossSavings,
		NetROI:       netROI,
		TimeSaved:    totalTimeSaved,
	}
}
