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

func (t *sliceType) init() (err error) {
	if t.elemType != nil {
		return
	}
	t.elemType, err = t.ctx.typeOfFromReflect(t.GetReflectType().Elem())
	return
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
		if html, err = t.elemType.HtmlView(v.Index(i)); err != nil {
			return
		}
		templateData = append(templateData, templateRow{
			Path:  i,
			Index: i,
			Html:  html,
		})
	}
	return theTemplate.ExecuteNameToHtml("Slice", templateData)
}
