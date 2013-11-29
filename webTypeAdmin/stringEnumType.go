package webTypeAdmin

import (
	"fmt"
	"html/template"
	"reflect"
)

var stringEnumReflectType = reflect.TypeOf((*StringEnum)(nil)).Elem()

/*
define a stringEnum like this
type Enum1 string
func (enum Enum1)GetEnumList()[]string{
	return []string{
		"abc",
		"bcd",
		"qwe",
	}
}
*/
type StringEnum interface {
	GetEnumList() []string //GetEnumList get all list of this type enum ,should only depend on type
}

type stringEnumType struct {
	commonType
	enumList []string
}

func (t *stringEnumType) init() {
	if t.enumList != nil {
		return
	}
	t.enumList = reflect.Zero(t.getReflectType()).Interface().(StringEnum).GetEnumList()
	if len(t.enumList) == 0 {
		panic("stringEnum.GetEnumList() return an array with 0 element")
	}
}
func (t *stringEnumType) Html(v reflect.Value) template.HTML {
	t.init()
	valueS := v.String()
	return theTemplate.MustExecuteNameToHtml("Select", selectTemplateData{
		List:  t.enumList,
		Value: valueS,
	})
}
func (t *stringEnumType) save(v reflect.Value, value string) error {
	t.init()
	if !isInStringSlice(t.enumList, value) {
		return fmt.Errorf("[stringEnumType.save] save value not in enum list value:%s", value)
	}
	v.SetString(value)
	return nil
}
func isInStringSlice(arr []string, target string) bool {
	for _, k := range arr {
		if target == k {
			return true
		}
	}
	return false
}
