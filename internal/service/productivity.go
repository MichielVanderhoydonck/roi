package service

import (
	"time"
)

// ProductivityInput holds the parameters for calculating Developer Productivity ROI.
type ProductivityInput struct {
	TimeBefore        time.Duration
	TimeAfter         time.Duration
	ExecutionsPerYear int
	HourlyRate        float64
	MaintenanceCost   float64
}

// ProductivityResult holds the output of the Developer Productivity ROI calculation.
type ProductivityResult struct {
	GrossSavings float64
	NetROI       float64
	TimeSaved    time.Duration
}

// ProductivityCalculator defines the interface for calculating Productivity ROI.
type ProductivityCalculator interface {
	Calculate(input ProductivityInput) ProductivityResult
}

// ProductivityService implements the ProductivityCalculator interface.
type ProductivityService struct{}

// NewProductivityService creates a new ProductivityService.
func NewProductivityService() *ProductivityService {
	return &ProductivityService{}
}

// Calculate computes the Developer Productivity ROI.
// Annual ROI = [(Time BEFORE - Time AFTER) * Executions * HourlyRate] - MaintenanceCost
func (s *ProductivityService) Calculate(input ProductivityInput) ProductivityResult {
	timeSavedPerExecution := input.TimeBefore - input.TimeAfter
	if timeSavedPerExecution < 0 {
		timeSavedPerExecution = 0
	}

	totalTimeSaved := timeSavedPerExecution * time.Duration(input.ExecutionsPerYear)
	totalHoursSaved := totalTimeSaved.Hours()

	grossSavings := totalHoursSaved * input.HourlyRate
	netROI := grossSavings - input.MaintenanceCost

	return ProductivityResult{
		GrossSavings: grossSavings,
		NetROI:       netROI,
		TimeSaved:    totalTimeSaved,
	}
}
