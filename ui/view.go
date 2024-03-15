package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	ErrorFetchingVersion = "error"
	SystemNotFound       = "not found"
)

func (m model) View() string {
	var s string

	s += "\n"
	s += " " + headerStyle.Render("ecsv")
	s += "\n\n"

	envs := []string{
		"qa",
		"staging",
		"prod",
	}
	s += fmt.Sprintf("%s", systemStyle.Render("system"))

	for _, env := range envs {
		s += fmt.Sprintf("%s    ", envStyle.Render(env))
	}
	s += "\n\n"
	for _, sys := range m.systemNames {
		s += fmt.Sprintf("%s", systemStyle.Render(sys))
		var style lipgloss.Style
		var versions []string
		for _, env := range envs {
			versions = append(versions, m.results[sys][env])
		}
		if allEqual(versions) {
			style = inSyncStyle
		} else {
			style = outOfSyncStyle
		}
		for _, env := range envs {
			s += fmt.Sprintf("%s    ", style.Render(m.results[sys][env]))
		}
		s += "\n"
	}
	return s
}

func allEqual(versions []string) bool {
	if len(versions) == 0 {
		return true
	}
	first := versions[0]
	for _, v := range versions[1:] {
		if v != first {
			return false
		}
	}
	return true
}
