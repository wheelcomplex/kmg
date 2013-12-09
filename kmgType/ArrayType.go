package kmgType

import (
	"reflect"
)

//path -> array index(number)
type ArrayType struct {
	reflectTypeGetterImp
	getElemByStringEditorabler
}

func (t *ArrayType) GetElemByString(v reflect.Value, k string) (ev reflect.Value, et KmgType, err error) {
	et, err = TypeOf(t.GetReflectType().Elem())
	if err != nil {
		return
	}
	ev, err = arrayGetSubValueByString(v, k)
	if err != nil {
		return
	}
	return
}
