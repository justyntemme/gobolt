package template

import (
	"fmt"
	"html/template"
)

// TODO return a template instead of a string
func ReturnBaseTemplate() string {
	return `
	<!doctype html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>{{ .Hostname }}</title>
		{{ .CSSImport }}
	</head>
	<body>
		{{ .Navigation }}
		<div>{{ .PageContent }}</div>
	</body>
	</html>
	`
}

// TODO return a template instead of a string
func ReturnNavTemplate() string {
	return `
		<nav>
			<ul>
			{{- range . }}
				<li><a href="{{ .URI }}">{{ .Title }}</a></li>
			{{- end }}
			</ul>
		</nav>
		`
}

func ReturnCSSImportTemplate(hostname string) template.HTML {
	// Generate the CSS import statement as a string
	return template.HTML(fmt.Sprintf(`<link rel="stylesheet" type="text/css" href="http://%s/css">`, hostname))
}
