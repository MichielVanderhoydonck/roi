package service_test

import (
	"testing"
	"time"

	"github.com/MichielVanderhoydonck/roi/internal/service"
)

func TestProductivityService_Calculate(t *testing.T) {
	svc := service.NewProductivityService()

	tests := []struct {
		name                 string
		input                service.ProductivityInput
		expectedTimeSaved    time.Duration
		expectedGrossSavings float64
		expectedNetROI       float64
	}{
		{
			name: "Standard improvement",
			input: service.ProductivityInput{
				TimeBefore:        4 * time.Hour,
				TimeAfter:         5 * time.Minute,
				ExecutionsPerYear: 1000,
				HourlyRate:        75.0,
				MaintenanceCost:   0.0,
			},
			expectedTimeSaved:    (3*time.Hour + 55*time.Minute) * 1000,
			expectedGrossSavings: ((3*time.Hour + 55*time.Minute) * 1000).Hours() * 75.0,
			expectedNetROI:       ((3*time.Hour + 55*time.Minute) * 1000).Hours() * 75.0,
		},
		{
			name: "Negative improvement (no time saved)",
			input: service.ProductivityInput{
				TimeBefore:        5 * time.Minute,
				TimeAfter:         10 * time.Minute,
				ExecutionsPerYear: 100,
				HourlyRate:        50.0,
				MaintenanceCost:   0.0,
			},
			expectedTimeSaved:    0,
			expectedGrossSavings: 0,
			expectedNetROI:       0,
		},
		{
			name: "Maintenance cost exceeding savings",
			input: service.ProductivityInput{
				TimeBefore:        1 * time.Hour,
				TimeAfter:         30 * time.Minute,
				ExecutionsPerYear: 10,
				HourlyRate:        100.0,
				MaintenanceCost:   1000.0,
			},
			expectedTimeSaved:    5 * time.Hour,
			expectedGrossSavings: 500.0,
			expectedNetROI:       -500.0,
		},
		{
			name: "Zero executions",
			input: service.ProductivityInput{
				TimeBefore:        1 * time.Hour,
				TimeAfter:         0,
				ExecutionsPerYear: 0,
				HourlyRate:        100.0,
				MaintenanceCost:   0,
			},
			expectedTimeSaved:    0,
			expectedGrossSavings: 0,
			expectedNetROI:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.Calculate(tt.input)
			if result.TimeSaved != tt.expectedTimeSaved {
				t.Errorf("got time saved %v, expected %v", result.TimeSaved, tt.expectedTimeSaved)
			}
			if result.GrossSavings != tt.expectedGrossSavings {
				t.Errorf("got gross savings %f, expected %f", result.GrossSavings, tt.expectedGrossSavings)
			}
			if result.NetROI != tt.expectedNetROI {
				t.Errorf("got net ROI %f, expected %f", result.NetROI, tt.expectedNetROI)
			}
		})
	}
}
