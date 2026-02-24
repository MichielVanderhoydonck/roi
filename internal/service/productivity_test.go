package service_test

import (
	"testing"
	"time"

	"github.com/michiel/roi/internal/domain"
	"github.com/michiel/roi/internal/service"
)

func TestProductivityService_Calculate(t *testing.T) {
	svc := service.NewProductivityService()

	input := domain.ProductivityInput{
		TimeBefore:        4 * time.Hour,
		TimeAfter:         5 * time.Minute,
		ExecutionsPerYear: 1000,
		HourlyRate:        75.0,
		MaintenanceCost:   0.0,
	}

	result := svc.Calculate(input)

	expectedTimeSavedPerExec := 3*time.Hour + 55*time.Minute
	expectedTotalTimeSaved := expectedTimeSavedPerExec * 1000
	expectedGrossSavings := expectedTotalTimeSaved.Hours() * 75.0

	if result.TimeSaved != expectedTotalTimeSaved {
		t.Errorf("expected time saved %v, got %v", expectedTotalTimeSaved, result.TimeSaved)
	}

	if result.GrossSavings != expectedGrossSavings {
		t.Errorf("expected gross savings %f, got %f", expectedGrossSavings, result.GrossSavings)
	}

	if result.NetROI != expectedGrossSavings {
		t.Errorf("expected net ROI %f, got %f", expectedGrossSavings, result.NetROI)
	}

	// Test case with negative savings
	inputNegative := domain.ProductivityInput{
		TimeBefore:        5 * time.Minute,
		TimeAfter:         4 * time.Hour,
		ExecutionsPerYear: 1000,
		HourlyRate:        75.0,
		MaintenanceCost:   100.0,
	}

	resultNegative := svc.Calculate(inputNegative)
	if resultNegative.TimeSaved != 0 {
		t.Errorf("expected no time saved for negative improvement, got %v", resultNegative.TimeSaved)
	}
	if resultNegative.GrossSavings != 0 {
		t.Errorf("expected no gross savings for negative improvement, got %f", resultNegative.GrossSavings)
	}
	if resultNegative.NetROI != -100.0 {
		t.Errorf("expected negative net ROI, got %f", resultNegative.NetROI)
	}
}
