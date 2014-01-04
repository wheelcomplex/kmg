package kmgType

import (
	"github.com/bronze1man/kmg/kmgTime"
	"reflect"
	"time"
)

var DateTimeReflectType = reflect.TypeOf((*time.Time)(nil)).Elem()

type DateTimeType struct {
	reflectTypeGetterImp
	saveScaleFromStringer
	saveScaleEditabler
}

func (t *DateTimeType) ToString(v reflect.Value) string {
	return v.Interface().(time.Time).Format(kmgTime.FormatMysql)
}
func (t *DateTimeType) SaveScale(v reflect.Value, value string) error {
	valueT, err := time.Parse(kmgTime.FormatMysql, value)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(valueT))
	return nil
}
