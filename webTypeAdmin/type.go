package webTypeAdmin

import (
	"reflect"
	//"github.com/bronze1man/kmg/kmgReflect"
	"fmt"
	"html/template"
	"time"
)

type typeInterface interface {
	getReflectType() reflect.Type
	Html(v reflect.Value) template.HTML // component html
	getSubValueByString(v reflect.Value, k string) (reflect.Value, error)
	delete(v reflect.Value, k string) error //delete some elem in slice or map
	create(v reflect.Value, k string) error
	save(v reflect.Value, value string) error
}

func newTypeFromReflect(rt reflect.Type) (t typeInterface, err error) {
	switch reflect.Zero(rt).Interface().(type) {
	case time.Time:
		t = &dateTimeType{commonType: commonType{rt}}
	default:
		switch rt.Kind() {
		case reflect.Ptr:
			t = &ptrType{commonType: commonType{rt}}
		case reflect.Bool:
			t = &boolType{commonType: commonType{rt}}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			t = &intType{commonType: commonType{rt}}
		case reflect.Float32, reflect.Float64:
			t = &floatType{commonType: commonType{rt}}
		case reflect.String:
			t = &stringType{commonType: commonType{rt}}
		case reflect.Array:
			t = &arrayType{commonType: commonType{rt}}
		case reflect.Slice:
			t = &sliceType{commonType: commonType{rt}}
		case reflect.Map:
			t = &mapType{commonType: commonType{rt}}
		case reflect.Struct:
			t = &structType{commonType: commonType{rt}}
		default:
			return nil, fmt.Errorf("not support type kind: %s", rt.Kind().String())
		}
	}
	//TODO Recursion checkout type
	return t, nil
}

func mustNewTypeFromReflect(rt reflect.Type) typeInterface {
	t, err := newTypeFromReflect(rt)
	if err != nil {
		panic(err)
	}
	return t
}

type commonType struct {
	reflect reflect.Type
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

//xxId refer to another object
//path -> pass to RefType
type fkRefType struct {
	commonType
}
type enumType struct {
	commonType
}
