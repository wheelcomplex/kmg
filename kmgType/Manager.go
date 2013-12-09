package kmgType

import (
	"fmt"
	"reflect"
)

type Manager struct {
	rootType  KmgType
	rootValue reflect.Value
}

func NewManager(ptr interface{}) (*Manager, error) {
	rt := reflect.TypeOf(ptr)
	if rt.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("NewManager need a ptr,but get %T", ptr)
	}
	manager := &Manager{}
	manager.rootValue = reflect.ValueOf(ptr)
	et, err := TypeOf(rt)
	if err != nil {
		return nil, err
	}
	manager.rootType = et
	return manager, nil
}
func (m *Manager) SaveByPath(path Path, value string) (err error) {
	pEv := &m.rootValue
	err = m.rootType.SaveByPath(pEv, path, value)
	if err != nil {
		return
	}
	if pEv != &m.rootValue {
		err = fmt.Errorf("[manager.SaveByPath] can not save")
		return
	}
	return nil
}
