package typeTransform

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgReflect"
	"github.com/bronze1man/kmg/kmgTime"
	"github.com/bronze1man/kmg/kmgType"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/*
try best transform one type to another type
special case:
"" => 0
"" => 0.0
*/
func Transform(in interface{}, out interface{}) (err error) {
	return Tran(reflect.ValueOf(in), reflect.ValueOf(out))
}

func Tran(in reflect.Value, out reflect.Value) (err error) {
	/*
		if in.Type() == out.Type() && out.CanSet() {
			out.Set(in)
			return nil
		}
	*/
	//fmt.Printf("%#v\n",out.Interface())
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
		case reflect.Float32, reflect.Float64:
			return StringToFloat(in, out)
		case reflect.Bool:
			return StringToBool(in, out)
		}
		if out.Type() == kmgType.DateTimeReflectType {
			var t time.Time
			t, err = kmgTime.ParseAutoInLocal(in.String())
			if err != nil {
				return
			}
			out.Set(reflect.ValueOf(t))
			return
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
		return Tran(in.Elem(), out)
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
	case reflect.Float64, reflect.Float32:
		switch out.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Uintptr:
			outf1 := in.Float()
			if math.Floor(outf1) != outf1 {
				return fmt.Errorf("[typeTransform.tran] it seems to lose some accuracy trying to convert from float to int,float:%f", outf1)
			}
			out.SetInt(int64(outf1))
			return
		case reflect.Float64, reflect.Float32:
			out.SetFloat(in.Float())
			return
		}
	}
	if out.Kind() == reflect.Ptr {
		if out.IsNil() {
			out.Set(reflect.New(out.Type().Elem()))
		}
		return Tran(in, out.Elem())
	}
	return fmt.Errorf("[typeTransform.tran] not support tran kind: [%s] to [%s]", in.Kind(), out.Kind())
}
func MapToMap(in reflect.Value, out reflect.Value) (err error) {
	out.Set(reflect.MakeMap(out.Type()))
	for _, key := range in.MapKeys() {
		oKey := reflect.New(out.Type().Key()).Elem()
		oVal := reflect.New(out.Type().Elem()).Elem()
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

/*
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
		//TODO remove this workaroud
		if oVal.Kind() == reflect.Float64 {
			f, ok := out.Type().FieldByName(sKey)
			if !ok {
				return
			}
			div := f.Tag.Get("floatDiv")
			if div != "" {
				divI, err := strconv.ParseFloat(div, 64)
				if err != nil {
					return err
				}
				oVal.SetFloat(oVal.Float() / divI)
			}
		}
	}
	return
} */

//every field in struct must in the map
//TODO 用户可以配置这个东西...
func MapToStruct(in reflect.Value, out reflect.Value) (err error) {
	oKey := reflect.New(reflect.TypeOf("")).Elem()
	out.Set(reflect.New(out.Type()).Elem())
	fieldNameMap := map[string]bool{}
	for _, key := range in.MapKeys() {
		err = Tran(key, oKey)
		if err != nil {
			return
		}
		sKey := oKey.String()
		fieldNameMap[sKey] = true
		oVal := out.FieldByName(sKey)
		if !oVal.IsValid() {
			continue
		}
		val := in.MapIndex(key)
		err = Tran(val, oVal)
		if err != nil {
			return
		}
		//TODO remove this workaroud
		if oVal.Kind() == reflect.Float64 {
			f, ok := out.Type().FieldByName(sKey)
			if !ok {
				return
			}
			div := f.Tag.Get("floatDiv")
			if div != "" {
				divI, err := strconv.ParseFloat(div, 64)
				if err != nil {
					return err
				}
				oVal.SetFloat(oVal.Float() / divI)
			}
		}
	}
	//if false {
	for _, structField := range kmgReflect.StructGetAllField(out.Type()) {
		if structField.Anonymous {
			continue
		}
		if !fieldNameMap[structField.Name] {
			return fmt.Errorf("[MapToStruct]type:%s field:%s not found", out.Type().Name(), structField.Name)
		}
	}
	//}
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
	inS = strings.TrimSpace(inS)
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

// "" => 0.0
func StringToFloat(in reflect.Value, out reflect.Value) (err error) {
	inS := in.String()
	if inS == "" {
		out.SetFloat(0.0)
		return nil
	}
	i, err := strconv.ParseFloat(inS, out.Type().Bits())
	if err != nil {
		return
	}
	out.SetFloat(i)
	return
}

// "" => false
func StringToBool(in reflect.Value, out reflect.Value) (err error) {
	inS := in.String()
	if inS == "" {
		out.SetBool(false)
		return nil
	}
	i, err := strconv.ParseBool(inS)
	if err != nil {
		return
	}
	out.SetBool(i)
	return
}
