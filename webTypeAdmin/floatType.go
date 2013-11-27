package webTypeAdmin

import (
	"html/template"
	"reflect"
	"strconv"
)

type floatType struct {
	commonType
}

func (t *floatType) Html(v reflect.Value) template.HTML {
	return theTemplate.MustExecuteNameToHtml("TextInput", v.Float())
}

func (t *floatType) save(v reflect.Value, value string) error {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(f))
	return nil
}
