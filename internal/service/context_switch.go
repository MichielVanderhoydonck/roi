package service

// ContextSwitchInput holds the parameters for calculating Context Switch ROI.
type ContextSwitchInput struct {
	ReducedIncidentsPerYear int
	HourlyRate              float64
}

// ContextSwitchResult holds the output of the Context Switch ROI calculation.
type ContextSwitchResult struct {
	AnnualSavings float64
	HoursSaved    float64
}

// ContextSwitchCalculator defines the interface for calculating Context Switch ROI.
type ContextSwitchCalculator interface {
	Calculate(input ContextSwitchInput) ContextSwitchResult
}

const contextSwitchPenaltyHours = 0.4

type ContextSwitchService struct{}

func NewContextSwitchService() *ContextSwitchService {
	return &ContextSwitchService{}
}

func (s *ContextSwitchService) Calculate(input ContextSwitchInput) ContextSwitchResult {
	hoursSaved := float64(input.ReducedIncidentsPerYear) * contextSwitchPenaltyHours
	savings := hoursSaved * input.HourlyRate

	return ContextSwitchResult{
		AnnualSavings: savings,
		HoursSaved:    hoursSaved,
	}
}
