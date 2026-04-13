package service_test

import (
	"testing"

	"github.com/MichielVanderhoydonck/roi/internal/service"
)

func TestFinOpsService_Calculate(t *testing.T) {
	svc := service.NewFinOpsService()

	tests := []struct {
		name     string
		input    service.FinOpsInput
		expected float64
	}{
		{
			name: "Standard savings",
			input: service.FinOpsInput{
				OldMonthlyBill: 20000.0,
				NewMonthlyBill: 15000.0,
			},
			expected: 60000.0,
		},
		{
			name: "Bill increase (negative savings should be 0)",
			input: service.FinOpsInput{
				OldMonthlyBill: 15000.0,
				NewMonthlyBill: 20000.0,
			},
			expected: 0,
		},
		{
			name: "Identical bills",
			input: service.FinOpsInput{
				OldMonthlyBill: 10000.0,
				NewMonthlyBill: 10000.0,
			},
			expected: 0,
		},
		{
			name: "Zero bills",
			input: service.FinOpsInput{
				OldMonthlyBill: 0,
				NewMonthlyBill: 0,
			},
			expected: 0,
		},
		{
			name: "High impact project",
			input: service.FinOpsInput{
				OldMonthlyBill: 1000000.0,
				NewMonthlyBill: 100000.0,
			},
			expected: 10800000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.Calculate(tt.input)
			if result.AnnualSavings != tt.expected {
				t.Errorf("got annual savings %f, expected %f", result.AnnualSavings, tt.expected)
			}
		})
	}
}
