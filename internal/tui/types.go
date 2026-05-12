package tui

import "charm.land/huh/v2"

type focusState int

const (
	focusMenu focusState = iota
	focusForm
)

type Sentiment int

const (
	SentimentNone Sentiment = iota
	SentimentGood
	SentimentBad
)

type Calculator interface {
	CreateForm() *huh.Form
	CalculateResult(form *huh.Form) (string, Sentiment)
	GetFormula(form *huh.Form) string
	GetContext(fieldKey string) string
	Reset()
}

type item struct {
	title, desc string
	calc        Calculator
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
