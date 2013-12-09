package kmgType

import (
	"fmt"
	"reflect"
)

//path -> array index(number)
type ArrayType struct {
	reflectTypeGetterImp
	getElemByStringEditorabler
	elemType KmgType
}

func (t *ArrayType) init() (err error) {
	if t.elemType != nil {
		return
	}
	t.elemType, err = TypeOf(t.GetReflectType().Elem())
	return
}
func (t *ArrayType) GetElemByString(v reflect.Value, k string) (ev reflect.Value, et KmgType, err error) {
	if err = t.init(); err != nil {
		return
	}
	et = t.elemType
	ev, err = arrayGetSubValueByString(v, k)
	if err != nil {
		return
	}
	return
}

func (t *ArrayType) DeleteByPath(v *reflect.Value, path Path) (err error) {
	if len(path) > 1 {
		return passThougthDeleteByPath(t, v, path)
	}
	return fmt.Errorf("can not delete from array, path:%s", path)
}
