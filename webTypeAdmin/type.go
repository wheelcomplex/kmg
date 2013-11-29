package webTypeAdmin

import (
	"reflect"
	//"github.com/bronze1man/kmg/kmgReflect"
	"fmt"
	"html/template"
)

type typeInterface interface {
	getReflectType() reflect.Type
	Html(v reflect.Value) template.HTML // component html
	getSubValueByString(v reflect.Value, k string) (reflect.Value, error)
	delete(v reflect.Value, k string) error //delete some elem in slice or map
	create(v reflect.Value, k string) error
	save(v reflect.Value, value string) error
}

type stringConverterType interface {
	fromString(s string) (reflect.Value, error)
	toString(v reflect.Value) string //is caller responsibility to ensure v is callee Type
}
type commonType struct {
	reflect reflect.Type
	ctx     *context
}

func (t *commonType) getReflectType() reflect.Type {
	return t.reflect
}
func (t *commonType) Html(v reflect.Value) template.HTML {
	return template.HTML(
		template.HTMLEscapeString(
			fmt.Sprintf("type %s not implement.", t.getReflectType().Kind().String()),
		))
}
func (t *commonType) getSubValueByString(v reflect.Value, k string) (reflect.Value, error) {
	return reflect.Value{}, fmt.Errorf("[commonType.getSubValueByString] type %s has not subtype while key:%s", t.getReflectType().Kind().String(), k)
}
func (t *commonType) delete(v reflect.Value, k string) error {
	return fmt.Errorf("[commonType.delete] type %s can not delete while key:%s", t.getReflectType().Kind().String(), k)
}
func (t *commonType) create(v reflect.Value, k string) error {
	return fmt.Errorf("[commonType.create] type %s can not create while key:%s", t.getReflectType().Kind().String(), k)
}
func (t *commonType) save(v reflect.Value, value string) error {
	return fmt.Errorf("[commonType.save] type %s can not save while value:%s", t.getReflectType().Kind().String(), value)
}
