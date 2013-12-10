package webTypeAdmin

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgType"
	"html/template"
	"reflect"
	//"runtime/debug"
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
	kmgType.StringType
	ctx      *context
	enumList []string
}

func (t *stringEnumType) init() (err error) {
	if t.enumList != nil {
		return
	}
	t.enumList = reflect.Zero(t.GetReflectType()).Interface().(StringEnum).GetEnumList()
	if len(t.enumList) == 0 {
		return fmt.Errorf("stringEnum.GetEnumList() return an array with 0 element")
	}
	return
}
func (t *stringEnumType) HtmlView(v reflect.Value) (html template.HTML, err error) {
	if err = t.init(); err != nil {
		return
	}
	valueS := v.String()
	return theTemplate.ExecuteNameToHtml("Select", selectTemplateData{
		List:  t.enumList,
		Value: valueS,
	})
}
