package ui

const (
	HTMLTemplText = `
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
