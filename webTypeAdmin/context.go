package webTypeAdmin

import (
	"reflect"
	//"time"
	"fmt"
	"github.com/bronze1man/kmg/kmgType"
)

//fkRefType need this stuff to get rootValue dependency
type context struct {
	kmgType.Context
	RootType adminType
}

func newContext(ptr interface{}) (ctx *context, err error) {
	kmgTypeCtx, err := kmgType.NewContext(ptr)
	if err != nil {
		return
	}
	ctx = &context{Context: *kmgTypeCtx}
	ctx.RootType, err = ctx.typeOfFromKmgType(ctx.Context.RootType)
	if err != nil {
		return
	}
	return
}
func (ctx *context) typeOfFromReflect(rt reflect.Type) (t adminType, err error) {
	kt, err := kmgType.TypeOf(rt)
	if err != nil {
		return
	}
	return ctx.typeOfFromKmgType(kt)
}
func (ctx *context) typeOfFromKmgType(kt kmgType.KmgType) (t adminType, err error) {
	return ctx.basicTypeOfFromKmgType(kt)
}
func (ctx *context) basicTypeOfFromKmgType(kt kmgType.KmgType) (t adminType, err error) {
	if kt.GetReflectType().Implements(fkRefReflectType) {
		st, ok := kt.(kmgType.StringConverterInterface)
		if !ok {
			err = fmt.Errorf("[fkRefType.new] underlyiny type is not stringConverterType")
			return
		}
		t = &fkRefType{
			KmgType:            kt,
			keyStringConverter: st,
			ctx:                ctx,
		}
		return
	}
	switch ckt := kt.(type) {
	case *kmgType.DateTimeType, *kmgType.FloatType,
		*kmgType.IntType:
		return &toStringTextHtmlView{kt.(kmgType.KmgTypeAndToStringInterface)}, nil
	case *kmgType.StringType:
		if kt.GetReflectType().Implements(stringEnumReflectType) {
			return &stringEnumType{StringType: *ckt, ctx: ctx}, nil
		} else {
			return &toStringTextHtmlView{kt.(kmgType.KmgTypeAndToStringInterface)}, nil
		}
	case *kmgType.BoolType:
		return &selectTextHtmlView{
			List: []string{"false", "true"},
			KmgTypeAndToStringInterface: kt.(kmgType.KmgTypeAndToStringInterface),
		}, nil

	case *kmgType.ArrayType:
		return &arrayType{ArrayType: *ckt, ctx: ctx}, nil
	case *kmgType.MapType:
		return &mapType{MapType: *ckt, ctx: ctx}, nil
	case *kmgType.PtrType:
		return &ptrType{PtrType: *ckt, ctx: ctx}, nil
	case *kmgType.SliceType:
		return &sliceType{SliceType: *ckt, ctx: ctx}, nil
	case *kmgType.StructType:
		return &structType{StructType: *ckt, ctx: ctx}, nil
	}
	err = fmt.Errorf("not support type kind: %s", kt.GetReflectType().Kind().String())
	return
}

func (ctx *context) GetElemByPath(p kmgType.Path) (v reflect.Value, t adminType, err error) {
	v, kt, err := ctx.Context.GetElemByPath(p)
	if err != nil {
		return
	}
	t, err = ctx.typeOfFromKmgType(kt)
	if err != nil {
		return
	}
	return
}

/*
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
*/
