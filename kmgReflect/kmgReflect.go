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
func GetTypeFullName(t reflect.Type) (name string, ok bool) {
	if t.Kind() == reflect.Ptr {
		return GetTypeFullName(t.Elem())
	}
	if t.Name() == "" {
		return "", false
	}

	if t.PkgPath() == "" {
		return t.Name(), true
	}
	return t.PkgPath() + "." + t.Name(), true
}
