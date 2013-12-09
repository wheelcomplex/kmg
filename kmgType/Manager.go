package kmgType

import (
	"fmt"
	"reflect"
)

type Context struct {
	RootType  KmgType
	RootValue reflect.Value
}

func NewContext(ptr interface{}) (*Context, error) {
	rt := reflect.TypeOf(ptr)
	if rt.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("NewContext need a ptr,but get %T", ptr)
	}
	manager := &Context{}
	manager.RootValue = reflect.ValueOf(ptr)
	et, err := TypeOf(rt)
	if err != nil {
		return nil, err
	}
	manager.RootType = et
	return manager, nil
}
func (m *Context) SaveByPath(path Path, value string) (err error) {
	pEv := &m.RootValue
	err = m.RootType.SaveByPath(pEv, path, value)
	if err != nil {
		return
	}
	if pEv != &m.RootValue {
		err = fmt.Errorf("[manager.SaveByPath] can not save")
		return
	}
	return nil
}

func (m *Context) DeleteByPath(path Path) (err error) {
	pEv := &m.RootValue
	err = m.RootType.DeleteByPath(pEv, path)
	if err != nil {
		return
	}
	if pEv != &m.RootValue {
		err = fmt.Errorf("[manager.DeleteByPath] can not save")
		return
	}
	return nil
}
