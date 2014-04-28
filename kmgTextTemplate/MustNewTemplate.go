package kmgTextTemplate

import (
	"text/template"
)

func MustNewTemplate(text string) (tmpl *template.Template) {
	return template.Must(template.New("").Parse(text))
}
