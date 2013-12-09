package kmgType

import (
	"fmt"
	"reflect"
)

//path -> key(Key type)
//key can be bool,string,stringEnum,int,float,
type MapType struct {
	reflectTypeGetterImp
	keyStringConverter StringConverterInterface
	keyType            KmgType
	elemType           KmgType
}

func (t *MapType) init() (err error) {
	if t.keyStringConverter != nil {
		return
	}
	t.keyType, err = TypeOf(t.GetReflectType().Key())
	if err != nil {
		return err
	}
	var ok bool
	t.keyStringConverter, ok = t.keyType.(StringConverterInterface)
	if !ok {
		return fmt.Errorf(
			"mapType key type not implement stringConverterType,key: %s",
			t.keyType.GetReflectType().Kind().String(),
		)
	}
	t.elemType, err = TypeOf(t.GetReflectType().Elem())
	if err != nil {
		return err
	}
	return nil
}
func (t *MapType) SaveByPath(v *reflect.Value, path Path, value string) (err error) {
	err = t.init()
	if err != nil {
		return
	}
	if len(path) == 0 {
		return fmt.Errorf("[mapType.save] get map with no path, value:%s", value)
	}
	OriginCanSet := v.CanSet()
	if v.IsNil() {
		if OriginCanSet {
			v.Set(reflect.MakeMap(t.GetReflectType()))
		} else {
			*v = reflect.MakeMap(t.GetReflectType())
		}
	} else {
		//copy all exist data,if this value can not set
		//this is not need some place
		if !OriginCanSet {
			output := reflect.MakeMap(t.GetReflectType())
			output.Set(*v)
			*v = output
		}
	}
	vk, err := t.keyStringConverter.FromString(path[0])
	if err != nil {
		return err
	}
	saveElemV := v.MapIndex(vk)
	KeyNotExist := false
	if !saveElemV.IsValid() {
		saveElemV = reflect.New(t.elemType.GetReflectType()).Elem()
		KeyNotExist = true
	}
	elemV := &saveElemV
	err = t.elemType.SaveByPath(elemV, path[1:], value)
	if err != nil {
		return err
	}
	if KeyNotExist {
		v.SetMapIndex(vk, saveElemV)
	}
	if elemV != &saveElemV {
		v.SetMapIndex(vk, *elemV)
	}
	return nil
}

func (t *MapType) DeleteByPath(v *reflect.Value, path Path) (err error) {
	return fmt.Errorf("[MapType.Delete] not implement,path:%s type:%s", path, v.Type().Kind())
}
