package template

func ReturnBaseTemplate() string {
	return ``
}

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
