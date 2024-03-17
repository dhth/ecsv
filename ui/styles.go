package ui

import "github.com/charmbracelet/lipgloss"

var (
	fgStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(lipgloss.Color("#282828"))

	fgStylePlain = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	headerStyle = fgStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#d3869b"))

	headerStylePlain = fgStylePlain.Copy().
				Align(lipgloss.Center)

	envStyle = fgStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#b8bb26")).
			Width(16)

	envStylePlain = fgStylePlain.Copy().
			Align(lipgloss.Center).
			Width(16)

	nonFgStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	systemStyle = nonFgStyle.Copy().
			Align(lipgloss.Left).
			Bold(true).
			Foreground(lipgloss.Color("#83a598")).
			Width(30)

	systemStylePlain = nonFgStyle.Copy().
				Align(lipgloss.Left).
				Width(30)

	inSyncStyle = nonFgStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Foreground(lipgloss.Color("#b8bb26")).
			Width(16)

	inSyncStylePlain = nonFgStyle.Copy().
				Align(lipgloss.Center).
				Width(16)

	outOfSyncStyle = nonFgStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Foreground(lipgloss.Color("#fb4934")).
			Width(16).
			Underline(true)

	outOfSyncStylePlain = nonFgStyle.Copy().
				Align(lipgloss.Center).
				Width(16)

	errorHeadingStyle = nonFgStyle.Copy().
				Bold(true).
				Foreground(lipgloss.Color("#fb4934"))

	errorDetailStyle = nonFgStyle.Copy().
				Foreground(lipgloss.Color("#665c54"))
)
