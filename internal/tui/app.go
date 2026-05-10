package tui

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/huh"
)

type App struct {
	menuList   list.Model
	activeCalc Calculator
	activeForm *huh.Form
	resultText string
	focus      focusState
	width      int
	height     int
}

func NewApp() *App {
	items := []list.Item{
		item{title: "Developer Productivity", desc: "Time Saved", calc: NewProductivityCalculator()},
		item{title: "Reliability", desc: "Cost of Downtime Avoided", calc: NewReliabilityCalculator()},
		item{title: "FinOps", desc: "Infrastructure Optimization", calc: NewFinOpsCalculator()},
		item{title: "SRE Toil Eradication", desc: "Automating Manual Work", calc: NewSRECalculator()},
		item{title: "Onboarding ROI", desc: "Time to First Value", calc: NewOnboardingCalculator()},
		item{title: "Cost of Context Switching", desc: "The Hidden Tax", calc: NewContextSwitchCalculator()},
		item{title: "Cost of Delay", desc: "The Product Velocity Metric", calc: NewCostOfDelayCalculator()},
	}

	m := list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.Title = "ROI Calculator"
	m.SetShowHelp(false)
	m.SetShowStatusBar(false)
	m.SetFilteringEnabled(false)

	initialCalc := items[0].(item).calc
	initialForm := initialCalc.CreateForm()

	return &App{
		focus:      focusMenu,
		menuList:   m,
		activeCalc: initialCalc,
		activeForm: initialForm,
	}
}

func (a *App) Run() error {
	p := tea.NewProgram(a)
	_, err := p.Run()
	return err
}

func (a *App) Init() tea.Cmd {
	return wrapCmd(a.activeForm.Init())
}

// Update handles application state changes
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}
		if handled, cmd := a.handleKeyMsg(msg); handled {
			return a, cmd
		}
	case tea.WindowSizeMsg:
		a.handleWindowSizeMsg(msg)
	}

	return a.updateFocus(msg)
}

func (a *App) handleKeyMsg(msg tea.KeyMsg) (bool, tea.Cmd) {
	switch a.focus {
	case focusMenu:
		return a.handleMenuKey(msg)
	case focusForm:
		return a.handleFormKey(msg)
	}
	return false, nil
}

func (a *App) handleMenuKey(msg tea.KeyMsg) (bool, tea.Cmd) {
	if msg.String() == "enter" || msg.String() == "right" || msg.String() == "tab" {
		a.focus = focusForm
		return true, nil
	}
	return false, nil
}

func (a *App) handleFormKey(msg tea.KeyMsg) (bool, tea.Cmd) {
	if msg.String() == "esc" {
		if a.activeForm.State == huh.StateCompleted {
			a.activeForm.State = huh.StateNormal
		}
		a.focus = focusMenu
		return true, nil
	}

	// When form is completed, enter returns to menu and resets state for editing
	if msg.String() == "enter" {
		if a.activeForm != nil && a.activeForm.State == huh.StateCompleted {
			a.focus = focusMenu
			a.activeForm.State = huh.StateNormal
			return true, nil
		}
	}

	if msg.String() == "ctrl+r" {
		a.activeForm = a.activeCalc.CreateForm()
		return true, wrapCmd(a.activeForm.Init())
	}

	return false, nil
}

func (a *App) handleWindowSizeMsg(msg tea.WindowSizeMsg) {
	a.width = msg.Width
	a.height = msg.Height

	mainHeight := max(0, a.height-2)
	h, v := lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.RoundedBorder()).GetFrameSize()

	leftWidth := max(30, min(45, a.width*30/100))
	leftContentWidth := max(0, leftWidth-h)

	a.menuList.SetSize(leftContentWidth, mainHeight-v)
}

func (a *App) updateFocus(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch a.focus {
	case focusMenu:
		oldIdx := a.menuList.Index()
		var cmd tea.Cmd
		a.menuList, cmd = a.menuList.Update(msg)
		cmds = append(cmds, cmd)

		if a.menuList.Index() != oldIdx {
			selectedItem := a.menuList.SelectedItem().(item)
			a.activeCalc = selectedItem.calc
			a.activeForm = a.activeCalc.CreateForm()
			cmds = append(cmds, wrapCmd(a.activeForm.Init()))
			a.resultText = ""
		}
	case focusForm:
		newModel, newCmd := a.activeForm.Update(mapV2MsgToV1(msg))
		if f, ok := newModel.(*huh.Form); ok {
			a.activeForm = f
			if f.State == huh.StateCompleted {
				a.resultText = a.activeCalc.CalculateResult(a.activeForm)
			}
		}
		cmds = append(cmds, wrapCmd(newCmd))
	}

	return a, tea.Batch(cmds...)
}

func (a *App) getActiveForm() *huh.Form {
	return a.activeForm
}

// View and Rendering
func (a *App) View() tea.View {
	if a.width == 0 {
		return tea.NewView("Starting...")
	}

	headerStr := a.renderHeader()
	footerStr := a.renderFooter()

	headerHeight := lipgloss.Height(headerStr)
	footerHeight := lipgloss.Height(footerStr)
	mainHeight := max(0, a.height-headerHeight-footerHeight)

	leftWidth, rightWidth := a.calculateLayoutWidths()

	leftPanel := a.renderLeftPanel(leftWidth, mainHeight)
	rightPanel := a.renderRightPanel(rightWidth, mainHeight)

	mainLayout := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left, headerStr, mainLayout, footerStr))
	v.AltScreen = true
	return v
}
