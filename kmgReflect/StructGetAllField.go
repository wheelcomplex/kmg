package kmgReflect

import (
	"reflect"
)

func StructGetAllField(t reflect.Type) []*reflect.StructField {
	fieldMap := map[string]bool{}
	return structGetAllFieldImp(t, fieldMap, []int{})
}

func structGetAllFieldImp(t reflect.Type, fieldMap map[string]bool, indexs []int) (output []*reflect.StructField) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	anonymousFieldList := []*reflect.StructField{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		f.Index = append(indexs, f.Index...)
		if f.Anonymous {
			anonymousFieldList = append(anonymousFieldList, &f)
		}
		if fieldMap[f.Name] {
			continue
		}
		fieldMap[f.Name] = true
		output = append(output, &f)
	}
	for _, f := range anonymousFieldList {
		output = append(output, structGetAllFieldImp(f.Type, fieldMap, f.Index)...)
	}
	return
}

//order of field is not stable...
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
