## {{.Version}} [{{.Repo}}-release-notes-{{.Version}}]

{{ if or .Feature .Enhancement }}
### Features and enhancements [{{.Repo}}-{{.Version}}-features-enhancements]
{{ if .Feature }}{{ range $k, $v := .Feature }}{{ range $item := $v }}
* {{ $item.Summary | beautify }} {{ linkPRSource $item.Component $item.LinkedPR }} {{ linkIssueSource $item.Component $item.LinkedIssue }}
{{- end }}{{- end }}{{- end }}{{ if .Enhancement }}{{ range $k, $v := .Enhancement }}{{ range $item := $v }}
* {{ $item.Summary | beautify }} {{ linkPRSource $item.Component $item.LinkedPR }} {{ linkIssueSource $item.Component $item.LinkedIssue }}
{{- end }}{{- end }}{{- end }}
{{- end }}

{{ if or .Security .BugFix }}
### Fixes [{{.Repo}}-{{.Version}}-fixes]
{{ if .Security }}{{ range $k, $v := .Security }}{{ range $item := $v }}
* {{ $item.Summary | beautify }} {{ linkPRSource $item.Component $item.LinkedPR }} {{ linkIssueSource $item.Component $item.LinkedIssue }}
{{- end }}{{- end }}{{- end }}{{ if .BugFix }}{{ range $k, $v := .BugFix }}{{ range $item := $v }}
* {{ $item.Summary | beautify }} {{ linkPRSource $item.Component $item.LinkedPR }} {{ linkIssueSource $item.Component $item.LinkedIssue }}
{{- end }}{{- end }}{{- end }}
{{- end }}
