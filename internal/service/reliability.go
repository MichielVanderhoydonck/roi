package service

import (
	"time"
)

// ReliabilityInput holds the parameters for calculating Reliability ROI.
type ReliabilityInput struct {
	OldMTTR          time.Duration
	NewMTTR          time.Duration
	IncidentsPerYear int
	DowntimeCost     float64
}

// ReliabilityResult holds the output of the Reliability ROI calculation.
type ReliabilityResult struct {
	DowntimeSavings float64
	TimeSaved       time.Duration
}

// ReliabilityCalculator defines the interface for calculating Reliability ROI.
type ReliabilityCalculator interface {
	Calculate(input ReliabilityInput) ReliabilityResult
}

// ReliabilityService implements the ReliabilityCalculator interface.
type ReliabilityService struct{}

// NewReliabilityService creates a new ReliabilityService.
func NewReliabilityService() *ReliabilityService {
	return &ReliabilityService{}
}

// Calculate computes the Reliability ROI.
// Downtime Savings = (Old MTTR - New MTTR) * Incidents * Downtime Cost
func (s *ReliabilityService) Calculate(input ReliabilityInput) ReliabilityResult {
	mttrReduction := input.OldMTTR - input.NewMTTR
	if mttrReduction < 0 {
		mttrReduction = 0
	}

	totalTimeSaved := mttrReduction * time.Duration(input.IncidentsPerYear)
	totalHoursSaved := totalTimeSaved.Hours()

	downtimeSavings := totalHoursSaved * input.DowntimeCost

	return ReliabilityResult{
		DowntimeSavings: downtimeSavings,
		TimeSaved:       totalTimeSaved,
	}
}
