# ROI Calculator

A cleanly architected, terminal-based user interface (TUI) for calculating returns on investment (ROI) across various business and engineering domains.

## Features & Use Cases

The application currently supports three main calculators:

1. **Developer Productivity ROI**: Calculates the gross and net savings of developer productivity initiatives. Enter the time saved per execution, executions per year, hourly rate, and maintenance cost to see the annual impact.
2. **Reliability ROI**: Calculates the cost savings of avoiding downtime. It compares the old Mean Time To Recovery (MTTR) with the new MTTR, taking into account the incident frequency and downtime cost per hour.
3. **FinOps ROI**: Calculates annual cloud savings by comparing old and new monthly cloud infrastructure bills.

## Architecture & Design

This project is built using **Go** and adheres to clean architecture principles:

* **`cmd/roi`**: The application entry point.
* **`internal/domain`**: Contains the core business logic, models, and interfaces for the different ROI calculators (Productivity, Reliability, FinOps).
* **`internal/service`**: Implements the domain interfaces, housing the specific calculation logic for each domain.
* **`internal/tui`**: The presentation layer. It leverages the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework for the overarching application state and layout, and [Huh](https://github.com/charmbracelet/huh) for the interactive forms. The UI is designed to be highly responsive, readable, and navigable entirely via the keyboard.

## Development Workflow

A `Makefile` is provided to simplify common development tasks.

* **Run the application**:
  ```bash
  make run
  ```
* **Build the binary**:
  ```bash
  make build
  ```
  The compiled executable will be placed in `bin/roi`.
* **Run tests**:
  ```bash
  make test
  ```
* **Format and lint code**:
  ```bash
  make fmt
  make vet
  ```
* **Clean build artifacts**:
  ```bash
  make clean
  ```

## Usage

Once running, use your keyboard to navigate:
* `up/down` or `j/k`: Navigate the main menu.
* `enter` / `tab` / `right`: Select a calculator and focus the form.
* `esc`: Go back to the main menu.
* `ctrl+r`: Clear the active form.
* `ctrl+c`: Quit the application.
