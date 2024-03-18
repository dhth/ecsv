package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func RenderUI(envSequence []string, systems []System, outFormat OutFormat, htmlTemplate string) {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	var opts []tea.ProgramOption
	if outFormat != UnspecifiedFmt {
		opts = append(opts, tea.WithoutRenderer())
		// TODO: this may be a hack, and will prevent using STDIN for
		// CLI mode, find a better way
		opts = append(opts, tea.WithInput(nil))
	}
	p := tea.NewProgram(newModel(envSequence, systems, outFormat, htmlTemplate), opts...)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting Bubble Tea program:", err)
		os.Exit(1)
	}

}
