package typeTransform

import (
	"fmt"
	"reflect"
	"strconv"
)

/*
try best transform one type to another type
special case:
"" => 0
*/
func Transform(in interface{}, out interface{}) (err error) {
	return Tran(reflect.ValueOf(in), reflect.ValueOf(out))
}

func Tran(in reflect.Value, out reflect.Value) (err error) {
	if in.Type() == out.Type() && out.CanSet() {
		out.Set(in)
		return nil
	}
	switch in.Kind() {
	case reflect.Map:
		switch out.Kind() {
		case reflect.Map:
			return MapToMap(in, out)
		case reflect.Struct:
			return MapToStruct(in, out)
		}
	case reflect.String:
		switch out.Kind() {
		case reflect.String:
			out.SetString(in.String())
			return
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Uintptr:
			return StringToInt(in, out)
		}
	case reflect.Ptr:
		switch out.Kind() {
		case reflect.Ptr:
			return Tran(in.Elem(), out.Elem())
		}
	case reflect.Slice:
		switch out.Kind() {
		case reflect.Slice:
			return SliceToSlice(in, out)
		}
	case reflect.Interface:
		switch out.Kind() {
		case reflect.Interface:
			return Tran(in.Elem(), out.Elem())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr:
		switch out.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Uintptr:
			out.SetInt(in.Int())
			return
		}
	}
	if out.Kind() == reflect.Ptr {
		return Tran(in, out.Elem())
	}
	return fmt.Errorf("[typeTransform.tran] not support tran kind: %s %s", in.Kind(), out.Kind())
}
func MapToMap(in reflect.Value, out reflect.Value) (err error) {
	oKey := reflect.New(out.Type().Key()).Elem()
	oVal := reflect.New(out.Type().Elem()).Elem()
	out.Set(reflect.MakeMap(out.Type()))
	for _, key := range in.MapKeys() {
		err = Tran(key, oKey)
		if err != nil {
			return
		}
		val := in.MapIndex(key)
		err = Tran(val, oVal)
		if err != nil {
			return
		}
		out.SetMapIndex(oKey, oVal)
	}
	return
}
func MapToStruct(in reflect.Value, out reflect.Value) (err error) {
	oKey := reflect.New(reflect.TypeOf("")).Elem()
	out.Set(reflect.New(out.Type()).Elem())
	for _, key := range in.MapKeys() {
		err = Tran(key, oKey)
		if err != nil {
			return
		}
		sKey := oKey.String()
		oVal := out.FieldByName(sKey)
		if !oVal.IsValid() {
			continue
		}
		val := in.MapIndex(key)
		err = Tran(val, oVal)
		if err != nil {
			return
		}
	}
	return
}
func SliceToSlice(in reflect.Value, out reflect.Value) (err error) {
	len := in.Len()
	out.Set(reflect.MakeSlice(out.Type(), len, len))
	for i := 0; i < len; i++ {
		val := in.Index(i)
		err = Tran(val, out.Index(i))
		if err != nil {
			return
		}
	}
	return
}

// "" => 0
func StringToInt(in reflect.Value, out reflect.Value) (err error) {
	inS := in.String()
	if inS == "" {
		out.SetInt(int64(0))
		return nil
	}
	i, err := strconv.ParseInt(inS, 10, out.Type().Bits())
	if err != nil {
		return
	}
	out.SetInt(i)
	return
}
