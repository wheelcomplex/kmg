package kmgReflect

import (
	"reflect"
)

func StructGetAllFieldMap(t reflect.Type) map[string]*reflect.StructField {
	fieldMap := map[string]*reflect.StructField{}
	structGetAllFieldMapImp(t, fieldMap, []int{})
	return fieldMap
}
func structGetAllFieldMapImp(t reflect.Type, fieldMap map[string]*reflect.StructField, indexs []int) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}
	anonymousFieldList := []*reflect.StructField{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		f.Index = append(indexs, f.Index...)
		if f.Anonymous {
			anonymousFieldList = append(anonymousFieldList, &f)
		}
		_, ok := fieldMap[f.Name]
		if ok {
			continue
		}
		fieldMap[f.Name] = &f
	}
	for _, f := range anonymousFieldList {
		structGetAllFieldMapImp(f.Type, fieldMap, f.Index)
	}
}
