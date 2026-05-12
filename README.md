# ROI
Holy shit, this guy's taking ROI off the grid!

## A ROI Calculator

A terminal-based user interface (TUI) for calculating returns on investment (ROI) across various business and engineering domains.

## Features & Use Cases

The application currently supports eight main calculators:

1. **Developer Productivity ROI**: Calculates the gross and net savings of developer productivity initiatives. Enter the time saved per execution, executions per year, hourly rate, and maintenance cost to see the annual impact.
2. **Reliability ROI**: Calculates the cost savings of avoiding downtime. It compares the old Mean Time To Recovery (MTTR) with the new MTTR, taking into account the incident frequency and downtime cost per hour.
3. **FinOps ROI**: Calculates annual cloud savings by comparing old and new monthly cloud infrastructure bills.
4. **SRE Toil Eradication ROI**: Calculates savings from automating manual repetitive work based on hours per week and developer rates.
5. **Onboarding ROI**: Calculates savings from reducing time-to-first-deploy for new hires. Compare old vs. new onboarding time based on expected new hires and daily rates.
6. **Cost of Context Switching**: Calculates the hidden tax of false-positive PagerDuty alerts, flaky tests, and the cost of flow interruption based on hourly rate.
7. **Cost of Delay**: Calculates the actual revenue lost when infrastructure or architecture bottlenecks delay a feature launch.
8. **DORA AI ROI**: Analyzes the impact of AI-assisted software development tools. Calculates first-year net benefit, ROI percentage, and payback period by weighing hard costs (licenses, infrastructure, training) and J-curve productivity drops against headcount reinvestment capacity, revenue impact from faster feature delivery, and downtime cost mitigation.

## Architecture & Design

The application is structured using principles of **Clean Architecture** and **Domain-Driven Design (DDD)**, ensuring the codebase is modular, testable, and highly maintainable.

### Domain Layer (`internal/service`)
Contains the core business logic for each calculator (e.g., `FinOpsService`, `ReliabilityService`). These services define strictly typed inputs and outputs, perform the mathematical ROI calculations, and are completely decoupled from any UI concerns. This allows the business logic to maintain 100% test coverage.

### Presentation Layer (`internal/tui`)
Handles the terminal user interface using [`bubbletea`](https://github.com/charmbracelet/bubbletea) and [`huh`](https://github.com/charmbracelet/huh) for form generation. 

#### Open-Closed Principle (Polymorphism)
The TUI layer is built around a polymorphic `Calculator` interface. Instead of the main `App` maintaining complex, hardcoded switch statements for every possible ROI tool, it relies on this interface:

```go
type Calculator interface {
	CreateForm() *huh.Form
	CalculateResult(form *huh.Form) (string, Sentiment)
	GetFormula(form *huh.Form) string
	GetContext(fieldKey string) string
	Reset()
}
```

The interface supports fluid layout updates, state resetting between context switches, and dynamic color-coded feedback via `Sentiment` returns (e.g., success vs. critical alert colors).

To add a new calculator, a developer simply needs to:
1. Create a new service in `internal/service/`.
2. Define a struct implementing the `Calculator` interface in `internal/tui/`.
3. Register it in the `menuList` inside `NewApp()`.

This ensures the core `App` logic never needs to be modified when expanding the tool's capabilities.

## Development Workflow

The project includes a `Makefile` to simplify common development tasks. Alternatively, standard `go` commands can be used directly.

### Using Make

* **Build binary**: `make build` (compiles to `bin/roi`)
* **Run locally**: `make run`
* **Run tests**: `make test`
* **Format code**: `make fmt`
* **Static analysis**: `make vet`
* **Clean artifacts**: `make clean`

### Using Go Directly

```bash
# Build and execute
go build -o roi ./cmd/roi && ./roi

# Or run directly without saving binary
go run ./cmd/roi
```

## Usage

The application features a real-time updating layout. As inputs are entered or modified, the formula, context panel, and ROI calculation outcomes update instantly with color-coded feedback (green borders/text for positive ROI, red/critical for negative ROI).

### Main Menu Navigation
* `up/down` or `j/k`: Navigate the list of calculators.
* `enter` / `tab` / `right`: Select a calculator and focus its interactive input form.
* `ctrl+c`: Quit the application.

### Form Navigation & Input
* `enter` / `tab` / `down`: Move to the next input field.
* `shift+tab` / `up`: Move to the previous input field.
* `esc`: Unfocus the form and return to the main menu.
* `ctrl+r`: Reset the active form to its default values.
