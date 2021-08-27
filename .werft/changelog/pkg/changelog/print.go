package changelog

import (
	"bytes"
	"strings"
	"text/template"
)

const title = `# Change Log
`

const templ = `{{ .Title }}

{{ if .Notes }}
{{ range .Notes }}
* {{ .MergeDay }} {{ .Description }} [[#{{ .Num }}]({{ .URI }})] - [{{ .Author }}]({{ .AuthorURL }})
{{ end }}
{{ end }}

{{ if .ExistingNotes }}
{{ .ExitingNotes }}
{{ end }}
`

type templateData struct {
	Notes         []ReleaseNote
	ExistingNotes string
	Title         string
}

// Print ...
func Print(notes []ReleaseNote, existingNotes string) (string, error) {
	if strings.HasPrefix(existingNotes, title) {
		existingNotes = existingNotes[len(title):]
	}
	data := templateData{
		Notes:         notes,
		ExistingNotes: existingNotes,
		Title:         title,
	}

	t := template.New("changes")
	res, err := t.Parse(templ)
	if err != nil {
		return "", err
	}

	b := bytes.NewBuffer(nil)
	err = res.Execute(b, data)
	if err != nil {
		return "", err
	}

	result := strings.ReplaceAll(b.String(), "\n\n", "\n")

	return result, nil
}
