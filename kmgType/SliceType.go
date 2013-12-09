package kmgType

import (
	"reflect"
)

//path -> slice index(number)
type SliceType struct {
	reflectTypeGetterImp
	getElemByStringEditorabler
}

func (t *SliceType) GetElemValueByString(v reflect.Value, k string) (ev reflect.Value, et KmgType, err error) {
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
