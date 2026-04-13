package service

// CostOfDelayInput holds the parameters for calculating Cost of Delay ROI.
type CostOfDelayInput struct {
	EstimatedMonthlyRevenue float64
	DaysDelayed             float64
}

// CostOfDelayResult holds the output of the Cost of Delay ROI calculation.
type CostOfDelayResult struct {
	CostOfDelay float64
}

// CostOfDelayCalculator defines the interface for calculating Cost of Delay ROI.
type CostOfDelayCalculator interface {
	Calculate(input CostOfDelayInput) CostOfDelayResult
}

type CostOfDelayService struct{}

func NewCostOfDelayService() *CostOfDelayService {
	return &CostOfDelayService{}
}

func (s *CostOfDelayService) Calculate(input CostOfDelayInput) CostOfDelayResult {
	costOfDelay := (input.EstimatedMonthlyRevenue / 30.0) * input.DaysDelayed

	return CostOfDelayResult{
		CostOfDelay: costOfDelay,
	}
}
