package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func RenderUI(envSequence []string, systems []System, outFormat OutFormat) {
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
		opts = []tea.ProgramOption{tea.WithoutRenderer()}
	}
	p := tea.NewProgram(newModel(envSequence, systems, outFormat), opts...)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting Bubble Tea program:", err)
		os.Exit(1)
	}

}
