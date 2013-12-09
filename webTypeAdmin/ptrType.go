package webTypeAdmin

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgType"
	"html/template"
	"reflect"
)

//just for same with golang reflect system
//path -> pass to RefType
type ptrType struct {
	kmgType.PtrType
	elemType adminType
	ctx      *context
}

func (t *ptrType) init() (err error) {
	if t.elemType != nil {
		return
	}
	t.elemType, err = t.ctx.typeOfFromReflect(t.GetReflectType().Elem())
	return
}
func (t *ptrType) HtmlView(v reflect.Value) (html template.HTML, err error) {
	if err = t.init(); err != nil {
		return
	}
	if v.IsNil() {
		return theTemplate.ExecuteNameToHtml("NilPtr", nil)
	}
	elemHtml, err := t.elemType.Html(v.Elem())
	if err != nil {
		return
	}
	return theTemplate.ExecuteNameToHtml("Ptr", elemHtml)
}
