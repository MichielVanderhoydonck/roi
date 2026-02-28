package tui

import (
	"fmt"
	"strconv"
	"time"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/MichielVanderhoydonck/roi/internal/domain"
	"github.com/MichielVanderhoydonck/roi/internal/service"
	teav1 "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type focusState int

const (
	focusMenu focusState = iota
	focusForm
)

type calcType string

const (
	calcProductivity calcType = "productivity"
	calcReliability  calcType = "reliability"
	calcFinOps       calcType = "finops"
	calcSRE          calcType = "sre"
	calcOnboarding   calcType = "onboarding"
)

type item struct {
	title, desc string
	calc        calcType
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type App struct {
	prodService       *service.ProductivityService
	relService        *service.ReliabilityService
	finService        *service.FinOpsService
	sreService        *service.SREToilService
	onboardingService *service.OnboardingService

	focus focusState

	menuList list.Model

	prodForm       *huh.Form
	relForm        *huh.Form
	finForm        *huh.Form
	sreForm        *huh.Form
	onboardingForm *huh.Form

	resultText string
	width      int
	height     int
}

func NewApp() *App {
	items := []list.Item{
		item{title: "Developer Productivity", desc: "Time Saved", calc: calcProductivity},
		item{title: "Reliability", desc: "Cost of Downtime Avoided", calc: calcReliability},
		item{title: "FinOps", desc: "Infrastructure Optimization", calc: calcFinOps},
		item{title: "SRE Toil Eradication", desc: "Automating Manual Work", calc: calcSRE},
		item{title: "Onboarding ROI", desc: "Time to First Value", calc: calcOnboarding},
	}

	m := list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.Title = "ROI Calculator"
	m.SetShowHelp(false)
	m.SetShowStatusBar(false)
	m.SetFilteringEnabled(false)

	return &App{
		prodService:       service.NewProductivityService(),
		relService:        service.NewReliabilityService(),
		finService:        service.NewFinOpsService(),
		sreService:        service.NewSREToilService(),
		onboardingService: service.NewOnboardingService(),
		focus:             focusMenu,
		menuList:          m,
		prodForm:          createProductivityForm(),
		relForm:           createReliabilityForm(),
		finForm:           createFinOpsForm(),
		sreForm:           createSREForm(),
		onboardingForm:    createOnboardingForm(),
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
			cmds = append(cmds, wrapCmd(a.prodForm.Init()), wrapCmd(a.relForm.Init()), wrapCmd(a.finForm.Init()), wrapCmd(a.sreForm.Init()), wrapCmd(a.onboardingForm.Init()))
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
	}
}

// Result Calculation
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
	}
}

func (a *App) calcProductivityResult() {
	tb, _ := time.ParseDuration(a.prodForm.GetString("timeBefore"))
	ta, _ := time.ParseDuration(a.prodForm.GetString("timeAfter"))
	execs, _ := strconv.Atoi(a.prodForm.GetString("executions"))
	hr, _ := strconv.ParseFloat(a.prodForm.GetString("hourlyRate"), 64)
	mc, _ := strconv.ParseFloat(a.prodForm.GetString("maintenance"), 64)

	res := a.prodService.Calculate(domain.ProductivityInput{
		TimeBefore:        tb,
		TimeAfter:         ta,
		ExecutionsPerYear: execs,
		HourlyRate:        hr,
		MaintenanceCost:   mc,
	})

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Primary)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success)
	a.resultText = fmt.Sprintf("%s\n\nTotal Time Saved: %s\nGross Savings:    %s\nNet ROI:          %s",
		titleStyle.Render("=== Productivity ROI Results ==="),
		res.TimeSaved.String(),
		valStyle.Render(fmt.Sprintf("$%.2f", res.GrossSavings)),
		valStyle.Render(fmt.Sprintf("$%.2f", res.NetROI)))
}

func (a *App) calcReliabilityResult() {
	om, _ := time.ParseDuration(a.relForm.GetString("oldMTTR"))
	nm, _ := time.ParseDuration(a.relForm.GetString("newMTTR"))
	inc, _ := strconv.Atoi(a.relForm.GetString("incidents"))
	dc, _ := strconv.ParseFloat(a.relForm.GetString("downtimeCost"), 64)

	res := a.relService.Calculate(domain.ReliabilityInput{
		OldMTTR:          om,
		NewMTTR:          nm,
		IncidentsPerYear: inc,
		DowntimeCost:     dc,
	})

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Warning)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success)
	a.resultText = fmt.Sprintf("%s\n\nTotal Downtime avoided: %s\nDowntime Savings:       %s",
		titleStyle.Render("=== Reliability ROI Results ==="),
		res.TimeSaved.String(),
		valStyle.Render(fmt.Sprintf("$%.2f", res.DowntimeSavings)))
}

func (a *App) calcFinOpsResult() {
	ob, _ := strconv.ParseFloat(a.finForm.GetString("oldBill"), 64)
	nb, _ := strconv.ParseFloat(a.finForm.GetString("newBill"), 64)

	res := a.finService.Calculate(domain.FinOpsInput{
		OldMonthlyBill: ob,
		NewMonthlyBill: nb,
	})

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Secondary)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success)
	a.resultText = fmt.Sprintf("%s\n\nAnnual Cloud Savings: %s",
		titleStyle.Render("=== FinOps ROI Results ==="),
		valStyle.Render(fmt.Sprintf("$%.2f", res.AnnualSavings)))
}

func (a *App) calcSREResult() {
	hpw, _ := strconv.ParseFloat(a.sreForm.GetString("hoursPerWeek"), 64)
	hr, _ := strconv.ParseFloat(a.sreForm.GetString("hourlyRate"), 64)
	cta, _ := strconv.ParseFloat(a.sreForm.GetString("costToAutomate"), 64)

	res := a.sreService.Calculate(domain.SREToilInput{
		HoursPerWeek:   hpw,
		HourlyRate:     hr,
		CostToAutomate: cta,
	})

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Primary)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success)
	a.resultText = fmt.Sprintf("%s\n\nTotal Hours Saved: %.1f\nNet Savings:       %s",
		titleStyle.Render("=== SRE Toil ROI Results ==="),
		res.HoursSaved,
		valStyle.Render(fmt.Sprintf("$%.2f", res.AnnualSavings)))
}

func (a *App) calcOnboardingResult() {
	od, _ := strconv.ParseFloat(a.onboardingForm.GetString("oldDays"), 64)
	nd, _ := strconv.ParseFloat(a.onboardingForm.GetString("newDays"), 64)
	nh, _ := strconv.Atoi(a.onboardingForm.GetString("newHires"))
	dr, _ := strconv.ParseFloat(a.onboardingForm.GetString("dailyRate"), 64)

	res := a.onboardingService.Calculate(domain.OnboardingInput{
		OldDays:   od,
		NewDays:   nd,
		NewHires:  nh,
		DailyRate: dr,
	})

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(DefaultTheme.Primary)
	valStyle := lipgloss.NewStyle().Foreground(DefaultTheme.Success)
	a.resultText = fmt.Sprintf("%s\n\nDays Saved per Hire: %.1f\nIdle Time Savings:   %s",
		titleStyle.Render("=== Onboarding ROI Results ==="),
		res.DaysSavedPerHire,
		valStyle.Render(fmt.Sprintf("$%.2f", res.AnnualSavings)))
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

func (a *App) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(DefaultTheme.Title).
		Bold(true).
		Padding(0, 1)

	focusText := " [ Menu ]"
	if a.focus == focusForm {
		focusText = " [ Calculator ]"
	}

	return lipgloss.JoinHorizontal(lipgloss.Left,
		headerStyle.Render("ROI Calculator"),
		lipgloss.NewStyle().Foreground(DefaultTheme.TextDim).Render(focusText),
	)
}

func (a *App) renderFooter() string {
	footerStyle := lipgloss.NewStyle().
		Foreground(DefaultTheme.TextDim).
		Padding(0, 1)

	var helpStr string
	if a.focus == focusMenu {
		helpStr = "(esc: close • enter/tab/→: select • ctrl+c: quit)"
	} else {
		activeForm := a.getActiveForm()
		if activeForm != nil && activeForm.State == huh.StateCompleted {
			helpStr = "(esc: back to menu • enter: edit form • ctrl+r: clear form • ctrl+c: quit)"
		} else {
			helpStr = "(esc: back to menu • ctrl+r: clear form • ctrl+c: quit)"
		}
	}
	return footerStyle.Render(helpStr)
}

func (a *App) calculateLayoutWidths() (left int, right int) {
	left = max(30, min(45, a.width*30/100))
	right = max(0, a.width-left)
	return left, right
}

func (a *App) renderLeftPanel(width, height int) string {
	borderColor := DefaultTheme.BorderDim
	if a.focus == focusMenu {
		borderColor = DefaultTheme.BorderActive
	}

	style := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)

	w, h := style.GetFrameSize()
	contentWidth := max(0, width-w)
	panelHeight := max(0, height-h)

	return style.
		Width(contentWidth).
		Height(panelHeight).
		Render(a.menuList.View())
}

func (a *App) renderRightPanel(width, height int) string {
	formulaStr, contextStr := a.getCurrentPanelStrings()

	bottomHeight := max(6, height*25/100)
	topHeight := max(0, height-bottomHeight)

	activeForm := a.getActiveForm()
	topContent := a.renderTopContent(formulaStr, activeForm)
	bottomContent := a.renderBottomContent(contextStr)

	borderColor := DefaultTheme.BorderDim
	if a.focus == focusForm {
		borderColor = DefaultTheme.BorderActive
	}

	topStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)

	bottomStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(DefaultTheme.BorderDim)

	topW, topH := topStyle.GetFrameSize()
	topBox := topStyle.
		Width(max(0, width-topW)).
		Height(max(0, topHeight-topH)).
		Render(topContent)

	bottomW, bottomH := bottomStyle.GetFrameSize()
	bottomBox := bottomStyle.
		Width(max(0, width-bottomW)).
		Height(max(0, bottomHeight-bottomH)).
		Render(bottomContent)

	return lipgloss.JoinVertical(lipgloss.Left, topBox, bottomBox)
}

func (a *App) getCurrentPanelStrings() (formulaStr, contextStr string) {
	selectedItem := a.menuList.SelectedItem().(item)
	activeForm := a.getActiveForm()
	var fieldKey string

	if activeForm.State != huh.StateCompleted {
		if field := activeForm.GetFocusedField(); field != nil {
			fieldKey = field.GetKey()
		}
	}

	switch selectedItem.calc {
	case calcProductivity:
		formulaStr = getProductivityFormula(activeForm)
		if activeForm.State != huh.StateCompleted {
			contextStr = getProductivityContext(fieldKey)
		}
	case calcReliability:
		formulaStr = getReliabilityFormula(activeForm)
		if activeForm.State != huh.StateCompleted {
			contextStr = getReliabilityContext(fieldKey)
		}
	case calcFinOps:
		formulaStr = getFinOpsFormula(activeForm)
		if activeForm.State != huh.StateCompleted {
			contextStr = getFinOpsContext(fieldKey)
		}
	case calcSRE:
		formulaStr = getSREFormula(activeForm)
		if activeForm.State != huh.StateCompleted {
			contextStr = getSREContext(fieldKey)
		}
	case calcOnboarding:
		formulaStr = getOnboardingFormula(activeForm)
		if activeForm.State != huh.StateCompleted {
			contextStr = getOnboardingContext(fieldKey)
		}
	}
	return formulaStr, contextStr
}

func (a *App) renderTopContent(formulaStr string, activeForm *huh.Form) string {
	formulaStyle := lipgloss.NewStyle().
		Foreground(DefaultTheme.Title).
		Bold(true).
		MarginBottom(1)

	content := formulaStyle.Render(formulaStr) + "\n\n" + activeForm.View()

	if activeForm.State == huh.StateCompleted {
		resultBox := lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(DefaultTheme.Primary).
			Render(a.resultText)
		content += "\n" + resultBox
	}
	return content
}

func (a *App) renderBottomContent(contextStr string) string {
	contextHeader := lipgloss.NewStyle().
		Bold(true).
		Foreground(DefaultTheme.TextDim).
		Render("FIELD CONTEXT") + "\n"

	contextStyle := lipgloss.NewStyle().
		Foreground(DefaultTheme.TextNormal)

	return contextHeader + contextStyle.Render(contextStr)
}

func formatFormulaValue(val, fallback string) string {
	if val == "" {
		return lipgloss.NewStyle().Foreground(DefaultTheme.TextDim).Render(fallback)
	}
	return lipgloss.NewStyle().Foreground(DefaultTheme.Success).Render(val)
}

// Interoperability Helpers for Bubble Tea v1 and v2
func wrapCmd(cmd teav1.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}
	return func() tea.Msg {
		msg := cmd()
		if b, ok := msg.(teav1.BatchMsg); ok {
			var cmds []tea.Cmd
			for _, c := range b {
				cmds = append(cmds, wrapCmd(c))
			}
			return tea.BatchMsg(cmds)
		}

		// Map standard v1 messages back to v2 if needed,
		// but typically huh commands just return huh specific messages or nil.
		// However, teav1.Quit needs to be mapped to tea.QuitMsg
		if _, ok := msg.(teav1.QuitMsg); ok {
			return tea.Quit()
		}
		return msg
	}
}

func mapV2MsgToV1(msg tea.Msg) teav1.Msg {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		var k teav1.KeyMsg
		switch msg.String() {
		case "ctrl+c":
			k.Type = teav1.KeyCtrlC
		case "enter":
			k.Type = teav1.KeyEnter
		case "esc":
			k.Type = teav1.KeyEsc
		case "tab":
			k.Type = teav1.KeyTab
		case "shift+tab":
			k.Type = teav1.KeyShiftTab
		case "up":
			k.Type = teav1.KeyUp
		case "down":
			k.Type = teav1.KeyDown
		case "left":
			k.Type = teav1.KeyLeft
		case "right":
			k.Type = teav1.KeyRight
		case "backspace":
			k.Type = teav1.KeyBackspace
		case "delete":
			k.Type = teav1.KeyDelete
		case "space":
			k.Type = teav1.KeySpace
			k.Runes = []rune{' '}
		default:
			k.Type = teav1.KeyRunes
			k.Runes = []rune(msg.Text)
		}

		if msg.Mod.Contains(tea.ModAlt) {
			k.Alt = true
		}
		return k
	case tea.WindowSizeMsg:
		return teav1.WindowSizeMsg{Width: msg.Width, Height: msg.Height}
	}
	return msg
}
