package service_test

import (
	"testing"

	"github.com/MichielVanderhoydonck/roi/internal/service"
)

func TestCostOfDelayService_Calculate(t *testing.T) {
	svc := service.NewCostOfDelayService()

	tests := []struct {
		name                string
		input               service.CostOfDelayInput
		expectedCostOfDelay float64
	}{
		{
			name: "Standard delay",
			input: service.CostOfDelayInput{
				EstimatedMonthlyRevenue: 30000,
				DaysDelayed:             10,
			},
			expectedCostOfDelay: 10000.0, // (30000/30) * 10 = 1000 * 10 = 10000
		},
		{
			name: "No delay",
			input: service.CostOfDelayInput{
				EstimatedMonthlyRevenue: 50000,
				DaysDelayed:             0,
			},
			expectedCostOfDelay: 0,
		},
		{
			name: "Zero revenue",
			input: service.CostOfDelayInput{
				EstimatedMonthlyRevenue: 0,
				DaysDelayed:             30,
			},
			expectedCostOfDelay: 0,
		},
		{
			name: "Large revenue project",
			input: service.CostOfDelayInput{
				EstimatedMonthlyRevenue: 900000,
				DaysDelayed:             5,
			},
			expectedCostOfDelay: 150000.0, // (900000/30) * 5 = 30000 * 5 = 150000
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.Calculate(tt.input)
			if result.CostOfDelay != tt.expectedCostOfDelay {
				t.Errorf("got cost of delay %f, expected %f", result.CostOfDelay, tt.expectedCostOfDelay)
			}
		})
	}
}
