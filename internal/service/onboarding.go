package service

// OnboardingInput holds the parameters for calculating Onboarding ROI.
type OnboardingInput struct {
	OldDays   float64
	NewDays   float64
	NewHires  int
	DailyRate float64
}

// OnboardingResult holds the output of the Onboarding ROI calculation.
type OnboardingResult struct {
	AnnualSavings    float64
	DaysSavedPerHire float64
}

// OnboardingCalculator defines the interface for calculating Onboarding ROI.
type OnboardingCalculator interface {
	Calculate(input OnboardingInput) OnboardingResult
}

type OnboardingService struct{}

func NewOnboardingService() *OnboardingService {
	return &OnboardingService{}
}

func (s *OnboardingService) Calculate(input OnboardingInput) OnboardingResult {
	daysSaved := input.OldDays - input.NewDays
	savings := daysSaved * float64(input.NewHires) * input.DailyRate

	return OnboardingResult{
		AnnualSavings:    savings,
		DaysSavedPerHire: daysSaved,
	}
}
