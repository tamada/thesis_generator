# {{ .Id }}

{{ .Year }}年度 {{ if eq .Degree "bachelor" }}卒業論文{{ else }}修士論文{{ end }}

{{ .Title }}

{{ .Author.Name }}

{{ .Author.Affiliation }}
