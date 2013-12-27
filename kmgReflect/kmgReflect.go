package kmgReflect

import (
	"reflect"
)

// try to get full name of a type
// if can not find it ,ok will be false
// only support pointer of a struct or a struct
// &ta{} -> main.ta
// ta{} -> main.ta
// []ta{} -> "" (not support)
// for debug use fmt.Printf("%T",xxx)
func GetTypeFullName(t reflect.Type) (name string) {
	if t.Kind() == reflect.Ptr {
		return GetTypeFullName(t.Elem())
	}
	if t.Name() == "" {
		return ""
	}

	if t.PkgPath() == "" {
		return t.Name()
	}
	return t.PkgPath() + "." + t.Name()
}

func IndirectType(v reflect.Type) reflect.Type {
	switch v.Kind() {
	case reflect.Ptr:
		return IndirectType(v.Elem())
	default:
		return v
	}
	return v
}
