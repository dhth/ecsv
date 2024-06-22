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

	headerStyle = fgStyle.
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#d3869b"))

	envStyle = fgStyle.
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#b8bb26")).
			Width(16)

	envStylePlain = fgStylePlain.
			Align(lipgloss.Center).
			Width(16)

	nonFgStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	systemStyle = nonFgStyle.
			Align(lipgloss.Left).
			Bold(true).
			Foreground(lipgloss.Color("#83a598")).
			Width(30)

	systemStylePlain = nonFgStyle.
				Align(lipgloss.Left).
				Width(30)

	inSyncStyle = nonFgStyle.
			Align(lipgloss.Center).
			Bold(true).
			Foreground(lipgloss.Color("#b8bb26")).
			Width(16)

	inSyncStylePlain = nonFgStyle.
				Align(lipgloss.Center).
				Width(16)

	outOfSyncStyle = nonFgStyle.
			Align(lipgloss.Center).
			Bold(true).
			Foreground(lipgloss.Color("#fb4934")).
			Width(16).
			Underline(true)

	outOfSyncStylePlain = nonFgStyle.
				Align(lipgloss.Center).
				Width(16)

	errorHeadingStyle = nonFgStyle.
				Bold(true).
				Foreground(lipgloss.Color("#fb4934"))

	errorDetailStyle = nonFgStyle.
				Foreground(lipgloss.Color("#665c54"))
)
