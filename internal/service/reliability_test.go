package service_test

import (
	"testing"
	"time"

	"github.com/michiel/roi/internal/domain"
	"github.com/michiel/roi/internal/service"
)

func TestReliabilityService_Calculate(t *testing.T) {
	svc := service.NewReliabilityService()

	input := domain.ReliabilityInput{
		OldMTTR:          2 * time.Hour,
		NewMTTR:          30 * time.Minute,
		IncidentsPerYear: 10,
		DowntimeCost:     50000.0,
	}

	result := svc.Calculate(input)

	expectedMTTRReduction := 1*time.Hour + 30*time.Minute
	expectedTotalTimeSaved := expectedMTTRReduction * 10
	expectedDowntimeSavings := expectedTotalTimeSaved.Hours() * 50000.0

	if result.TimeSaved != expectedTotalTimeSaved {
		t.Errorf("expected total time saved %v, got %v", expectedTotalTimeSaved, result.TimeSaved)
	}

	if result.DowntimeSavings != expectedDowntimeSavings {
		t.Errorf("expected downtime savings %f, got %f", expectedDowntimeSavings, result.DowntimeSavings)
	}

	// Test negative savings (worse MTTR)
	inputNegative := domain.ReliabilityInput{
		OldMTTR:          30 * time.Minute,
		NewMTTR:          2 * time.Hour,
		IncidentsPerYear: 10,
		DowntimeCost:     50000.0,
	}

	resultNegative := svc.Calculate(inputNegative)
	if resultNegative.TimeSaved != 0 {
		t.Errorf("expected 0 time saved for negative improvement, got %v", resultNegative.TimeSaved)
	}
	if resultNegative.DowntimeSavings != 0 {
		t.Errorf("expected 0 downtime savings, got %f", resultNegative.DowntimeSavings)
	}
}
