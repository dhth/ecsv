package ui

import "github.com/charmbracelet/lipgloss"

var (
	fgStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(lipgloss.Color("#282828"))

	headerStyle = fgStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#d3869b"))

	envStyle = fgStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#b8bb26")).
			Width(16)

	nonFgStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	systemStyle = nonFgStyle.Copy().
			Align(lipgloss.Left).
			Bold(true).
			Foreground(lipgloss.Color("#83a598")).
			Width(30)

	inSyncStyle = nonFgStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Foreground(lipgloss.Color("#b8bb26")).
			Width(16)

	outOfSyncStyle = nonFgStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Foreground(lipgloss.Color("#fb4934")).
			Width(16).
			Underline(true)
)
