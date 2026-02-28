package service

import "github.com/MichielVanderhoydonck/roi/internal/domain"

type CostOfDelayService struct{}

func NewCostOfDelayService() *CostOfDelayService {
	return &CostOfDelayService{}
}

func (s *CostOfDelayService) Calculate(input domain.CostOfDelayInput) domain.CostOfDelayResult {
	costOfDelay := (input.EstimatedMonthlyRevenue / 30.0) * input.DaysDelayed

	return domain.CostOfDelayResult{
		CostOfDelay: costOfDelay,
	}
}
