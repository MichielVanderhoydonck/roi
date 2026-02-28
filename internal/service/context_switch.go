package service

import "github.com/MichielVanderhoydonck/roi/internal/domain"

const contextSwitchPenaltyHours = 0.4

type ContextSwitchService struct{}

func NewContextSwitchService() *ContextSwitchService {
	return &ContextSwitchService{}
}

func (s *ContextSwitchService) Calculate(input domain.ContextSwitchInput) domain.ContextSwitchResult {
	hoursSaved := float64(input.ReducedIncidentsPerYear) * contextSwitchPenaltyHours
	savings := hoursSaved * input.HourlyRate

	return domain.ContextSwitchResult{
		AnnualSavings: savings,
		HoursSaved:    hoursSaved,
	}
}
