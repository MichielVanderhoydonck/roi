package service_test

import (
	"testing"

	"github.com/MichielVanderhoydonck/roi/internal/service"
)

func TestSREToilService_Calculate(t *testing.T) {
	svc := service.NewSREToilService()
	input := service.SREToilInput{
		HoursPerWeek:   5,
		HourlyRate:     75,
		CostToAutomate: 1500, // e.g. 20 hours * $75
	}
	res := svc.Calculate(input)

	// 5 * 52 = 260 hours
	// 260 * 75 = 19500
	// 19500 - 1500 = 18000
	if res.AnnualSavings != 18000.0 {
		t.Errorf("Expected AnnualSavings %f, got %f", 18000.0, res.AnnualSavings)
	}
	if res.HoursSaved != 260.0 {
		t.Errorf("Expected HoursSaved %f, got %f", 260.0, res.HoursSaved)
	}
}
