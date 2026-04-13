package tui

type focusState int

const (
	focusMenu focusState = iota
	focusForm
)

type calcType string

const (
	calcProductivity  calcType = "productivity"
	calcReliability   calcType = "reliability"
	calcFinOps        calcType = "finops"
	calcSRE           calcType = "sre"
	calcOnboarding    calcType = "onboarding"
	calcContextSwitch calcType = "context_switch"
	calcCostOfDelay   calcType = "cost_of_delay"
)

type item struct {
	title, desc string
	calc        calcType
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
