package service

import (
	"time"

	"github.com/michiel/roi/internal/domain"
)

// ReliabilityService implements the ReliabilityCalculator interface.
type ReliabilityService struct{}

// NewReliabilityService creates a new ReliabilityService.
func NewReliabilityService() *ReliabilityService {
	return &ReliabilityService{}
}

// Calculate computes the Reliability ROI.
// Downtime Savings = (Old MTTR - New MTTR) * Incidents * Downtime Cost
func (s *ReliabilityService) Calculate(input domain.ReliabilityInput) domain.ReliabilityResult {
	mttrReduction := input.OldMTTR - input.NewMTTR
	if mttrReduction < 0 {
		mttrReduction = 0
	}

	totalTimeSaved := mttrReduction * time.Duration(input.IncidentsPerYear)
	totalHoursSaved := totalTimeSaved.Hours()

	downtimeSavings := totalHoursSaved * input.DowntimeCost

	return domain.ReliabilityResult{
		DowntimeSavings: downtimeSavings,
		TimeSaved:       totalTimeSaved,
	}
}
