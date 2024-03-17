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

const (
	templateText = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://cdn.tailwindcss.com"></script>
    <title>{{.Title}}</title>
</head>
<body class="bg-slate-900 text-xl">
    <div class="container mx-auto p-4">
        <h1 class="text-stone-50 text-2xl mb-4 font-bold">{{.Title}}</h1>
        <br>
        <table class="table-auto font-bold text-left">
            <thead>
                <tr class="text-stone-50 bg-slate-700">
                    {{range .Columns -}}
                    <th class="px-4 py-2">{{.}}</th>
                    {{end -}}
                </tr>
            </thead>
            <tbody>
                {{range .Rows -}}
                    {{if .InSync}}
                <tr class="text-green-600">
                    {{else}}
                <tr class="text-red-600">
                    {{end}}
                    {{range .Data -}}
                    <td class="px-4 py-2">{{.}}</td>
                    {{end -}}
                </tr>
                {{end -}}
            </tbody>
        </table>
        <br>
        <br>
        <p class="text-stone-300 italic">Generated at {{.Timestamp}}</p>

        {{if .Errors }}
        <br>
        <hr>
        <br>
        <p class="text-red-600 font-bold italic">Errors</p>
            <br>
            {{range $index, $error := .Errors -}}
            <p class="text-gray-400 italic">{{$index}}: {{$error}}</p>
            <br>
            {{end -}}
        {{end -}}
    </div>
</body>
</html>
`
	errorTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
<p>Something went wrong generating the HTML</p>
<p>Error: %s</p>
</body>
</html>
`
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
			versions = append(versions, m.results[sys][env])
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
			versions = append(versions, m.results[sys][env])
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

	tmpl, err := template.New("ecsv").Parse(templateText)
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
