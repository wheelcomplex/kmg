package webTypeAdmin

import (
	"github.com/bronze1man/kmg/kmgReflect"
	"github.com/bronze1man/kmg/kmgType"
	"html/template"
	"reflect"
	//"fmt"
)

//path -> field name
//TODO embed field
type structType struct {
	kmgType.StructType
	ctx    *context
	fields []structField //here need a stable order.
}
type structField struct {
	Name string
	Type adminType
	*reflect.StructField
}

func (t *structType) init() (err error) {
	if t.fields != nil {
		return
	}
	//t.fieldsMap = map[string]*structField{}
	for _, v := range kmgReflect.StructGetAllField(t.GetReflectType()) {
		at, err := t.ctx.typeOfFromReflect(v.Type)
		if err != nil {
			return err
		}
		sf := structField{
			Name:        v.Name,
			Type:        at,
			StructField: v,
		}
		t.fields = append(t.fields, sf)
	}
	return nil
}
func (t *structType) HtmlView(v reflect.Value) (html template.HTML, err error) {
	err = t.init()
	if err != nil {
		return
	}
	type templateRow struct {
		Path string
		Name string
		Html template.HTML
	}
	var templateData []templateRow
	for _, sf := range t.fields {
		var thisHtml template.HTML
		thisHtml, err = sf.Type.HtmlView(v.FieldByIndex(sf.Index))
		if err != nil {
			return
		}
		templateData = append(templateData, templateRow{
			Path: sf.Name,
			Name: sf.Name,
			Html: thisHtml,
		})
	}
	return theTemplate.ExecuteNameToHtml("Struct", templateData)
}
