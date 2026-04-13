package service_test

import (
	"testing"

	"github.com/MichielVanderhoydonck/roi/internal/service"
)

func TestSREToilService_Calculate(t *testing.T) {
	svc := service.NewSREToilService()

	tests := []struct {
		name                  string
		input                 service.SREToilInput
		expectedAnnualSavings float64
		expectedHoursSaved    float64
	}{
		{
			name: "Standard toil eradication",
			input: service.SREToilInput{
				HoursPerWeek:   5,
				HourlyRate:     75,
				CostToAutomate: 1500,
			},
			expectedAnnualSavings: 18000.0,
			expectedHoursSaved:    260.0,
		},
		{
			name: "High automation cost",
			input: service.SREToilInput{
				HoursPerWeek:   2,
				HourlyRate:     50,
				CostToAutomate: 10000,
			},
			expectedAnnualSavings: -4800.0, // (2*52*50) - 10000 = 5200 - 10000
			expectedHoursSaved:    104.0,
		},
		{
			name: "Zero toil",
			input: service.SREToilInput{
				HoursPerWeek:   0,
				HourlyRate:     100,
				CostToAutomate: 0,
			},
			expectedAnnualSavings: 0,
			expectedHoursSaved:    0,
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
