package ajkApi

import "reflect"

type ApiFuncMeta struct {
	IsMethod     bool
	Func         reflect.Value
	AttachObject reflect.Value // the object method attached
	MethodName   string
}
