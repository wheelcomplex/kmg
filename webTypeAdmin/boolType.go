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
	return theTemplate.MustExecuteNameToHtml("Select", selectTemplateData{
		List:  []string{"false", "true"},
		Value: t.toString(v),
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

func (t *boolType) fromString(s string) (reflect.Value, error) {
	valueT, err := strconv.ParseBool(s)
	if err != nil {
		return reflect.Value{}, err
	}
	rv := reflect.New(t.getReflectType()).Elem()
	rv.SetBool(valueT)
	return rv, nil
}
func (t *boolType) toString(v reflect.Value) string {
	return strconv.FormatBool(v.Bool())
}

func (t *boolType) Save(v *reflect.Value, path Path, value string) error {
	if err := scaleValueSaveHandle(t, &v, path); err != nil {
		return err
	}
	valueT, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	v.SetBool(valueT)
	return nil
}
