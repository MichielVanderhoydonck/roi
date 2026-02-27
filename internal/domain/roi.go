package domain

import "time"

// ProductivityInput holds the parameters for calculating Developer Productivity ROI.
type ProductivityInput struct {
	TimeBefore        time.Duration
	TimeAfter         time.Duration
	ExecutionsPerYear int
	HourlyRate        float64
	MaintenanceCost   float64
}

// ProductivityResult holds the output of the Developer Productivity ROI calculation.
type ProductivityResult struct {
	GrossSavings float64
	NetROI       float64
	TimeSaved    time.Duration
}

// ProductivityCalculator defines the interface for calculating Productivity ROI.
type ProductivityCalculator interface {
	Calculate(input ProductivityInput) ProductivityResult
}

// ReliabilityInput holds the parameters for calculating Reliability ROI.
type ReliabilityInput struct {
	OldMTTR          time.Duration
	NewMTTR          time.Duration
	IncidentsPerYear int
	DowntimeCost     float64
}

// ReliabilityResult holds the output of the Reliability ROI calculation.
type ReliabilityResult struct {
	DowntimeSavings float64
	TimeSaved       time.Duration
}

// ReliabilityCalculator defines the interface for calculating Reliability ROI.
type ReliabilityCalculator interface {
	Calculate(input ReliabilityInput) ReliabilityResult
}

// FinOpsInput holds the parameters for calculating FinOps ROI.
type FinOpsInput struct {
	OldMonthlyBill float64
	NewMonthlyBill float64
}

// FinOpsResult holds the output of the FinOps ROI calculation.
type FinOpsResult struct {
	AnnualSavings float64
}

// FinOpsCalculator defines the interface for calculating FinOps ROI.
type FinOpsCalculator interface {
	Calculate(input FinOpsInput) FinOpsResult
}

// SREToilInput holds the parameters for calculating SRE Toil Eradication ROI.
type SREToilInput struct {
	HoursPerWeek   float64
	HourlyRate     float64
	CostToAutomate float64
}

// SREToilResult holds the output of the SRE Toil Eradication ROI calculation.
type SREToilResult struct {
	AnnualSavings float64
	HoursSaved    float64
}

// SREToilCalculator defines the interface for calculating SRE Toil Eradication ROI.
type SREToilCalculator interface {
	Calculate(input SREToilInput) SREToilResult
}
