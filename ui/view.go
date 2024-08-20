package ui

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/ecsv/internal/types"
	"github.com/olekukonko/tablewriter"
)

const (
	errorMsg       = "error"
	systemNotFound = "not found"
)

var (
	errCouldntParseHTMLTemplate = errors.New("couldn't parse HTML template")
	errCouldntRenderHTML        = errors.New("couldn't render HTML")
)

//go:embed assets/template.html
var htmlTemplate string

func GetOutput(config Config, results map[string]map[string]types.SystemResult) (string, error) {
	switch config.OutputFmt {
	case types.TabularFmt:
		return getTabularOutput(config, results), nil
	case types.HTMLFmt:
		return getHTMLOutput(config, results)
	default:
		return getTerminalOutput(config, results), nil
	}
}

func getTabularOutput(config Config, results map[string]map[string]types.SystemResult) string {
	rows := make([][]string, len(results))

	for _, sys := range config.SystemKeys {
		var row []string

		var versions []string
		for _, env := range config.EnvSequence {
			r, ok := results[sys][env]
			if !ok {
				versions = append(versions, "")
				continue
			}
			if r.Err != nil {
				versions = append(versions, errorMsg)
			} else {
				if !r.Found {
					versions = append(versions, systemNotFound)
				} else {
					versions = append(versions, r.Version)
				}
			}
		}
		var inSync string
		if allEqual(versions) {
			inSync = "YES"
		} else {
			inSync = "NO"
		}
		row = append(row, sys)
		row = append(row, inSync)
		row = append(row, versions...)
		rows = append(rows, row)
	}

	b := bytes.Buffer{}
	table := tablewriter.NewWriter(&b)

	var headers []string
	headers = append(headers, "system")
	headers = append(headers, "in-sync")
	headers = append(headers, config.EnvSequence...)
	table.SetHeader(headers)

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.AppendBulk(rows)

	table.Render()

	return b.String()
}

func getTerminalOutput(config Config, results map[string]map[string]types.SystemResult) string {
	var s string

	s += "\n"
	s += " " + headerStyle.Render("ecsv")
	s += "\n\n"

	s += systemStyle.Render("system")

	for _, env := range config.EnvSequence {
		s += fmt.Sprintf("%s    ", envStyle.Render(env))
	}
	s += "\n\n"
	errorIndex := 0
	var errors []error

	for _, sys := range config.SystemKeys {
		s += systemStyle.Render(sys)
		var versions []string
		var style lipgloss.Style
		styleSet := false
		for _, env := range config.EnvSequence {
			r, ok := results[sys][env]
			if !ok {
				versions = append(versions, "")
				continue
			}
			if r.Err != nil {
				versions = append(versions, fmt.Sprintf("%s [%d]", errorMsg, errorIndex))
				errors = append(errors, r.Err)
				errorIndex++
				style = errorStyle
				styleSet = true
			} else {
				if !r.Found {
					versions = append(versions, systemNotFound)
				} else {
					versions = append(versions, r.Version)
				}
			}
		}

		if !styleSet {
			if allEqual(versions) {
				style = inSyncStyle
			} else {
				style = outOfSyncStyle
			}
		}
		for _, v := range versions {
			s += fmt.Sprintf("%s    ", style.Render(v))
		}
		s += "\n"
	}

	if len(errors) > 0 {
		s += "\n"
		s += errorHeadingStyle.Render("Errors")
		s += "\n"
		for index, err := range errors {
			s += errorDetailStyle.Render(fmt.Sprintf("[%d]: %s", index, err.Error()))
			s += "\n"
		}
	}
	return s
}

func getHTMLOutput(config Config, results map[string]map[string]types.SystemResult) (string, error) {
	var columns []string
	rows := make([]HTMLDataRow, len(config.SystemKeys))

	data := HTMLData{
		Title: "ecsv",
	}

	columns = append(columns, "system")
	columns = append(columns, config.EnvSequence...)

	errorIndex := 0
	var errors []error
	for i, sys := range config.SystemKeys {
		var rowData []string
		rowData = append(rowData, sys)
		var versions []string
		var inSync bool
		for _, env := range config.EnvSequence {
			r, ok := results[sys][env]
			if !ok {
				versions = append(versions, "")
				continue
			}
			if r.Err != nil {
				versions = append(versions, fmt.Sprintf("%s [%d]", errorMsg, errorIndex))
				errorIndex++
				errors = append(errors, r.Err)
				inSync = false
			} else {
				if !r.Found {
					versions = append(versions, systemNotFound)
				} else {
					versions = append(versions, r.Version)
				}
			}
		}

		if allEqual(versions) {
			inSync = true
		}
		rowData = append(rowData, versions...)

		rows[i] = HTMLDataRow{
			Data:   rowData,
			InSync: inSync,
		}
	}
	data.Columns = columns
	data.Rows = rows
	if len(errors) > 0 {
		data.Errors = errors
	}
	data.Timestamp = time.Now().Format("2006-01-02 15:04:05 MST")

	var tmpl *template.Template
	var err error
	if config.HTMLTemplate != "" {
		tmpl, err = template.New("ecsv").Parse(config.HTMLTemplate)
	} else {
		tmpl, err = template.New("ecsv").Parse(htmlTemplate)
	}
	if err != nil {
		return "", fmt.Errorf("%w: %s", errCouldntParseHTMLTemplate, err.Error())
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errCouldntRenderHTML, err.Error())
	}

	return buf.String(), nil
}

func allEqual(versions []string) bool {
	if len(versions) <= 1 {
		return true
	}
	var firstNonEmpty string
	for _, v := range versions {
		if v == errorMsg || v == systemNotFound {
			return false
		}
		if v != "" {
			firstNonEmpty = v
			break
		}
	}
	if firstNonEmpty == "" {
		return true
	}

	for _, v := range versions[1:] {
		if v == errorMsg || v == systemNotFound {
			return false
		}
		if v != firstNonEmpty || v == errorMsg || v == systemNotFound {
			return false
		}
	}
	return true
}
