package webTypeAdmin

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgReflect"
	"html/template"
	"reflect"
)

//path -> field name
type structType struct {
	commonType
	fields []structField //here need a stable order.
}
type structField struct {
	Name string
	Type typeInterface
	*reflect.StructField
}

func (t *structType) init() {
	if t.fields != nil {
		return
	}
	for _, v := range kmgReflect.StructGetAllField(t.getReflectType()) {
		sf := structField{
			Name:        v.Name,
			Type:        mustNewTypeFromReflect(v.Type),
			StructField: v,
		}
		t.fields = append(t.fields, sf)
	}
}
func (t *structType) Html(v reflect.Value) template.HTML {
	t.init()
	type templateRow struct {
		Path string
		Name string
		Html template.HTML
	}
	var templateData []templateRow
	for _, sf := range t.fields {
		templateData = append(templateData, templateRow{
			Path: sf.Name,
			Name: sf.Name,
			Html: sf.Type.Html(v.FieldByIndex(sf.Index)),
		})
	}
	return theTemplate.MustExecuteNameToHtml("Struct", templateData)
}
func (t *structType) getSubValueByString(v reflect.Value, k string) (reflect.Value, error) {
	ev := v.FieldByName(k)
	if !ev.IsValid() {
		return reflect.Value{}, fmt.Errorf("field %s not find in struct", k)
	}
	return ev, nil
}
