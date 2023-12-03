# {{ .Id }}

{{ .Year }}年度 {{ if eq .Degree "bachelor" }}卒業論文{{ else }}修士論文{{ end }}

## Dissertation 情報

* **タイトル**: {{ .Title }}
* **著者**: {{ .Author.Name }}
* **所属**: {{ .Author.Affiliation }}
* **指導教員**: {{- range $a := .Supervisors }}
  1. {{ $a }}
{{- end }}
