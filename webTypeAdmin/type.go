package webTypeAdmin

import (
	"github.com/bronze1man/kmg/kmgType"
	"html/template"
	"reflect"
)

type adminType interface {
	kmgType.KmgType
	HtmlView(v reflect.Value) (html template.HTML, err error) // component html
}

//int float datatime string
type toStringTextHtmlView struct {
	kmgType.KmgTypeAndToStringInterface
}

func (t toStringTextHtmlView) HtmlView(v reflect.Value) (html template.HTML, err error) {
	return theTemplate.ExecuteNameToHtml("TextInput", t.ToString(v))
}

//bool fkRef map
type selectTextHtmlView struct {
	List []string
	kmgType.KmgTypeAndToStringInterface
}

func (t selectTextHtmlView) HtmlView(v reflect.Value) (html template.HTML, err error) {
	return theTemplate.ExecuteNameToHtml("Select", selectTemplateData{
		List:  t.List,
		Value: t.ToString(v),
	})
}

/*
type typeInterface interface {
	getReflectType() reflect.Type
	Html(v reflect.Value) template.HTML // component html
	getSubValueByString(v reflect.Value, k string) (reflect.Value, error)
	delete(v reflect.Value, k string) error //delete some elem in slice or map
	create(v reflect.Value, k string) error
	save(v reflect.Value, value string) error
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

func (t *commonType) Delete(inV reflect.Value, path Path) (*reflect.Value, error) {
	return nil, fmt.Errorf("[commonType.Delete] Can not create while type:%s path:%s", t.getReflectType().Kind(), path)
}
func (t *commonType) Save(v *reflect.Value, path Path, value string) error {
	return fmt.Errorf("[commonType.Save] Can not save while type:%s path:%s value:%s", t.getReflectType().Kind(), path, value)
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

func scaleValueSaveHandle(t typeInterface, v **reflect.Value, path Path) error {
	if len(path) != 0 {
		return fmt.Errorf("[scaleValueSaveHandle] get string with some path,path:%s", path)
	}
	if !*v.CanSet() {
		**v = reflect.New(t.getReflectType()).Elem()
	}
	return nil
}

type HtmlViewType interface {
	HtmlView(v reflect.Value) template.HTML // component html
}
type StringConverterType interface {
	fromString(s string) (reflect.Value, error)
	toString(v reflect.Value) string //is caller responsibility to ensure v is callee Type
}

//bool datetime fkref float int string stringEnum
type ScaleType interface {
	SaveScale(v reflect.Value, value string) error
}

type EditAbleType interface {
	// v must not be nil, v.IsValid must be true
	// every type must consider reflect.Value.CanSet()
	// every type in the path except least one(include ptr),must eat a path.
	// try best to create all path, if it is not exist.
	// if this value can set its own value,set it,if this value can not set its own value,replace v with a new value
	Save(v *reflect.Value, path Path, value string) (err error)
	Delete(v *reflect.Value, path Path) (err error)
}

//array slice struct map
type ElemByStringGetter interface {
	getSubValueByString(v reflect.Value, k string) (reflect.Value, error)
}

type AdminType struct {
	commonType
	HtmlViewType
	StringConverterType
	EditAbleType
}
*/
