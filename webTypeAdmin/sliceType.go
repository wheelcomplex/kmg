package webTypeAdmin

import (
	"fmt"
	"html/template"
	"reflect"
	"strconv"
)

//path -> slice index(number)
type sliceType struct {
	commonType
	elemType typeInterface
}

func (t *sliceType) init() {
	if t.elemType != nil {
		return
	}
	t.elemType = mustNewTypeFromReflect(t.getReflectType().Elem())
}
func (t *sliceType) Html(v reflect.Value) template.HTML {
	t.init()
	type templateRow struct {
		Path int
		Html template.HTML
	}
	var templateData []templateRow
	len := v.Len()
	for i := 0; i < len; i++ {
		templateData = append(templateData, templateRow{
			Path: i,
			Html: t.elemType.Html(v.Index(i)),
		})
	}
	return theTemplate.MustExecuteNameToHtml("Slice", templateData)
}

func (t *sliceType) getSubValueByString(v reflect.Value, k string) (reflect.Value, error) {
	t.init()
	i, err := t.parseKey(v, k)
	if err != nil {
		return reflect.Value{}, nil
	}
	return v.Index(i), nil
}
func (t *sliceType) delete(v reflect.Value, k string) error {
	t.init()
	i, err := t.parseKey(v, k)
	if err != nil {
		return err
	}
	v.Set(
		reflect.AppendSlice(v.Slice(0, i), v.Slice(i+1, v.Len())),
	)
	return nil
}
func (t *sliceType) create(v reflect.Value, k string) error {
	t.init()
	v.Set(
		reflect.Append(v, reflect.New(t.elemType.getReflectType()).Elem()),
	)
	return nil
}
func (t *sliceType) parseKey(v reflect.Value, k string) (int, error) {
	i64, err := strconv.ParseInt(k, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("[sliceType.parseKey] index is not int k:%s", k)
	}
	i := int(i64)
	if i >= v.Len() || i < 0 {
		return 0, fmt.Errorf("[sliceType.parseKey] index is not of range k:%s,len:%d", k, v.Len())
	}
	return i, nil
}
