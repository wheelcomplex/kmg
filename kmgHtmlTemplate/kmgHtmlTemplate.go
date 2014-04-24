package kmgHtmlTemplate

import (
	"bytes"
	"fmt"
	"html/template"
)

type Template struct {
	*template.Template
}

func (templ *Template) ExecuteNameToByte(name string, data interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := templ.ExecuteTemplate(buf, name, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (templ *Template) MustExecuteNameToHtml(name string, data interface{}) template.HTML {
	b, err := templ.ExecuteNameToByte(name, data)
	if err != nil {
		panic(err)
	}
	return template.HTML(b)
}

func (templ *Template) ExecuteNameToHtml(name string, data interface{}) (html template.HTML, err error) {
	b, err := templ.ExecuteNameToByte(name, data)
	if err != nil {
		return
	}
	return template.HTML(b), nil
}

func (templ *Template) MustExecuteToHtml(data interface{}) template.HTML {
	h, err := templ.ExecuteToByte(data)
	if err != nil {
		panic(err)
	}
	return template.HTML(h)
}

func MustNewSingle(templ string) *Template {
	return &Template{template.Must(template.New("single").Parse(templ))}
}

func (templ *Template) ExecuteToByte(data interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := templ.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type TemplateFile struct {
	Name    string
	Content string
}

//something like template.ParseFile but in memory
func FileSet(templs []TemplateFile) (ot *Template, err error) {
	var t *template.Template
	for _, file := range templs {
		if t == nil {
			t, err = template.New(file.Name).Parse(file.Content)
			if err != nil {
				return nil, fmt.Errorf("[kmgHtmlTemplate.FileSet] template.Parse name[%s] err:%s", file.Name, err)
			}
			continue
		}
		_, err = t.New(file.Name).Parse(file.Content)
		if err != nil {
			return nil, fmt.Errorf("[kmgHtmlTemplate.FileSet] template.Parse name[%s] err:%s", file.Name, err)
		}
	}
	return &Template{t}, nil
}

func MustFileSet(templs []TemplateFile) *Template {
	t, err := FileSet(templs)
	if err != nil {
		panic(err)
	}
	return t
}
