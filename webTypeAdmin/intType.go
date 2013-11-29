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

func (t *intType) fromString(s string) (reflect.Value, error) {
	valueT, err := strconv.ParseInt(s, 10, t.getReflectType().Bits())
	if err != nil {
		return reflect.Value{}, err
	}
	rv := reflect.New(t.getReflectType()).Elem()
	rv.SetInt(valueT)
	return rv, nil
}
func (t *intType) toString(v reflect.Value) string {
	return strconv.FormatInt(v.Int(), t.getReflectType().Bits())
}
