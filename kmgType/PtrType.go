package kmgType

import (
	"fmt"
	"reflect"
)

//path -> pass to RefType
type PtrType struct {
	reflectTypeGetterImp
	elemType KmgType
}

func (t *PtrType) init() (err error) {
	if t.elemType != nil {
		return
	}
	t.elemType, err = TypeOf(t.GetReflectType().Elem())
	return
}
func (t *PtrType) GetElemType() KmgType {
	return t.elemType
}
func (t *PtrType) SaveByPath(inV *reflect.Value, path Path, value string) error {
	err := t.init()
	if err != nil {
		return err
	}
	//create
	if inV.IsNil() {
		if inV.CanSet() {
			inV.Set(reflect.New(t.GetReflectType().Elem()))
		} else {
			*inV = reflect.New(t.GetReflectType().Elem())
		}
	}
	//a elem of a ptr CanSet must be true.
	elemV := inV.Elem()
	return t.elemType.SaveByPath(&elemV, path[1:], value)
}
func (t *PtrType) DeleteByPath(v *reflect.Value, path Path) (err error) {
	return fmt.Errorf("[MapType.Delete] not implement,path:%s type:%s", path, v.Type().Kind())
}
