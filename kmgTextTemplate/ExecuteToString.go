package kmgTextTemplate

import (
	"bytes"
	"text/template"
)

func ExecuteToString(tmpl *template.Template, data interface{}) (output string, err error) {
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return
	}
	return buf.String(), nil
}
