package service

import "github.com/MichielVanderhoydonck/roi/internal/domain"

type OnboardingService struct{}

func NewOnboardingService() *OnboardingService {
	return &OnboardingService{}
}

func (s *OnboardingService) Calculate(input domain.OnboardingInput) domain.OnboardingResult {
	daysSaved := input.OldDays - input.NewDays
	savings := daysSaved * float64(input.NewHires) * input.DailyRate

	return domain.OnboardingResult{
		AnnualSavings:    savings,
		DaysSavedPerHire: daysSaved,
	}
}
