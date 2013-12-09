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
	f, err := strconv.ParseFloat(value, t.getReflectType().Bits())
	if err != nil {
		return err
	}
	v.SetFloat(f)
	return nil
}

func (t *floatType) fromString(s string) (reflect.Value, error) {
	valueT, err := strconv.ParseFloat(s, t.getReflectType().Bits())
	if err != nil {
		return reflect.Value{}, err
	}
	rv := reflect.New(t.getReflectType()).Elem()
	rv.SetFloat(valueT)
	return rv, nil
}
func (t *floatType) toString(v reflect.Value) string {
	return strconv.FormatFloat(v.Float(), 'g', -1, t.getReflectType().Bits())
}

func (t *floatType) Save(v *reflect.Value, path Path, value string) error {
	if err := scaleValueSaveHandle(t, &v, path); err != nil {
		return err
	}
	f, err := strconv.ParseFloat(value, t.getReflectType().Bits())
	if err != nil {
		return err
	}
	v.SetFloat(f)
	return nil
}
