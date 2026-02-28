# ROI
Holy shit, this guy's taking ROI off the grid!

## A ROI Calculator

A terminal-based user interface (TUI) for calculating returns on investment (ROI) across various business and engineering domains.

## Features & Use Cases

The application currently supports three main calculators:

1. **Developer Productivity ROI**: Calculates the gross and net savings of developer productivity initiatives. Enter the time saved per execution, executions per year, hourly rate, and maintenance cost to see the annual impact.
2. **Reliability ROI**: Calculates the cost savings of avoiding downtime. It compares the old Mean Time To Recovery (MTTR) with the new MTTR, taking into account the incident frequency and downtime cost per hour.
3. **FinOps ROI**: Calculates annual cloud savings by comparing old and new monthly cloud infrastructure bills.
4. **SRE Toil Eradication ROI**: Calculates savings from automating manual repetitive work based on hours per week and developer rates.
5. **Onboarding ROI**: Calculates savings from reducing time-to-first-deploy for new hires. Compare old vs. new onboarding time based on expected new hires and daily rates.
6. **Cost of Context Switching**: Calculates the hidden tax of false-positive PagerDuty alerts, flaky tests, and the cost of flow interruption based on hourly rate.
7. **Cost of Delay**: Calculates the actual revenue lost when infrastructure or architecture bottlenecks delay a feature launch.

## Development Workflow

See `Makefile` for common development tasks.

## Usage

Once running, use your keyboard to navigate:
* `up/down` or `j/k`: Navigate the main menu.
* `enter` / `tab` / `right`: Select a calculator and focus the form.
* `esc`: Go back to the main menu.
* `ctrl+r`: Clear the active form.
* `ctrl+c`: Quit the application.
