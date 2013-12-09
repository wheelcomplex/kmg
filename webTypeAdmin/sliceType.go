package webTypeAdmin

import (
	"github.com/bronze1man/kmg/kmgType"
	"html/template"
	"reflect"
)

//path -> slice index(number)
type sliceType struct {
	kmgType.SliceType
	elemType adminType
	ctx      *context
}

func (t *sliceType) init() {
	if t.elemType != nil {
		return
	}
	t.elemType = t.ctx.typeOfFromReflect(t.GetReflectType().Elem())
}
func (t *sliceType) HtmlView(v reflect.Value) (html template.HTML, err error) {
	if err = t.init(); err != nil {
		return
	}
	type templateRow struct {
		Path  int
		Index int
		Html  template.HTML
	}
	var templateData []templateRow
	len := v.Len()
	for i := 0; i < len; i++ {
		templateData = append(templateData, templateRow{
			Path:  i,
			Index: i,
			Html:  t.elemType.Html(v.Index(i)),
		})
	}
	return theTemplate.ExecuteNameToTemplate("Slice", templateData)
}
