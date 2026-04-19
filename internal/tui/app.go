package tui

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	"github.com/charmbracelet/huh"
)

type App struct {
	menuList             list.Model
	onboardingService    *service.OnboardingService
	finForm              *huh.Form
	sreService           *service.SREToilService
	prodService          *service.ProductivityService
	contextSwitchService *service.ContextSwitchService
	costOfDelayService   *service.CostOfDelayService
	costOfDelayForm      *huh.Form
	relService           *service.ReliabilityService
	finService           *service.FinOpsService
	prodForm             *huh.Form
	relForm              *huh.Form
	sreForm              *huh.Form
	onboardingForm       *huh.Form
	contextSwitchForm    *huh.Form
	resultText           string
	focus                focusState
	width                int
	height               int
}

func NewApp() *App {
	items := []list.Item{
		item{title: "Developer Productivity", desc: "Time Saved", calc: calcProductivity},
		item{title: "Reliability", desc: "Cost of Downtime Avoided", calc: calcReliability},
		item{title: "FinOps", desc: "Infrastructure Optimization", calc: calcFinOps},
		item{title: "SRE Toil Eradication", desc: "Automating Manual Work", calc: calcSRE},
		item{title: "Onboarding ROI", desc: "Time to First Value", calc: calcOnboarding},
		item{title: "Cost of Context Switching", desc: "The Hidden Tax", calc: calcContextSwitch},
		item{title: "Cost of Delay", desc: "The Product Velocity Metric", calc: calcCostOfDelay},
	}

	m := list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.Title = "ROI Calculator"
	m.SetShowHelp(false)
	m.SetShowStatusBar(false)
	m.SetFilteringEnabled(false)

	return &App{
		prodService:          service.NewProductivityService(),
		relService:           service.NewReliabilityService(),
		finService:           service.NewFinOpsService(),
		sreService:           service.NewSREToilService(),
		onboardingService:    service.NewOnboardingService(),
		contextSwitchService: service.NewContextSwitchService(),
		costOfDelayService:   service.NewCostOfDelayService(),
		focus:                focusMenu,
		menuList:             m,
		prodForm:             createProductivityForm(),
		relForm:              createReliabilityForm(),
		finForm:              createFinOpsForm(),
		sreForm:              createSREForm(),
		onboardingForm:       createOnboardingForm(),
		contextSwitchForm:    createContextSwitchForm(),
		costOfDelayForm:      createCostOfDelayForm(),
	}
}

func (a *App) Run() error {
	p := tea.NewProgram(a)
	_, err := p.Run()
	return err
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(
		wrapCmd(a.prodForm.Init()),
		wrapCmd(a.relForm.Init()),
		wrapCmd(a.finForm.Init()),
		wrapCmd(a.sreForm.Init()),
		wrapCmd(a.onboardingForm.Init()),
		wrapCmd(a.contextSwitchForm.Init()),
		wrapCmd(a.costOfDelayForm.Init()),
	)
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
		a.resetFormState()
		a.focus = focusMenu
		return true, nil
	}

	// When form is completed, enter returns to menu and resets state for editing
	if msg.String() == "enter" {
		activeForm := a.getActiveForm()
		if activeForm != nil && activeForm.State == huh.StateCompleted {
			a.focus = focusMenu
			activeForm.State = huh.StateNormal
			return true, nil
		}
	}

	if msg.String() == "ctrl+r" {
		return true, a.clearActiveForm()
	}

	return false, nil
}

func (a *App) resetFormState() {
	if a.prodForm.State == huh.StateCompleted {
		a.prodForm.State = huh.StateNormal
	}
	if a.relForm.State == huh.StateCompleted {
		a.relForm.State = huh.StateNormal
	}
	if a.finForm.State == huh.StateCompleted {
		a.finForm.State = huh.StateNormal
	}
	if a.sreForm.State == huh.StateCompleted {
		a.sreForm.State = huh.StateNormal
	}
	if a.onboardingForm.State == huh.StateCompleted {
		a.onboardingForm.State = huh.StateNormal
	}
	if a.contextSwitchForm.State == huh.StateCompleted {
		a.contextSwitchForm.State = huh.StateNormal
	}
	if a.costOfDelayForm.State == huh.StateCompleted {
		a.costOfDelayForm.State = huh.StateNormal
	}
}

func (a *App) clearActiveForm() tea.Cmd {
	switch a.menuList.SelectedItem().(item).calc {
	case calcProductivity:
		a.prodForm = createProductivityForm()
		return wrapCmd(a.prodForm.Init())
	case calcReliability:
		a.relForm = createReliabilityForm()
		return wrapCmd(a.relForm.Init())
	case calcFinOps:
		a.finForm = createFinOpsForm()
		return wrapCmd(a.finForm.Init())
	case calcSRE:
		a.sreForm = createSREForm()
		return wrapCmd(a.sreForm.Init())
	case calcOnboarding:
		a.onboardingForm = createOnboardingForm()
		return wrapCmd(a.onboardingForm.Init())
	case calcContextSwitch:
		a.contextSwitchForm = createContextSwitchForm()
		return wrapCmd(a.contextSwitchForm.Init())
	case calcCostOfDelay:
		a.costOfDelayForm = createCostOfDelayForm()
		return wrapCmd(a.costOfDelayForm.Init())
	}
	return nil
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
			a.prodForm = createProductivityForm()
			a.relForm = createReliabilityForm()
			a.finForm = createFinOpsForm()
			a.sreForm = createSREForm()
			a.onboardingForm = createOnboardingForm()
			a.contextSwitchForm = createContextSwitchForm()
			a.costOfDelayForm = createCostOfDelayForm()
			cmds = append(cmds, wrapCmd(a.prodForm.Init()), wrapCmd(a.relForm.Init()), wrapCmd(a.finForm.Init()), wrapCmd(a.sreForm.Init()), wrapCmd(a.onboardingForm.Init()), wrapCmd(a.contextSwitchForm.Init()), wrapCmd(a.costOfDelayForm.Init()))
			a.resultText = ""
		}
	case focusForm:
		selectedItem := a.menuList.SelectedItem().(item)
		form := a.getActiveForm()

		newModel, newCmd := form.Update(mapV2MsgToV1(msg))
		if f, ok := newModel.(*huh.Form); ok {
			a.setActiveForm(selectedItem.calc, f)
			if f.State == huh.StateCompleted {
				a.calculateResult(selectedItem.calc)
			}
		}
		cmds = append(cmds, wrapCmd(newCmd))
	}

	return a, tea.Batch(cmds...)
}

func (a *App) getActiveForm() *huh.Form {
	selectedItem := a.menuList.SelectedItem().(item)
	switch selectedItem.calc {
	case calcProductivity:
		return a.prodForm
	case calcReliability:
		return a.relForm
	case calcFinOps:
		return a.finForm
	case calcSRE:
		return a.sreForm
	case calcOnboarding:
		return a.onboardingForm
	case calcContextSwitch:
		return a.contextSwitchForm
	case calcCostOfDelay:
		return a.costOfDelayForm
	}
	return nil
}

func (a *App) setActiveForm(calc calcType, f *huh.Form) {
	switch calc {
	case calcProductivity:
		a.prodForm = f
	case calcReliability:
		a.relForm = f
	case calcFinOps:
		a.finForm = f
	case calcSRE:
		a.sreForm = f
	case calcOnboarding:
		a.onboardingForm = f
	case calcContextSwitch:
		a.contextSwitchForm = f
	case calcCostOfDelay:
		a.costOfDelayForm = f
	}
}

// calculateResult dispatches the result calculation to the appropriate domain handler.
func (a *App) calculateResult(calc calcType) {
	switch calc {
	case calcProductivity:
		a.calcProductivityResult()
	case calcReliability:
		a.calcReliabilityResult()
	case calcFinOps:
		a.calcFinOpsResult()
	case calcSRE:
		a.calcSREResult()
	case calcOnboarding:
		a.calcOnboardingResult()
	case calcContextSwitch:
		a.calcContextSwitchResult()
	case calcCostOfDelay:
		a.calcCostOfDelayResult()
	}
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
