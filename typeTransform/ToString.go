package typeTransform

import "reflect"

func ToString(in interface{}) (out string, err error) {
	err = DefaultTransformer.Transform(in, &out)
	return
}

func ToStringReflect(in reflect.Value) (out string, err error) {
	err = DefaultTransformer.Tran(in, reflect.ValueOf(&out))
	return
}
