package webTypeAdmin

import (
	"reflect"
	//"time"
	"fmt"
)

//fkRefType need this stuff to get rootValue dependency
type context struct {
	rootType  typeInterface
	rootValue reflect.Value
}

func (ctx *context) newTypeFromReflect(rt reflect.Type) (t typeInterface, err error) {
	t, err = ctx.newBasicTypeFromReflect(rt)
	if err != nil {
		return nil, err
	}
	if rt.Implements(fkRefReflectType) {
		ut := t
		st, ok := ut.(stringConverterType)
		if !ok {
			panic("[fkRefType.new] underlyiny type is not stringConverterType")
		}
		t = &fkRefType{
			commonType:          commonType{rt, ctx},
			underlyingType:      ut,
			stringConverterType: st,
		}
	}
	return t, nil
}
func (ctx *context) newBasicTypeFromReflect(rt reflect.Type) (t typeInterface, err error) {
	ct := commonType{rt, ctx}
	switch rt {
	case dateTimeReflectType:
		t = &dateTimeType{commonType: ct}
	default:
		switch rt.Kind() {
		case reflect.Ptr:
			t = &ptrType{commonType: ct}
		case reflect.Bool:
			t = &boolType{commonType: ct}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Uintptr:
			t = &intType{commonType: ct}
		case reflect.Float32, reflect.Float64:
			t = &floatType{commonType: ct}
		case reflect.String:
			if rt.Implements(stringEnumReflectType) {
				t = &stringEnumType{commonType: ct}
			} else {
				t = &stringType{commonType: ct}
			}
		case reflect.Array:
			t = &arrayType{commonType: ct}
		case reflect.Slice:
			t = &sliceType{commonType: ct}
		case reflect.Map:
			t = &mapType{commonType: ct}
		case reflect.Struct:
			t = &structType{commonType: ct}
		default:
			return nil, fmt.Errorf("not support type kind: %s", rt.Kind().String())
		}
	}
	//TODO Recursion checkout type
	return t, nil
}
func (ctx *context) mustNewTypeFromReflect(rt reflect.Type) typeInterface {
	t, err := ctx.newTypeFromReflect(rt)
	if err != nil {
		panic(err)
	}
	return t
}
