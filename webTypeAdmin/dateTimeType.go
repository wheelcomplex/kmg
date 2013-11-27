package webTypeAdmin

import (
	"github.com/bronze1man/kmg/kmgTime"
	"html/template"
	"reflect"
	"time"
)

type dateTimeType struct {
	commonType
}

func (t *dateTimeType) Html(v reflect.Value) template.HTML {
	return theTemplate.MustExecuteNameToHtml("TextInput", v.Interface().(time.Time).Format(kmgTime.FormatMysql))
}

func (t *dateTimeType) save(v reflect.Value, value string) error {
	valueT, err := time.Parse(kmgTime.FormatMysql, value)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(valueT))
	return nil
}
