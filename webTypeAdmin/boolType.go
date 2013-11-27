package webTypeAdmin

import (
	"html/template"
	"reflect"
	"strconv"
)

type boolType struct {
	commonType
}

func (t *boolType) Html(v reflect.Value) template.HTML {
	var valueS string
	if v.Bool() {
		valueS = "true"
	} else {
		valueS = "false"
	}
	return theTemplate.MustExecuteNameToHtml("Select", selectTemplateData{
		List:  []string{"false", "true"},
		Value: valueS,
	})
}
func (t *boolType) save(v reflect.Value, value string) error {
	valueT, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	v.SetBool(valueT)
	return nil
}
