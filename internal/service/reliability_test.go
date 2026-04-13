package service_test

import (
	"testing"
	"time"

	"github.com/MichielVanderhoydonck/roi/internal/service"
)

func TestReliabilityService_Calculate(t *testing.T) {
	svc := service.NewReliabilityService()

	tests := []struct {
		name                    string
		input                   service.ReliabilityInput
		expectedTimeSaved       time.Duration
		expectedDowntimeSavings float64
	}{
		{
			name: "Standard improvement",
			input: service.ReliabilityInput{
				OldMTTR:          2 * time.Hour,
				NewMTTR:          30 * time.Minute,
				IncidentsPerYear: 10,
				DowntimeCost:     50000.0,
			},
			expectedTimeSaved:       15 * time.Hour,
			expectedDowntimeSavings: 15.0 * 50000.0,
		},
		{
			name: "No improvement (identical MTTR)",
			input: service.ReliabilityInput{
				OldMTTR:          1 * time.Hour,
				NewMTTR:          1 * time.Hour,
				IncidentsPerYear: 5,
				DowntimeCost:     10000.0,
			},
			expectedTimeSaved:       0,
			expectedDowntimeSavings: 0,
		},
		{
			name: "Negative improvement (worse MTTR)",
			input: service.ReliabilityInput{
				OldMTTR:          30 * time.Minute,
				NewMTTR:          2 * time.Hour,
				IncidentsPerYear: 10,
				DowntimeCost:     50000.0,
			},
			expectedTimeSaved:       0,
			expectedDowntimeSavings: 0,
		},
		{
			name: "Zero incidents",
			input: service.ReliabilityInput{
				OldMTTR:          2 * time.Hour,
				NewMTTR:          1 * time.Hour,
				IncidentsPerYear: 0,
				DowntimeCost:     50000.0,
			},
			expectedTimeSaved:       0,
			expectedDowntimeSavings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.Calculate(tt.input)
			if result.TimeSaved != tt.expectedTimeSaved {
				t.Errorf("got time saved %v, expected %v", result.TimeSaved, tt.expectedTimeSaved)
			}
			if result.DowntimeSavings != tt.expectedDowntimeSavings {
				t.Errorf("got downtime savings %f, expected %f", result.DowntimeSavings, tt.expectedDowntimeSavings)
			}
		})
	}
}
