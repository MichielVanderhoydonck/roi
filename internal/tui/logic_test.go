package tui

import (
	"testing"
	tea "charm.land/bubbletea/v2"
)

func TestApp_CalculateResultDispatch(t *testing.T) {
	a := NewApp()

	if a.resultText == "" {
		t.Errorf("expected pre-calculated default resultText, got empty string")
	}

	a.activeCalc = NewProductivityCalculator()
	a.activeForm = a.activeCalc.CreateForm()

	res, _ := a.activeCalc.CalculateResult(a.activeForm)
	if res == "" {
		t.Error("expected a valid result string from CalculateResult")
	}
}

func TestApp_ResetState(t *testing.T) {
	a := NewApp()
	a.focus = focusForm
	a.activeForm.State = 1 // huh.StateCompleted
	
	handled, _ := a.handleFormKey(tea.KeyPressMsg{Text: "esc"})
	
	if !handled {
		t.Error("expected handleFormKey to handle esc")
	}
	
	if a.activeForm.State != 0 {
		t.Errorf("expected form state to be Normal, got %v", a.activeForm.State)
	}
	if a.focus != focusMenu {
		t.Errorf("expected focus to be menu, got %v", a.focus)
	}
}

func TestApp_SentimentColors(t *testing.T) {
	// Test positive sentiment with Productivity
	pCalc := NewProductivityCalculator()
	pCalc.timeBefore = "10h"
	pCalc.timeAfter = "1h"
	pCalc.executions = "100"
	pCalc.hourlyRate = "100"
	pCalc.maintenance = "500" // gross savings = 9h * 100 * 100 = 90,000 > 500
	form := pCalc.CreateForm()
	_, sentiment := pCalc.CalculateResult(form)
	if sentiment != SentimentGood {
		t.Errorf("expected SentimentGood for positive net ROI, got %v", sentiment)
	}

	// Test negative sentiment with Productivity
	pCalc.maintenance = "1000000" // massive maintenance cost causes negative net ROI
	form = pCalc.CreateForm()
	_, sentiment = pCalc.CalculateResult(form)
	if sentiment != SentimentBad {
		t.Errorf("expected SentimentBad for negative net ROI, got %v", sentiment)
	}

	// Test negative sentiment with Cost of Delay
	dCalc := NewCostOfDelayCalculator()
	dCalc.monthlyRevenue = "300000"
	dCalc.daysDelayed = "15"
	form = dCalc.CreateForm()
	_, sentiment = dCalc.CalculateResult(form)
	if sentiment != SentimentBad {
		t.Errorf("expected SentimentBad for cost of delay revenue lost, got %v", sentiment)
	}
}
