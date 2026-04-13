package service_test

import (
	"testing"

	"github.com/MichielVanderhoydonck/roi/internal/service"
)

func TestContextSwitchService_Calculate(t *testing.T) {
	svc := service.NewContextSwitchService()

	tests := []struct {
		name                  string
		input                 service.ContextSwitchInput
		expectedAnnualSavings float64
		expectedHoursSaved    float64
	}{
		{
			name: "Reduce 100 incidents per year",
			input: service.ContextSwitchInput{
				ReducedIncidentsPerYear: 100,
				HourlyRate:              100,
			},
			expectedAnnualSavings: 4000.0, // 100 * 0.4 * 100 = 40 * 100 = 4000
			expectedHoursSaved:    40.0,
		},
		{
			name: "Zero incidents reduced",
			input: service.ContextSwitchInput{
				ReducedIncidentsPerYear: 0,
				HourlyRate:              120,
			},
			expectedAnnualSavings: 0,
			expectedHoursSaved:    0,
		},
		{
			name: "High hourly rate",
			input: service.ContextSwitchInput{
				ReducedIncidentsPerYear: 50,
				HourlyRate:              200,
			},
			expectedAnnualSavings: 4000.0, // 50 * 0.4 * 200 = 20 * 200 = 4000
			expectedHoursSaved:    20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.Calculate(tt.input)
			if result.AnnualSavings != tt.expectedAnnualSavings {
				t.Errorf("got annual savings %f, expected %f", result.AnnualSavings, tt.expectedAnnualSavings)
			}
			if result.HoursSaved != tt.expectedHoursSaved {
				t.Errorf("got hours saved %f, expected %f", result.HoursSaved, tt.expectedHoursSaved)
			}
		})
	}
}
