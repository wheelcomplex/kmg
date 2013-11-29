package webTypeAdmin

import (
	"html/template"
	"reflect"
	//"fmt"
)

type stringType struct {
	commonType
}

func (t *stringType) Html(v reflect.Value) template.HTML {
	return theTemplate.MustExecuteNameToHtml("TextInput", v.String())
}

func (t *stringType) save(v reflect.Value, value string) error {
	v.SetString(value)
	return nil
}

func (t *stringType) fromString(s string) (reflect.Value, error) {
	rv := reflect.New(t.getReflectType()).Elem()
	rv.SetString(s)
	return rv, nil
}
func (t *stringType) toString(v reflect.Value) string {
	return v.String()
}
