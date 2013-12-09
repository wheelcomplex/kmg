package kmgType

import (
	"fmt"
	"reflect"
)

//path -> field name
//TODO embed field
type StructType struct {
	reflectTypeGetterImp
	getElemByStringEditorabler
}

func (t *StructType) GetElemByString(v reflect.Value, k string) (ev reflect.Value, et KmgType, err error) {
	ev = v.FieldByName(k)
	if !ev.IsValid() {
		err = fmt.Errorf("field %s not find in struct", k)
		return
	}
	et, err = TypeOf(ev.Type())
	return
}
