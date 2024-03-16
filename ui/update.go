package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case processFinishedMsg:
		if msg.err != nil {
			m.results[msg.systemKey][msg.env] = string(ErrorFetchingVersion)
		} else {
			if !msg.found {
				m.results[msg.systemKey][msg.env] = string(SystemNotFound)
			} else {
				m.results[msg.systemKey][msg.env] = msg.version
			}
		}
		m.numResults += 1

		if m.numResults >= m.numResultsToGet {
			return m, quitProg()
		}
	case quitProgMsg:
		if !m.outputPrinted {
			switch m.outFormat {
			case PlainTextFmt:
				v := m.renderPlainText()
				fmt.Print(v)
				m.outputPrinted = true
			case HTMLFmt:
				v := m.renderHTML()
				fmt.Print(v)
				m.outputPrinted = true
			}
		}
		return m, tea.Quit
	}
	return m, nil
}
