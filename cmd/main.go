// ©AngelaMos | 2026 "yoshi"
// main.go

package main

import (
	"fmt"
	"os"

	ui "github.com/CarterPerez-dev/yoshi-audit/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	app := ui.NewApp()
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
