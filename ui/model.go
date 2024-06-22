package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	outFormat       OutFormat
	envSequence     []string
	results         map[string]map[string]string
	systemNames     []string
	systems         []System
	awsConfigs      map[string]AWSConfig
	numResultsToGet int
	numResults      int
	printWhenReady  bool
	outputPrinted   bool
	htmlTemplate    string
	errors          []error
	terminalWidth   int
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, system := range m.systems {
		cmds = append(cmds, m.fetchSystemVersion(system))
	}
	return tea.Batch(cmds...)
}
