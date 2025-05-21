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
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

const (
	errorMsg       = "error"
	systemNotFound = "not found"
)

var (
	ErrCouldntParseHTMLTemplate        = errors.New("couldn't parse HTML template")
	ErrCouldntParseBuiltInHTMLTemplate = errors.New("couldn't parse built-in HTML template")
	ErrCouldntCreateTable              = errors.New("couldn't create data table")
	errCouldntRenderTable              = errors.New("couldn't render table")
	ErrCouldntRenderHTML               = errors.New("couldn't render HTML")
)

//go:embed assets/template.html
var builtInHTMLTemplate string

func GetOutput(config VersionsUIConfig,
	versionResults map[string]map[string]types.VersionResult,
	changesResults []types.ChangesResult,
) (string, error) {
	switch config.OutputFmt {
	case types.TabularFmt:
		return getTabularOutput(config, versionResults)
	case types.HTMLFmt:
		return getHTMLOutput(config, versionResults, changesResults)
	default:
		return getTerminalOutput(config, versionResults), nil
	}
}

func getTabularOutput(config VersionsUIConfig, results map[string]map[string]types.VersionResult) (string, error) {
	rows := make([][]string, len(results))

	for _, sys := range config.SystemKeys {
		var row []string

		var versions []versionInfo
		for _, env := range config.EnvSequence {
			r, ok := results[sys][env]
			if !ok {
				versions = append(versions, versionInfo{})
				continue
			}
			if r.Err != nil {
				versions = append(versions, versionInfo{errMsg: errorMsg})
			} else {
				if !r.Found {
					versions = append(versions, versionInfo{notFound: true})
				} else {
					versions = append(versions, versionInfo{version: r.Version, registeredAt: r.RegisteredAt})
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
		for _, v := range versions {
			if v.errMsg != "" {
				row = append(row, v.errMsg)
			} else if v.notFound {
				row = append(row, systemNotFound)
			} else if v.version == "" {
				row = append(row, "")
			} else {
				if config.ShowRegisteredAt {
					duration := int(time.Since(*v.registeredAt).Seconds())
					durationMsg := fmt.Sprintf("(%s ago)", HumanizeDuration(duration))
					row = append(row, fmt.Sprintf("%s %s", v.version, durationMsg))
				} else {
					row = append(row, v.version)
				}
			}
		}
		rows = append(rows, row)
	}

	var headers []string
	headers = append(headers, "system")
	headers = append(headers, "in-sync")
	headers = append(headers, config.EnvSequence...)

	var style tw.BorderStyle
	switch config.Style {
	case types.BlankStyle:
		style = tw.StyleNone
	case types.DotsStyle:
		style = tw.StyleDotted
	case types.SharpStyle:
		style = tw.StyleDefault
	default:
		style = tw.StyleASCII
	}

	b := bytes.Buffer{}
	table := tablewriter.NewTable(&b,
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{
					Alignment:  tw.AlignRight,
					AutoWrap:   tw.WrapNone,
					AutoFormat: tw.Off,
				},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{
					Alignment: tw.AlignRight,
					AutoWrap:  tw.WrapNone,
				},
			},
		}),
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{Symbols: tw.NewSymbols(style)})),
		tablewriter.WithHeader(headers),
	)

	if err := table.Bulk(rows); err != nil {
		return "", fmt.Errorf("%w: %w", errCouldntRenderTable, err)
	}

	if err := table.Render(); err != nil {
		return "", fmt.Errorf("%w: %s", errCouldntRenderTable, err.Error())
	}

	return b.String(), nil
}

type versionInfo struct {
	version      string
	errMsg       string
	registeredAt *time.Time
	notFound     bool
}

func getTerminalOutput(config VersionsUIConfig, results map[string]map[string]types.VersionResult) string {
	var s string

	s += "\n"
	s += " " + headerStyle.Render("ecsv")
	s += "\n\n"

	var envSt lipgloss.Style
	var resultSt lipgloss.Style

	if config.ShowRegisteredAt {
		envSt = envStyle
		resultSt = resultStyle
	} else {
		envSt = envStyle.Width(18)
		resultSt = resultStyle.Width(22)
	}

	s += systemStyle.Render("system")

	for _, env := range config.EnvSequence {
		s += fmt.Sprintf("%s    ", envSt.Render(env))
	}
	s += "\n\n"
	errorIndex := 0
	var errors []error

	for _, sys := range config.SystemKeys {
		s += systemStyle.Render(sys)
		var versions []versionInfo
		for _, env := range config.EnvSequence {
			r, ok := results[sys][env]
			if !ok {
				versions = append(versions, versionInfo{})
				continue
			}
			if r.Err != nil {
				versions = append(versions, versionInfo{errMsg: fmt.Sprintf("%s [%d]", errorMsg, errorIndex)})
				errors = append(errors, r.Err)
				errorIndex++
			} else {
				if !r.Found {
					versions = append(versions, versionInfo{notFound: true})
				} else {
					versions = append(versions, versionInfo{version: r.Version, registeredAt: r.RegisteredAt})
				}
			}
		}

		var style lipgloss.Style
		if allEqual(versions) {
			style = inSyncStyle
		} else {
			style = outOfSyncStyle
		}

		for _, v := range versions {
			if v.errMsg != "" {
				s += resultSt.Render(errorStyle.Render(v.errMsg))
			} else if v.notFound {
				s += resultSt.Render(errorStyle.Render(systemNotFound))
			} else if v.version == "" {
				s += resultSt.Render("")
			} else {
				if config.ShowRegisteredAt {
					duration := int(time.Since(*v.registeredAt).Seconds())
					durationMsg := fmt.Sprintf("(%s ago)", HumanizeDuration(duration))
					s += resultSt.Render(fmt.Sprintf("%s %s", style.Render(v.version), durationStyle.Render(durationMsg)))
				} else {
					s += resultSt.Render(style.Render(v.version))
				}
			}
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

func getHTMLOutput(config VersionsUIConfig,
	versionResults map[string]map[string]types.VersionResult,
	changesResults []types.ChangesResult,
) (string, error) {
	var columns []string
	rows := make([]VersionRow, len(config.SystemKeys))

	data := HTMLData{
		Title:    config.HTMLTitle,
		TitleURL: config.HTMLTitleURL,
		Changes:  changesResults,
	}

	columns = append(columns, "system")
	columns = append(columns, config.EnvSequence...)

	errorIndex := 0
	var errors []error
	for i, sys := range config.SystemKeys {
		var rowData []string
		rowData = append(rowData, sys)
		var versions []versionInfo
		var inSync bool
		for _, env := range config.EnvSequence {
			r, ok := versionResults[sys][env]
			if !ok {
				versions = append(versions, versionInfo{})
				continue
			}
			if r.Err != nil {
				versions = append(versions, versionInfo{errMsg: fmt.Sprintf("%s [%d]", errorMsg, errorIndex)})
				errorIndex++
				errors = append(errors, r.Err)
				inSync = false
			} else {
				if !r.Found {
					versions = append(versions, versionInfo{notFound: true})
				} else {
					versions = append(versions, versionInfo{version: r.Version, registeredAt: r.RegisteredAt})
				}
			}
		}

		if allEqual(versions) {
			inSync = true
		}
		for _, v := range versions {
			if v.errMsg != "" {
				rowData = append(rowData, v.errMsg)
			} else if v.notFound {
				rowData = append(rowData, systemNotFound)
			} else if v.version == "" {
				rowData = append(rowData, "")
			} else {
				if config.ShowRegisteredAt {
					duration := int(time.Since(*v.registeredAt).Seconds())
					durationMsg := fmt.Sprintf("(%s ago)", HumanizeDuration(duration))
					rowData = append(rowData, fmt.Sprintf("%s %s", v.version, durationMsg))
				} else {
					rowData = append(rowData, v.version)
				}
			}
		}

		rows[i] = VersionRow{
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
		if err != nil {
			return "", fmt.Errorf("%w: %s", ErrCouldntParseHTMLTemplate, err.Error())
		}
	} else {
		tmpl, err = template.New("ecsv").Parse(builtInHTMLTemplate)
		if err != nil {
			return "", fmt.Errorf("%w: %s", ErrCouldntParseBuiltInHTMLTemplate, err.Error())
		}
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrCouldntRenderHTML, err.Error())
	}

	return buf.String(), nil
}
