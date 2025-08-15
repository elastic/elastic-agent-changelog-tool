{{ if .KnownIssue -}}{{ range $k, $v := .KnownIssue }}{{ range $item := $v }}

::::{dropdown} {{ $item.Summary | beautify }}
**Applies to**: {{.Version}}

{{ if $item.Description }}{{ $item.Description }}{{ end }}

For more information, check {{ linkPRSource $item.Component $item.LinkedPR }}{{ linkIssueSource $item.Component $item.LinkedIssue }}.

{{ if not $item.Impact }}% {{ end }}**Impact**<br>{{ if $item.Impact }}{{ $item.Impact }}{{ else }}_Add a description of the impact_{{ end }}

{{ if not $item.Action }}% {{ end }}**Action**<br>{{ if $item.Action }}{{ $item.Action }}{{ else }}_Add a description of the what action to take_{{ end }}
::::
{{- end }}{{- end }}{{- end }}