package service_test

import (
	"testing"

	"github.com/MichielVanderhoydonck/roi/internal/domain"
	"github.com/MichielVanderhoydonck/roi/internal/service"
)

func TestFinOpsService_Calculate(t *testing.T) {
	svc := service.NewFinOpsService()

	input := domain.FinOpsInput{
		OldMonthlyBill: 20000.0,
		NewMonthlyBill: 15000.0,
	}

	result := svc.Calculate(input)

	expectedMonthlySavings := 5000.0
	expectedAnnualSavings := expectedMonthlySavings * 12

	if result.AnnualSavings != expectedAnnualSavings {
		t.Errorf("expected annual savings %f, got %f", expectedAnnualSavings, result.AnnualSavings)
	}

	// Negative savings test
	inputNegative := domain.FinOpsInput{
		OldMonthlyBill: 15000.0,
		NewMonthlyBill: 20000.0,
	}

	resultNegative := svc.Calculate(inputNegative)
	if resultNegative.AnnualSavings != 0 {
		t.Errorf("expected 0 annual savings for bill increase, got %f", resultNegative.AnnualSavings)
	}
}
