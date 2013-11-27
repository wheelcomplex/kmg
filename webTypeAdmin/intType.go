package webTypeAdmin

import (
	"html/template"
	"reflect"
	"strconv"
)

type intType struct {
	commonType
}

func (t *intType) Html(v reflect.Value) template.HTML {
	return theTemplate.MustExecuteNameToHtml("TextInput", v.Int())
}

func (t *intType) save(v reflect.Value, value string) error {
	i, err := strconv.ParseInt(value, 10, t.getReflectType().Bits())
	if err != nil {
		return err
	}
	v.SetInt(i)
	return nil
}
