package typeTransform

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgType"
	"reflect"
)

//golang的kind分的太细
type Kind uint

const (
	Invalid Kind = iota
	String
	Int
	Float
	Ptr
	Bool
	Time
	Interface
	Map
	Struct
	Slice
	Array
)

type TransformerFunc func(traner Transformer, in reflect.Value, out reflect.Value) (err error)
type Transformer map[Kind]map[Kind]TransformerFunc

func (t Transformer) Transform(in interface{}, out interface{}) (err error) {
	return t.Tran(reflect.ValueOf(in), reflect.ValueOf(out))
}
func (t Transformer) Tran(in reflect.Value, out reflect.Value) (err error) {
	iKind := GetReflectKind(in)
	oKind := GetReflectKind(out)
	m1, ok := t[iKind]
	if !ok {
		return fmt.Errorf("[typeTransform.tran] not support tran kind: [%s] to [%s]", in.Kind(), out.Kind())
	}
	m2, ok := m1[oKind]
	if !ok {
		return fmt.Errorf("[typeTransform.tran] not support tran kind: [%s] to [%s]", in.Kind(), out.Kind())
	}
	return m2(t, in, out)
}
func (t Transformer) Clone() Transformer {
	out1 := Transformer{}
	for inKind, m1 := range t {
		out2 := map[Kind]TransformerFunc{}
		for outKind, m2 := range m1 {
			out2[outKind] = m2
		}
		out1[inKind] = out2
	}
	return out1
}
func GetReflectKind(in reflect.Value) Kind {
	t := in.Type()
	if t == kmgType.DateTimeReflectType {
		return Time
	}
	switch t.Kind() {
	case reflect.String:
		return String
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr:
		return Int
	case reflect.Float32, reflect.Float64:
		return Float
	case reflect.Ptr:
		return Ptr
	case reflect.Bool:
		return Bool
	case reflect.Interface:
		return Interface
	case reflect.Map:
		return Map
	case reflect.Struct:
		return Struct
	case reflect.Slice:
		return Slice
	case reflect.Array:
		return Array
	}
	panic(fmt.Errorf("not implement type %s", t.Kind().String()))
}

var DefaultTransformer = Transformer{
	Map: map[Kind]TransformerFunc{
		Map:    MapToMap,
		Struct: MapToStruct,
		Ptr:    NonePtrToPtr,
	},
	String: map[Kind]TransformerFunc{
		String: StringToString,
		Int:    StringToInt,
		Float:  StringToFloat,
		Bool:   StringToBool,
		Time:   StringToTime,
		Ptr:    NonePtrToPtr,
	},
	Ptr: map[Kind]TransformerFunc{
		Ptr: PtrToPtr, //TODO reference to self..
	},
	Slice: map[Kind]TransformerFunc{
		Slice: SliceToSlice,
		Ptr:   NonePtrToPtr,
	},
	Interface: map[Kind]TransformerFunc{
		String: InterfaceToNoneInterface,
		Int:    InterfaceToNoneInterface,
		Float:  InterfaceToNoneInterface,
		Bool:   InterfaceToNoneInterface,
		Time:   InterfaceToNoneInterface,
		Struct: InterfaceToNoneInterface,
		Map:    InterfaceToNoneInterface,
		Ptr:    InterfaceToNoneInterface,
	},
	Int: map[Kind]TransformerFunc{
		Int: IntToInt,
		Ptr: NonePtrToPtr,
	},
	Float: map[Kind]TransformerFunc{
		Int:   FloatToInt,
		Float: FloatToFloat,
		Ptr:   NonePtrToPtr,
	},
}
