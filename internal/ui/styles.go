package ui

import "github.com/charmbracelet/lipgloss"

var (
	fgStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(lipgloss.Color("#282828"))

	headerStyle = fgStyle.
			Align(lipgloss.Left).
			Bold(true).
			Background(lipgloss.Color("#d3869b"))

	envStyle = fgStyle.
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#b8bb26")).
			Width(30)

	nonFgStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	systemStyle = nonFgStyle.
			Align(lipgloss.Left).
			Bold(true).
			Foreground(lipgloss.Color("#83a598")).
			Width(30)

	versionStyle = nonFgStyle.
			Width(18).
			Align(lipgloss.Left).
			Bold(true).
			Foreground(lipgloss.Color("#b8bb26"))

	resultStyle = lipgloss.NewStyle().
			Width(34)

	durationStyle = lipgloss.NewStyle().
			Width(12).
			Foreground(lipgloss.Color("#665c54"))

	inSyncStyle = versionStyle.
			Foreground(lipgloss.Color("#b8bb26"))

	outOfSyncStyle = versionStyle.
			Foreground(lipgloss.Color("#fb4934")).
			Underline(true)

	errorStyle = versionStyle.
			Foreground(lipgloss.Color("#fabd2f")).
			Underline(true)

	errorHeadingStyle = nonFgStyle.
				Bold(true).
				Foreground(lipgloss.Color("#fb4934"))

	errorDetailStyle = nonFgStyle.
				Foreground(lipgloss.Color("#665c54"))
)
