package ui

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	ErrorFetchingVersion = "error"
	SystemNotFound       = "not found"
)

func (m model) renderPlainText() string {
	var s string

	s += " " + "ecsv"
	s += "\n\n"

	s += fmt.Sprintf("%s", systemStylePlain.Render("system"))

	for _, env := range m.envSequence {
		s += fmt.Sprintf("%s    ", envStylePlain.Render(env))
	}
	s += "\n\n"
	for _, sys := range m.systemNames {
		s += fmt.Sprintf("%s", systemStylePlain.Render(sys))
		var style lipgloss.Style
		var versions []string
		for _, env := range m.envSequence {
			if m.results[sys][env] != "" {
				versions = append(versions, m.results[sys][env])
			}
		}
		if allEqual(versions) {
			style = inSyncStylePlain
		} else {
			style = outOfSyncStylePlain
		}
		for _, env := range m.envSequence {
			s += fmt.Sprintf("%s    ", style.Render(m.results[sys][env]))
		}
		s += "\n"
	}

	if len(m.errors) > 0 {
		s += "\n"
		s += "Errors"
		s += "\n"
		for index, err := range m.errors {
			s += fmt.Sprintf("[%2d]: %s", index+1, err.Error())
			s += "\n"
		}
	}
	return s
}

func (m model) renderHTML() string {

	var columns []string
	var rows []HTMLDataRow
	var inSync bool

	data := HTMLData{
		Title: "ecsv",
	}

	columns = append(columns, "system")
	for _, env := range m.envSequence {
		columns = append(columns, env)
	}

	for _, sys := range m.systemNames {
		var rowData []string
		rowData = append(rowData, sys)
		var versions []string
		for _, env := range m.envSequence {
			if m.results[sys][env] != "" {
				versions = append(versions, m.results[sys][env])
			}
		}
		if allEqual(versions) {
			inSync = true
		} else {
			inSync = false
		}
		for _, env := range m.envSequence {
			rowData = append(rowData, m.results[sys][env])
		}
		rows = append(rows, HTMLDataRow{
			Data:   rowData,
			InSync: inSync,
		})
	}
	data.Columns = columns
	data.Rows = rows
	if len(m.errors) > 0 {
		data.Errors = &m.errors
	}
	data.Timestamp = time.Now().Format("2006-01-02 15:04:05 MST")

	var tmpl *template.Template
	var err error
	if m.htmlTemplate == "" {
		tmpl, err = template.New("ecsv").Parse(HTMLTemplText)
	} else {
		tmpl, err = template.New("ecsv").Parse(m.htmlTemplate)
	}
	if err != nil {
		return fmt.Sprintf(string(errorTemplate), err.Error())
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Sprintf(string(errorTemplate), err.Error())
	}

	return buf.String()
}

func (m model) renderCLIUI() string {
	var s string

	s += "\n"
	s += " " + headerStyle.Render("ecsv")
	s += "\n\n"

	s += fmt.Sprintf("%s", systemStyle.Render("system"))

	for _, env := range m.envSequence {
		s += fmt.Sprintf("%s    ", envStyle.Render(env))
	}
	s += "\n\n"
	for _, sys := range m.systemNames {
		s += fmt.Sprintf("%s", systemStyle.Render(sys))
		var style lipgloss.Style
		var versions []string
		for _, env := range m.envSequence {
			if m.results[sys][env] != "" {
				versions = append(versions, m.results[sys][env])
			}
		}
		if allEqual(versions) {
			style = inSyncStyle
		} else {
			style = outOfSyncStyle
		}
		for _, env := range m.envSequence {
			s += fmt.Sprintf("%s    ", style.Render(m.results[sys][env]))
		}
		s += "\n"
	}

	if len(m.errors) > 0 {
		s += "\n"
		s += errorHeadingStyle.Render("Errors")
		s += "\n"
		for index, err := range m.errors {
			s += errorDetailStyle.Render(fmt.Sprintf("[%2d]: %s", index+1, err.Error()))
			s += "\n"
		}
	}
	return s
}

func (m model) View() string {
	return m.renderCLIUI()
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
