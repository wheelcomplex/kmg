package webTypeAdmin

import (
	"html/template"
	"reflect"
)

type stringType struct {
	commonType
}

func (t *stringType) Html(v reflect.Value) template.HTML {
	return theTemplate.MustExecuteNameToHtml("TextInput", v.String())
}

func (t *stringType) save(v reflect.Value, value string) error {
	v.Set(reflect.ValueOf(value))
	return nil
}
