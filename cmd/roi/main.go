package main

import (
	"fmt"
	"os"

	"github.com/MichielVanderhoydonck/roi/internal/tui"
)

func main() {
	app := tui.NewApp()
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
