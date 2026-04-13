package service_test

import (
	"testing"

	"github.com/MichielVanderhoydonck/roi/internal/service"
)

func TestOnboardingService_Calculate(t *testing.T) {
	svc := service.NewOnboardingService()

	tests := []struct {
		name                   string
		input                  service.OnboardingInput
		expectedAnnualSavings  float64
		expectedDaysSaved      float64
	}{
		{
			name: "Standard onboarding improvement",
			input: service.OnboardingInput{
				OldDays:   10,
				NewDays:   5,
				NewHires:  4,
				DailyRate: 500,
			},
			expectedAnnualSavings: 10000.0, // (10-5) * 4 * 500 = 5 * 4 * 500 = 10000
			expectedDaysSaved:      5.0,
		},
		{
			name: "No improvement",
			input: service.OnboardingInput{
				OldDays:   7,
				NewDays:   7,
				NewHires:  10,
				DailyRate: 400,
			},
			expectedAnnualSavings: 0,
			expectedDaysSaved:      0,
		},
		{
			name: "Zero hires",
			input: service.OnboardingInput{
				OldDays:   20,
				NewDays:   10,
				NewHires:  0,
				DailyRate: 600,
			},
			expectedAnnualSavings: 0,
			expectedDaysSaved:      10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.Calculate(tt.input)
			if result.AnnualSavings != tt.expectedAnnualSavings {
				t.Errorf("got annual savings %f, expected %f", result.AnnualSavings, tt.expectedAnnualSavings)
			}
			if result.DaysSavedPerHire != tt.expectedDaysSaved {
				t.Errorf("got days saved %f, expected %f", result.DaysSavedPerHire, tt.expectedDaysSaved)
			}
		})
	}
}
