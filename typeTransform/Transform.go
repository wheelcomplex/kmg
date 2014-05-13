package typeTransform

import (
	"fmt"
	//"github.com/bronze1man/kmg/kmgReflect"
	"github.com/bronze1man/kmg/kmgTime"
	//"github.com/bronze1man/kmg/kmgType"
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
	return DefaultTransformer.Transform(in,out)
}

func MapToMap(t Transformer,in reflect.Value, out reflect.Value) (err error) {
	out.Set(reflect.MakeMap(out.Type()))
	for _, key := range in.MapKeys() {
		oKey := reflect.New(out.Type().Key()).Elem()
		oVal := reflect.New(out.Type().Elem()).Elem()
		err = t.Tran(key, oKey)
		if err != nil {
			return
		}
		val := in.MapIndex(key)
		err = t.Tran(val, oVal)
		if err != nil {
			return
		}
		out.SetMapIndex(oKey, oVal)
	}
	return
}

func StringToString(t Transformer,in reflect.Value,out reflect.Value)(err error){
	out.SetString(in.String())
	return nil
}

func StringToTime(traner Transformer,in reflect.Value,out reflect.Value)(err error){
	var t time.Time
	t, err = kmgTime.ParseAutoInLocal(in.String())
	if err != nil {
		return
	}
	out.Set(reflect.ValueOf(t))
	return
}

func PtrToPtr(t Transformer,in reflect.Value,out reflect.Value)(err error){
	t.Tran(in.Elem(), out.Elem())
	return
}

func MapToStruct(t Transformer,in reflect.Value, out reflect.Value) (err error) {
	oKey := reflect.New(reflect.TypeOf("")).Elem()
	out.Set(reflect.New(out.Type()).Elem())
	fieldNameMap := map[string]bool{}
	for _, key := range in.MapKeys() {
		err = t.Tran(key, oKey)
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
		err = t.Tran(val, oVal)
		if err != nil {
			return
		}
	}
	return
}
func SliceToSlice(t Transformer,in reflect.Value, out reflect.Value) (err error) {
	len := in.Len()
	out.Set(reflect.MakeSlice(out.Type(), len, len))
	for i := 0; i < len; i++ {
		val := in.Index(i)
		err = t.Tran(val, out.Index(i))
		if err != nil {
			return
		}
	}
	return
}

// "" => 0
func StringToInt(t Transformer,in reflect.Value, out reflect.Value) (err error) {
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
func StringToFloat(t Transformer,in reflect.Value, out reflect.Value) (err error) {
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
func StringToBool(t Transformer,in reflect.Value, out reflect.Value) (err error) {
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

func IntToInt(t Transformer,in reflect.Value,out reflect.Value)(err error){
	out.SetInt(in.Int())
	return nil
}

func FloatToInt(t Transformer,in reflect.Value,out reflect.Value)(err error){
	outf1 := in.Float()
	if math.Floor(outf1) != outf1 {
		return fmt.Errorf("[typeTransform.tran] it seems to lose some accuracy trying to convert from float to int,float:%f", outf1)
	}
	out.SetInt(int64(outf1))
	return
}

func FloatToFloat(t Transformer,in reflect.Value,out reflect.Value)(err error){
	out.SetFloat(in.Float())
	return
}

func NonePtrToPtr(t Transformer,in reflect.Value,out reflect.Value)(err error){
	if out.IsNil() {
		out.Set(reflect.New(out.Type().Elem()))
	}
	return t.Tran(in, out.Elem())
}
func InterfaceToNoneInterface(t Transformer,in reflect.Value,out reflect.Value)(err error){
	return t.Tran(in.Elem(), out)
}
