package typeTransform

import "reflect"

func ToString(in interface{}) (out string, err error) {
	err = Transform(in, out)
	return
}

func ToStringReflect(in reflect.Value) (out string, err error) {
	err = Tran(in, reflect.ValueOf(&out))
	return
}
