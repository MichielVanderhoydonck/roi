package tui

import "github.com/charmbracelet/huh"

type focusState int

const (
	focusMenu focusState = iota
	focusForm
)

type Calculator interface {
	CreateForm() *huh.Form
	CalculateResult(form *huh.Form) string
	GetFormula(form *huh.Form) string
	GetContext(fieldKey string) string
}

type item struct {
	title, desc string
	calc        Calculator
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
